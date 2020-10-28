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
			<link rel="icon" type="image/png" href="{url_icon}">
			<!-- Include javascript redoc lib -->
			<link href="https://fonts.googleapis.com/css2?family=Nunito:wght@300;600&amp;family=Open+Sans:wght@300;600&amp;family=Roboto+Mono&amp;display=swap" rel="stylesheet">
			<script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
		</head>
		<body>
			<rapi-doc 
				spec-url=".{url_docs}" 
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
    			src=".{url_logo}"
  			/> 
			</rapi-doc>
		</body>
	`,
	"redoc": `
		<head>
			<title> {title} </title>
			<link rel="icon" type="image/png" href="{url_icon}">
			<!-- Include javascript redoc lib -->
			<script type="text/javascript" src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"></script>
		</head>
		<body>
			<!-- Redoc UI shows here below -->
			<div id="redoc_ui"></div>
			<!-- Init redoc UI -->
			<script type="text/javascript">
				Redoc.init(".{url_docs}", {settings}, document.getElementById("redoc_ui")); 
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
		"{url_logo}", joinPath(path, "logo.png"),
		"{url_icon}", joinPath(path, "favicon.ico"),
		"{url_docs}", joinPath(path, "docs.yaml"),
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
			description["tags"] = []string{"API"}
		}
	}
	return description, nil
}

func parseRoutePattern(pattern string) (string, []map[string]interface{}) {
	params := make([]map[string]interface{}, 0)
	var path string

	for _, subName := range strings.Split(pattern[1:], "/") {
		if subName[0] != '{' || subName[len(subName)-1] != '}' {
			path += "/" + subName
			continue
		}

		var index int = strings.Index(subName, ":")
		if index < 0 {
			path += "/" + subName
			continue
		}

		var name string = subName[1:index]
		var format string = subName[index+1 : len(subName)-1]

		params = append(params, map[string]interface{}{
			"in":       "path",
			"name":     name,
			"required": true,
			"schema": map[string]interface{}{
				"type":   "string",
				"format": format,
			},
		})

		path += "/{" + name + "}"
	}
	return path, params
}

func walkRoute(parent string, p map[string]interface{}, parseTMP map[string][]*ast.CommentGroup, r chi.Routes) (map[string]interface{}, error) {
	for _, route := range r.Routes() {
		pattern, params := parseRoutePattern(route.Pattern)

		if strings.HasSuffix(pattern, "favicon.png") {
			continue
		}

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

				// add parameters
				if params != nil {
					d["parameters"] = params
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

// parseTag parse format docs:"description: TITLE,required" to
// map[string]string
// 	descritption: " TITLE"
//	required: ""
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

// isIntType checks if type is a interger
func isIntType(t reflect.Type) bool {
	return t.Kind() >= reflect.Int && t.Kind() <= reflect.Uint64
}

// isFloatType checks if type is a number
func isFloatType(t reflect.Type) bool {
	return t.Kind() >= reflect.Float32 && t.Kind() <= reflect.Float64
}

// isArrType checks if type is a arr
func isArrType(t reflect.Type) bool {
	return t.Kind() == reflect.Array || t.Kind() == reflect.Slice
}

// parseDefinitions parse definition models for a map[Type]
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

		// Stop recursivePlug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }
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
	settings.Set("tags", [](interface{}){
		map[string]interface{}{
			"name": "API",
		},
	})
	raw := make(map[string]interface{})
	if err = settings.Decode(raw); err != nil {
		return doc, err
	}

	buffer, err := yaml.Marshal(raw)
	return string(buffer), err
}

func readImage(handle HandlerImage, logo io.Writer) error {
	if handle != nil {
		if err := png.Encode(logo, handle()); err != nil {
			return err
		}
		return nil
	}
	return errors.New("Handle is nil")
}

func joinPath(p0, p1 string) string {
	l := len(p0)
	if p0[l-1] == '/' {
		return p0 + p1
	}
	return p0 + "/" + p1
}

// AddRouteDoc adds documention to route
func AddRouteDoc(root *chi.Mux, docpath string, settings *DocSettings) error {
	var urlDoc string = docpath

	var html string = replaceHTML(htmls[settings.Render], settings.Title, urlDoc, settings)

	// set logo swagger
	settings.Set("info.x-logo.url", joinPath(docpath, "logo.png"))

	docs, err := genRouteYAML(settings, root)
	if err != nil {
		return err
	}

	// Create page index
	root.Get(docpath, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	})

	// Read static logo
	var logo bytes.Buffer
	err = readImage(settings.handlerLogo, &logo)
	if err == nil {
		// Set logo
		//Adds logo router
		root.Get(joinPath(docpath, "logo.png"), func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "image/png")
			w.Write(logo.Bytes())
		})
	}

	// Read static icon
	var icon bytes.Buffer
	if err := readImage(settings.handlerIcon, &icon); err != nil {
		root.Get(joinPath(docpath, "favicon.png"), func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "image/png")
			w.Write(icon.Bytes())
		})
	}

	// Create route for docs generation
	root.Get(joinPath(docpath, "docs.yaml"), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/x-yaml")
		w.Write([]byte(docs))
	})
	return nil
}
