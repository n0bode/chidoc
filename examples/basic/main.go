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
	Message string      `json:"message,omitempty" docs:"description: Texto da resposta"`
	Data    interface{} `json:"data" docs:"description:Data from request"`
}

func HTTPSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Data: data,
	})
}

// User struct for users
type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" docs:"len:5,required"`
	Age      int    `json:"age"`
	ParentID int64  `json:"parent_id"`
	Weapon   int    `docs:"enum:Weapons"`
}

var data []User = []User{
	User{
		ID:       1,
		Name:     "Thror",
		Age:      248,
		ParentID: 0,
	},
	User{
		ID:       2,
		Name:     "Thrain",
		Age:      206,
		ParentID: 1,
	},
	User{
		ID:       3,
		Name:     "Thorin",
		Age:      195,
		ParentID: 2,
	},
	User{
		ID:       4,
		Name:     "Frerin",
		Age:      48,
		ParentID: 2,
	},
	User{
		ID:       5,
		Name:     "Dis",
		Age:      200,
		ParentID: 2,
	},
	User{
		ID:       6,
		Name:     "Kili",
		Age:      200,
		ParentID: 5,
	},
	User{
		ID:       7,
		Name:     "Fili",
		Age:      200,
		ParentID: 5,
	},
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

// GetChildren returns user's children
// summary: get children by parent
func GetChildren(w http.ResponseWriter, r *http.Request) {
	var userID string = chi.URLParam(r, "userID")

	var users []User
	for _, user := range data {
		if strconv.FormatInt(user.ParentID, 10) == userID {
			users = append(users, user)
		}
	}
	HTTPSuccess(w, users)
}

// GetChildrenByID returns user's children
// summary: get children by parent
func GetChildrenByID(w http.ResponseWriter, r *http.Request) {
	var parentID string = chi.URLParam(r, "userID")
	var userID string = chi.URLParam(r, "childrenID")

	var users []User
	for _, user := range data {
		if strconv.FormatInt(user.ParentID, 10) == parentID && userID == strconv.FormatInt(user.ID, 10) {
			users = append(users, user)
		}
	}

	HTTPSuccess(w, users)
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
	router.Get("/api/user/{userID:[0-9]+}/children", GetChildren)
	router.Get("/api/user/{userID:[0-9]+}/children/{childrenID:[0-9]+}", GetChildrenByID)
	// ..

	// Create doc settings
	// doc settings is a struct to generate YAML format for Redoc
	docSettings := chidoc.NewDocSettings("Title for API", chidoc.RapidRender)
	// Here, you set up models, that you gonna use documention, like
	/*
	  schema:
	   "$ref": "#/components/schemes/HTTPResponse"
	*/
	docSettings.SetDefinitions(Response{}, User{}, chidoc.Enum("Weapons", `
	1 - Hammer
	2 - Staff
	3 - Sword
	4 - Mace
	`, 1, 2, 3, 4))
	docSettings.SetTheme(chidoc.DarkTheme)

	// Here adds security
	docSettings.SetAuths(chidoc.NewAuthAPIKey("Auth", "Token", "Authorization", chidoc.InHeader))

	docSettings.SetLogo(chidoc.ImageFromURL("https://i.imgur.com/7lZu0wq.png"))
	if err := chidoc.AddRouteDoc(router, "/", docSettings, "docs"); err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}
