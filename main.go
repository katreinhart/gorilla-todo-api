package main

import (
	"net/http"
	"os"

	jwtmiddleware "github.com/aiden0z/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/katreinhart/gorilla-api/controller"
	"github.com/katreinhart/gorilla-api/routing"
)

func main() {
	var port string

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = "8080"
	}

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", routing.HomeHandler)

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/todos", controller.FetchAllTodos).Methods("GET")
	api.HandleFunc("/todos", controller.CreateTodo).Methods("POST")
	api.HandleFunc("/todos/{id}", controller.FetchSingleTodo).Methods("GET")
	api.HandleFunc("/todos/{id}", controller.UpdateTodo).Methods("PUT")
	api.HandleFunc("/todos/{id}", controller.DeleteTodo).Methods("DELETE")

	u := r.PathPrefix("/users").Subrouter()
	u.HandleFunc("/", controller.FetchAllUsers).Methods("GET")
	u.HandleFunc("/me", controller.FetchMyInfo).Methods("GET")

	s := r.PathPrefix("/auth").Subrouter()

	s.HandleFunc("/register", controller.CreateUser).Methods("POST")
	s.HandleFunc("/login", controller.LoginUser).Methods("POST")

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		UserProperty:  "user",
	})

	muxRouter := http.NewServeMux()
	muxRouter.Handle("/", r)
	muxRouter.Handle("/api/", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(api),
	))

	muxRouter.Handle("/users/", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(u),
	))

	n := negroni.Classic()
	n.UseHandler(muxRouter)

	http.ListenAndServe(":"+port, handlers.RecoveryHandler()(n))

}
