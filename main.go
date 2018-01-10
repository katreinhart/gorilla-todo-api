package main

import (
	"net/http"
	"os"

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
	r.HandleFunc("/todos", controller.FetchAllTodos).Methods("GET")
	r.HandleFunc("/todos", controller.CreateTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", controller.FetchSingleTodo).Methods("GET")
	r.HandleFunc("/todos/{id}", controller.UpdateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", controller.DeleteTodo).Methods("DELETE")

	s := r.PathPrefix("/users").Subrouter()

	s.HandleFunc("/register", controller.CreateUser).Methods("POST")
	s.HandleFunc("/login", controller.LoginUser).Methods("POST")

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	http.ListenAndServe(":"+port, handlers.RecoveryHandler()(loggedRouter))
}
