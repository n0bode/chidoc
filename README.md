# Documentation generator for go-chi (xi)

## How is it works?

I have a problem to generate documentation for my API's,
so, i thought with myself, why don't i create a lib for auto-generate
documentation for my API's. Well, this lib is that!

## YAML route

Unfortunately, i didn't a better method to replace this:
```go
// [NAME_FUNCTION] comentary for SDK documentation
// # Use a space between slash and your YAML command 
// # Here begin YAML
// summary: my route to do anything
// # Pay attention, use space to indent
// responses:
//  '200':
//   description: type here
//   schema:
//    $ref: #/components/schemes/Response
```

## Example
```go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/n0bode/chidoc"
)

// Response
type Response struct {
	// doc uses json field with default tag
	Message string `json:"message"`
}

// GETSay says hello for anyone
// summary: Endpoint to say hello world
// responses:
//  '200':
//    description: It said hello world for you
//    content:
//     application/json:
//      schema:
//       "$ref": "#/components/schemes/Response"
func GETSay(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Response{
		Message: "Hi! don't worry, it's running",
	})
}

func main() {
	router := chi.NewRouter()
	// ... API
	router.Get("/api/say", GETSay)

	// ..

	// Create doc settings
	// doc settings is a struct to generate YAML format for Redoc
	docSettings := chidoc.NewDocSettings("Title for API")
	// Here, you set up models, that you gonna use documention, like
	/*
	  schema:
	   "$ref": "#/components/schemes/Response"
	*/
	docSettings.SetDefinitions(Response{})

	// Here adds security
	docSettings.SetAuths(chidoc.NewAuthBasic("auth", "this is a simple auth"))

	if err := chidoc.AddRouteDoc(router, "/my-doc-path", docSettings); err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}

```
