package chidoc

import (
	"bytes"
	"encoding/json"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"image/png"
	"io"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/go-chi/chi"
)

var htmls = map[DocRender]string{
	"rapidoc": `
		<head>
			<title> {title} </title>
			<link rel="icon" type="image/png" href="{url_docs}/favicon.png">
			<!-- Include javascript redoc lib -->
			<link href="https://fonts.googleapis.com/css2?family=Nunito:wght@300;600&amp;family=Open+Sans:wght@300;600&amp;family=Roboto+Mono&amp;display=swap" rel="stylesheet">
			<script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>

		</head>
		<body>
			<rapi-doc 
				spec-url=".{path}/{docs}" 
				mono-font="{theme.fontname}" 
				regular-font="{theme.fonttype}" 
				text-color="{theme.textcolor}" 
				bg-color="{theme.backgroundcolor}" 
				theme="{theme.schema}" 
				render-style="{theme.renderstyle}"
				font-size="{theme.fontsize}"
				show-header="{theme.header}"
				primary-color="{theme.primarycolor}"
				header-color="{theme.headercolor}"
				schema-style="{theme.schematype}"
				nav-bg-color="" 
				nav-text-color="" 
				nav-hover-bg-color="" 
				nav-hover-text-color="" 
				nav-accent-color=""
			> 
			<img 
    			slot="nav-logo" 
    			src=".{path}/{logo}"
  			/> 
			</rapi-doc>
		</body>
	`,
	"redoc": `
		<head>
			<title> {title} </title>
			<link rel="icon" type="image/png" href="{url_docs}/favicon.png">
			<!-- Include javascript redoc lib -->
			<script type="text/javascript" src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"></script>
		</head>
		<body>
			<!-- Redoc UI shows here below -->
			<div id="redoc_ui"></div>
			<!-- Init redoc UI -->
			<script type="text/javascript">
				Redoc.init(".{path}/{docs}", {settings}, document.getElementById("redoc_ui")); 
			</script>
		</body>
	`,
}

func themeToList(prefix string, theme Theme) (arr []string) {
	t := reflect.ValueOf(theme)
	for i := 0; i < t.Type().NumField(); i++ {
		f := t.Type().FieldByIndex([]int{i})
		v := t.FieldByIndex([]int{i})
		var key string = "{" + prefix + "." + strings.ToLower(f.Name) + "}"
		if f.Tag.Get("doc") == "attribute" {
			if v.IsZero() {
				continue
			}
			arr = append(arr, key, "\""+f.Name+"\"=\""+v.String()+"\"")
			continue
		}
		arr = append(arr, key, v.String())
	}
	return arr
}

func replaceHTML(html, title, path string, settings *DocSettings) string {
	dumps, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return ""
	}

	themeAtts := themeToList("theme", settings.Theme)

	attrs := []string{
		"{title}", title,
		"{path}", path,
		"{logo}", "logo.png",
		"{docs}", "docs.yaml",
		"{settings}", string(dumps),
	}

	r := strings.NewReplacer(
		append(attrs, themeAtts...)...,
	)
	return r.Replace(html)
}

func splitFuncName(name string) string {
	var arr []string = strings.Split(name, ".")
	return strings.Split(arr[len(arr)-1], "-")[0]
}

func infoFunc(handler http.Handler) (name, filename string, line int) {
	valueOf := reflect.ValueOf(handler)
	funcPC := runtime.FuncForPC(valueOf.Pointer())
	filename, line = funcPC.FileLine(0)
	return splitFuncName(funcPC.Name()), filename, line
}

func routeDescription(handler http.Handler, tmp map[string][]*ast.CommentGroup) (map[string]interface{}, error) {
	fname, filename, _ := infoFunc(handler)

	comments, exists := tmp[filename]
	if !exists {
		fset := token.NewFileSet()
		parse, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		tmp[filename] = parse.Comments
		comments = parse.Comments
	}

	var data string
	for _, group := range comments {
		if !strings.HasPrefix(group.Text(), fname) {
			continue
		}
		data = "#" + group.Text()
	}

	description := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(data), &description); err != nil {
		return nil, err
	}

	if description != nil {
		if _, exists := description["tags"]; !exists {
			description["tags"] = []string{}
		}
	}
	return description, nil
}

func walkRoute(parent string, p map[string]interface{}, parseTMP map[string][]*ast.CommentGroup, r chi.Routes) (map[string]interface{}, error) {
	for _, route := range r.Routes() {
		var pattern string = route.Pattern
		if strings.HasSuffix(pattern, "/*") {
			pattern = pattern[:len(pattern)-2]
		}
		var path string = parent + pattern
		if route.SubRoutes == nil {
			doc := make(map[string]interface{})
			for method, handler := range route.Handlers {
				d, err := routeDescription(handler, parseTMP)
				if err != nil {
					return nil, err
				}
				doc[strings.ToLower(method)] = d
			}
			p[path] = doc
			continue
		}
		walkRoute(path, p, parseTMP, route.SubRoutes)
	}

	return p, nil
}

func parseTag(tag string) (m map[string]string) {
	var token string
	var key string = "name"

	// Flag to mark endline
	tag += ","
	m = make(map[string]string)
	for i := 0; i < len(tag); i++ {
		if tag[i] == ',' {
			if len(key) == 0 {
				key = token
				token = ""
			}
			m[key] = token
			key = ""
			continue
		}

		if tag[i] == ':' {
			key = token
			token = ""
			continue
		}
		token += string(tag[i])
	}
	return
}

func isIntType(t reflect.Type) bool {
	return t.Kind() >= reflect.Int && t.Kind() <= reflect.Uint64
}

func isFloatType(t reflect.Type) bool {
	return t.Kind() >= reflect.Float32 && t.Kind() <= reflect.Float64
}

func isArrType(t reflect.Type) bool {
	return t.Kind() == reflect.Array || t.Kind() == reflect.Slice
}

func parseDefinition(schemes, m map[string]interface{}, t reflect.Type) map[string]interface{} {
	// if it was a pointer
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch {
	case isIntType(t):
		m["type"] = "integer"
	case isFloatType(t):
		m["type"] = "number"
	case t.Kind() == reflect.Bool:
		m["type"] = "boolean"
	case isArrType(t):
		m["type"] = "array"
		m["items"] = parseDefinition(schemes, make(map[string]interface{}), t.Elem())
	case t.Kind() == reflect.String:
		m["type"] = "string"
	case t == reflect.TypeOf(time.Time{}):
		m["type"] = "string"
		m["format"] = "date-time"
	case t.Kind() == reflect.Struct:
		var req []string
		props := make(map[string]interface{})

		// Stop recursive
		if _, exists := schemes[t.Name()]; exists {
			m["$ref"] = "#/components/schemes/" + t.Name()
			break
		}
		schemes[t.Name()] = "not recursive here!"

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Anonymous {
				continue
			}
			var name string = f.Name
			aa := make(map[string]interface{})
			tagJSON := parseTag(f.Tag.Get("json"))
			name = tagJSON["name"]

			if tagJSON["name"] == "-" {
				continue
			}

			docs := parseTag(f.Tag.Get("docs"))
			if _, required := docs["required"]; required {
				req = append(req)
			}

			if description, exists := docs["description"]; exists {
				aa["description"] = description
			}
			ff := parseDefinition(schemes, aa, f.Type)
			props[name] = ff
		}

		m["type"] = "object"
		// Properties
		if len(props) != 0 {
			m["properties"] = props
		}

		// Required fields
		if len(req) != 0 {
			m["required"] = req
		}
		schemes[t.Name()] = m
	default:
		m["type"] = "object"
	}
	return m
}

func genRouteYAML(settings *DocSettings, r *chi.Mux) (doc string, err error) {
	paths, err := walkRoute("", make(map[string]interface{}), make(map[string][]*ast.CommentGroup), r)
	if err != nil {
		return doc, err
	}

	// Parse definitions to YAML
	schemes := make(map[string]interface{})
	for _, d := range settings.definitions {
		var t reflect.Type = reflect.TypeOf(d)
		parseDefinition(schemes, make(map[string]interface{}), t)
	}
	settings.Set("components.schemes", schemes)

	// Parse authorization to YAML
	auths := make(map[string]interface{})
	for _, a := range settings.auths {
		if err = a.Decode(auths); err != nil {
			return doc, err
		}
	}
	settings.Set("components.securitySchemes", auths)
	settings.Set("paths", paths)

	raw := make(map[string]interface{})
	if err = settings.Decode(raw); err != nil {
		return doc, err
	}

	buffer, err := yaml.Marshal(raw)
	return string(buffer), err
}

func readImage(handle HandlerImage, logo io.Writer) error {
	if handle != nil {
		image := handle()
		return png.Encode(logo, image)
	}
	return errors.New("Handle is nil")
}

// AddRouteDoc adds documention to route
func AddRouteDoc(root *chi.Mux, docpath string, settings *DocSettings) error {
	var urlDoc string = docpath

	var html string = replaceHTML(htmls[settings.Render], settings.Title, urlDoc, settings)

	// Create page index
	root.Get(docpath, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	})

	// Read static logo
	var logo bytes.Buffer
	if err := readImage(settings.handlerLogo, &logo); err == nil {
		// Set logo
		settings.Set("info.x-logo.url", docpath+"/logo.png")
		//Adds logo router
		root.Get(docpath+"/logo.png", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "image/png")
			w.Write(logo.Bytes())
		})
	}

	// Read static icon
	var icon bytes.Buffer
	if err := readImage(settings.handlerIcon, &icon); err != nil {
		root.Get(docpath+"/favicon.png", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "image/png")
			w.Write(icon.Bytes())
		})
	}

	docs, err := genRouteYAML(settings, root)
	if err != nil {
		return err
	}

	// Create route for docs generation
	root.Get(docpath+"/docs.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/x-yaml")
		w.Write([]byte(docs))
	})
	return nil
}
