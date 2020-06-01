package chidoc

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/go-chi/chi"
)

const (
	// PageHTML sad
	pageHTML = `
		<head>
			<title> {title} </title>
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
	`

	headerYAML = `
	swagger: '3.0'
	schemes:
		- https
		- http
	info:
		{info}
	`
)

func replaceHTML(title, urlDocs string, settings *DocSettings) string {
	dumps, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return ""
	}
	r := strings.NewReplacer("{title}", title, "{url_docs}", urlDocs, "{settings}", string(dumps))
	return r.Replace(pageHTML)
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
			var name string = t.Name()

			tagJSON := parseTag(f.Tag.Get("json"))
			if tagJSON["name"] != "-" {
				name = tagJSON["name"]
			}

			docs := parseTag(f.Tag.Get("docs"))
			if _, required := docs["required"]; required {
				req = append(req)
			}

			ff := parseDefinition(schemes, make(map[string]interface{}), f.Type)
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

// AddRouteDoc adds documention to route
func AddRouteDoc(root *chi.Mux, docpath string, settings *DocSettings) error {
	var html string = replaceHTML(settings.Title, docpath+"/docs.yaml", settings)
	docs, err := genRouteYAML(settings, root)
	if err != nil {
		return err
	}

	// Create page index
	root.Get(docpath, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	})

	// Create route for docs generation
	root.Get(docpath+"/docs.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/x-yaml")
		w.Write([]byte(docs))
	})
	return nil
}
