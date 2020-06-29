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
	docSettings := chidoc.NewDocSettings("Title for API", chidoc.RapidRender)
	// Here, you set up models, that you gonna use documention, like
	/*
	  schema:
	   "$ref": "#/components/schemes/HTTPResponse"
	*/
	docSettings.SetDefinitions(Response{})

	// Here adds security
	docSettings.SetAuths(chidoc.NewAuthBasic("auth", "this is a simple auth"))

	if err := chidoc.AddRouteDoc(router, "", docSettings); err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}
