package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/joho/godotenv"
)

var db *gorm.DB

type todoModel struct {
	gorm.Model
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type transformedTodo struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func init() {
	err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	hostname := os.Getenv("HOST")
	dbname := os.Getenv("DBNAME")
	username := os.Getenv("USER")
	password := os.Getenv("PASSWORD")

	dbString := "host=" + hostname + " user=" + username + " dbname=" + dbname + " sslmode=disable password=" + password
	fmt.Println(dbString)
	db, err = gorm.Open("postgres", dbString)
	if err != nil {
		panic("Unable to connect to DB")
	}

	db.AutoMigrate(&todoModel{})
}

func main() {
	environment := os.Getenv("ENVIRONMENT")
	var port string
	if environment == "development" {
		port = ":3000"
	} else if environment == "production" {
		port = ":8080"
	}

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/todos", fetchAllTodos).Methods("GET")
	r.HandleFunc("/todos", createTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", fetchSingleTodo).Methods("GET")
	r.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	http.ListenAndServe(port, handlers.RecoveryHandler()(loggedRouter))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func fetchAllTodos(w http.ResponseWriter, r *http.Request) {
	var todos []todoModel
	var _todos []transformedTodo

	db.Find(&todos)

	if len(todos) <= 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Todos not found"))
		return
	}

	for _, item := range todos {
		_todos = append(_todos, transformedTodo{ID: item.ID, Completed: item.Completed, Title: item.Title})
	}

	js, err := json.Marshal(_todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	var todo todoModel

	err := json.Unmarshal(b, &todo)

	if err != nil {
		panic("unable to marshal input into todoModel")
	}

	db.Save(&todo)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Todo successfully created!"))
}

func fetchSingleTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var todo todoModel
	db.First(&todo, id)

	if todo.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Todo not found"))
		return
	}

	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: todo.Completed}

	js, err := json.Marshal(_todo)
	if err != nil {
		panic("Unable to convert todo to JSON format")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	var todo, updatedTodo todoModel
	db.First(&todo, id)

	if todo.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Todo not found"))
		return
	}

	err := json.Unmarshal(b, &updatedTodo)
	if err != nil {
		panic("Unable to marshal todo into json")
	}

	db.Model(&todo).Update("title", updatedTodo.Title)
	db.Model(&todo).Update("completed", updatedTodo.Completed)

	js, err := json.Marshal(&todo)
	if err != nil {
		panic("Unable to marshal todo into json")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var todo todoModel
	db.First(&todo, id)

	if todo.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Todo not found"))
		return
	}

	db.Delete(&todo)

	js, err := json.Marshal(&todo)
	if err != nil {
		panic("Unable to marshal todo into json")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
