package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/n0bode/chidoc"
	"github.com/n0bode/chidoc/examples/compose/db"
)

// Response
type Response struct {
	// doc uses json field with default tag
	Message string      `json:"message,omitempty" docs:"description: Texto da resposta"`
	Data    interface{} `json:"data,omitempty" docs:"description:Data from request"`
}

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email"`
}

func HTTPSuccess(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Data: data,
	})
}

func HTTPError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Message: message,
	})
}

// GetAllUsers gets all users in your app
// summary: gets all users in your app
// security:
// - oauth: []
// responses:
//  '200':
//    description: Returns users by id
//    content:
//     application/json:
//      schema:
//       type: array
//       items:
//        "$ref": "#/components/schemes/UserOrm"
func GetAllUsers(conn *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		HTTPSuccess(w, conn.Filter(func(user db.UserOrm) (db.UserOrm, bool) {
			user.Password = ""
			return user, false
		}), 200)
	}
}

// PostUser creates a new user
// summary: creates a new user
// security:
// - Auth: []
// responses:
//  '201':
//    description: Created a new user
//    content:
//     application/json:
//      schema:
//       "$ref": "#/components/schemes/UserOrm"
// requestBody:
//  description: Optional description in *Markdown*
//  required: true
//  content:
//   application/json:
//    schema:
//     $ref: '#/components/schemes/User'
func PostUser(conn *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			HTTPSuccess(w, "payload is not valid", 400)
			return
		}

		users := conn.Filter(func(result db.UserOrm) (db.UserOrm, bool) {
			if result.Username == user.Username {
				return result, false
			}
			// pass all
			return result, true
		})

		if len(users) > 0 {
			HTTPSuccess(w, "username already exists", 400)
			return
		}

		// create user
		HTTPSuccess(w, conn.AddUser(user.Username, user.Password, user.Name, user.Email), 201)
	}
}

// PutUser update fields user
// summary: updates fields user
// security:
// - oauth: []
// responses:
//  '200':
//    description: Updated user
//    content:
//     application/json:
//      schema:
//       "$ref": "#/components/schemes/UserOrm"
//  '400':
//    description: Check message field response
//    content:
//     application/json:
//      schema:
//       "$ref": "#/components/schemes/Response"
// requestBody:
//  description: Optional description in *Markdown*
//  required: true
//  content:
//   application/json:
//    schema:
//     $ref: '#/components/schemes/User'
func PutUser(conn *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

		var update User
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			HTTPSuccess(w, "payload is not valid", 400)
			return
		}

		// get user by id
		user, updated := conn.Update(func(user db.UserOrm) (db.UserOrm, bool) {
			// not found
			if user.ID != userID {
				return user, false
			}

			// update
			user.UpdateAt = time.Now()
			user.Email = update.Email
			user.Name = update.Name
			return user, true
		})

		if !updated {
			HTTPError(w, "user not exists", 400)
			return
		}
		HTTPSuccess(w, user, 201)
	}
}

func PostToken(conn *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, pass, ok := r.BasicAuth()
		if !ok {
			HTTPError(w, "unsupported_grant", 401)
			return
		}

		users := conn.Filter(func(user db.UserOrm) (db.UserOrm, bool) {
			if user.Username == username && user.Password == pass {
				return user, false
			}
			return user, true
		})

		if len(users) == 0 {
			HTTPError(w, "unsupported_grant", 401)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token":  "aaaabbbcc",
			"refresh_token": "aaaabbbcca",
			"token_type":    "Bearer",
			"expires_in":    3600,
		})
	}
}

func Routes() *chi.Mux {
	db := db.NewDB()

	db.AddUser("admin", "admin", "Zeus", "zeus@olympus.com")

	router := chi.NewRouter()

	// ... API
	router.Get("/users", GetAllUsers(db))
	router.Post("/users", PostUser(db))
	router.Put("/users/{id:[0-9]+}", PutUser(db))
	router.Post("/token", PostToken(db))

	return router
}

func main() {
	routes := Routes()

	// Create doc settings
	// doc settings is a struct to generate YAML format for Redoc
	docSettings := chidoc.NewDocSettings("Compose Func Users", chidoc.RapidRender)
	docSettings.SetTheme(chidoc.DarkTheme)
	docSettings.SetDefinitions(db.UserOrm{}, User{}, Response{})

	// Here adds security
	docSettings.SetAuths(chidoc.NewAuthOAuth("oauth", "http://localhost:8000/token", "asd", map[string]string{
		"Roles": "username to acessa page",
		"Test":  "password user",
	}))

	docSettings.SetLogo(chidoc.ImageFromURLScaled("https://i.imgur.com/7lZu0wq.png", 0.5))
	if err := chidoc.AddRouteDoc(routes, "/", docSettings, "docs"); err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8000", routes); err != nil {
		log.Fatal(err)
	}
}
