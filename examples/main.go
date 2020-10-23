package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/n0bode/chidoc"
)

// Response
type Response struct {
	// doc uses json field with default tag
	Message   string     `json:"message" docs:"description: Texto da resposta"`
	Responses []Response `json:"responses" docs:"description:Lista de respostas"`
}

// User struct for users
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// GETSay says hello for anyone
// summary: Endpoint to say hello world
// security:
// - Auth: []
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

// GETUserByID gets user by id
// summary: get users by id
// security:
// - Auth: []
// responses:
//  '200':
//    description: Returns users by id
//    content:
//     application/json:
//      schema:
//       "$ref": "#/components/schemes/User"
//  '400':
//    description: User is dead
//    content:
//     application/json:
//      schema:
//       "$ref": "#/components/schemes/Response"
func GETUserByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	age, _ := strconv.ParseInt(userID, 10, 32)

	w.Header().Add("Content-Type", "application/json")
	if (10 * age) >= 100 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "I think this users is dead!",
		})
		return
	}

	json.NewEncoder(w).Encode(User{
		Name: "Long Jhonson",
		Age:  int(10 * age),
	})
}

func main() {
	router := chi.NewRouter()
	// ... API
	router.Get("/api/say", GETSay)
	router.Get("/api/user/{userID:[0-9]+}", GETUserByID)
	// ..

	// Create doc settings
	// doc settings is a struct to generate YAML format for Redoc
	docSettings := chidoc.NewDocSettings("Title for API", chidoc.RedocRender)
	// Here, you set up models, that you gonna use documention, like
	/*
	  schema:
	   "$ref": "#/components/schemes/HTTPResponse"
	*/
	docSettings.SetDefinitions(Response{}, User{})
	docSettings.SetTheme(chidoc.DarkTheme)

	// Here adds security
	docSettings.SetAuths(chidoc.NewAuthAPIKey("Auth", "Token", "Authorization", chidoc.InHeader))

	docSettings.SetLogo(chidoc.ImageFromURL("https://i.imgur.com/7lZu0wq.png"))
	if err := chidoc.AddRouteDoc(router, "/docs", docSettings); err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}
