package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/katreinhart/gorilla-api/controller"
	"github.com/katreinhart/gorilla-api/routing"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB

type (
	// todoModel struct {
	// 	gorm.Model
	// 	Title     string `json:"title"`
	// 	Completed bool   `json:"completed"`
	// }

	// transformedTodo struct {
	// 	ID        uint   `json:"id"`
	// 	Title     string `json:"title"`
	// 	Completed bool   `json:"completed"`
	// }

	// userModel struct {
	// 	gorm.Model
	// 	Email    string `json:"email"`
	// 	Password string `json:"password"`
	// }

	// transformedUser struct {
	// 	ID    uint   `json:"id"`
	// 	Email string `json:"email"`
	// 	Token string `json:"token"`
	// }

	token struct {
		Sub uint      `json:"sub"`
		Exp time.Time `json:"exp"`
	}
)

// func init() {
// 	_ = godotenv.Load()

// 	hostname := os.Getenv("HOST")
// 	dbname := os.Getenv("DBNAME")
// 	username := os.Getenv("USER")
// 	password := os.Getenv("PASSWORD")

// 	dbString := "host=" + hostname + " user=" + username + " dbname=" + dbname + " sslmode=disable password=" + password

// 	var err error
// 	db, err = gorm.Open("postgres", dbString)
// 	if err != nil {
// 		panic("Unable to connect to DB")
// 	}

// 	db.AutoMigrate(&todoModel{})
// 	db.AutoMigrate(&userModel{})
// }

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
	// r.HandleFunc("/todos/{id}", fetchSingleTodo).Methods("GET")
	// r.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	// r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")

	// s := r.PathPrefix("/users").Subrouter()

	// s.HandleFunc("/register", createUser).Methods("POST")
	// s.HandleFunc("/login", loginUser).Methods("POST")

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	http.ListenAndServe(":"+port, handlers.RecoveryHandler()(loggedRouter))
}

// func fetchSingleTodo(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	var todo todoModel
// 	db.First(&todo, id)

// 	if todo.ID == 0 {
// 		w.WriteHeader(http.StatusNotFound)
// 		w.Write([]byte("Todo not found"))
// 		return
// 	}

// 	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: todo.Completed}

// 	js, err := json.Marshal(_todo)
// 	if err != nil {
// 		panic("Unable to convert todo to JSON format")
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(js)
// }

// func updateTodo(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(r.Body)
// 	b := []byte(buf.String())

// 	var todo, updatedTodo todoModel
// 	db.First(&todo, id)

// 	if todo.ID == 0 {
// 		w.WriteHeader(http.StatusNotFound)
// 		w.Write([]byte("Todo not found"))
// 		return
// 	}

// 	err := json.Unmarshal(b, &updatedTodo)
// 	if err != nil {
// 		panic("Unable to marshal todo into json")
// 	}

// 	db.Model(&todo).Update("title", updatedTodo.Title)
// 	db.Model(&todo).Update("completed", updatedTodo.Completed)

// 	js, err := json.Marshal(&todo)
// 	if err != nil {
// 		panic("Unable to marshal todo into json")
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(js)
// }

// func deleteTodo(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	var todo todoModel
// 	db.First(&todo, id)

// 	if todo.ID == 0 {
// 		w.WriteHeader(http.StatusNotFound)
// 		w.Write([]byte("Todo not found"))
// 		return
// 	}

// 	db.Delete(&todo)

// 	js, err := json.Marshal(&todo)
// 	if err != nil {
// 		panic("Unable to marshal todo into json")
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(js)
// }

// func createUser(w http.ResponseWriter, r *http.Request) {

// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(r.Body)
// 	b := []byte(buf.String())

// 	var user userModel

// 	err := json.Unmarshal(b, &user)

// 	if err != nil {
// 		panic("unable to marshal input into todoModel")
// 	}

// 	fmt.Println(user)
// 	hash, err := hashPassword(user.Password)
// 	if err != nil {
// 		panic("unable to retrieve password")
// 	}

// 	user.Password = hash
// 	db.Save(&user)

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("User successfully created!"))
// }

// func loginUser(w http.ResponseWriter, r *http.Request) {
// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(r.Body)
// 	b := []byte(buf.String())

// 	var user userModel

// 	err := json.Unmarshal(b, &user)

// 	if err != nil {
// 		panic("unable to marshal input into todoModel")
// 	}

// 	id := user.ID

// 	var dbUser userModel
// 	db.First(&dbUser, id)
// 	if dbUser.ID == 0 {
// 		panic("Unable to find user in db")
// 	}

// 	match := checkPasswordHash(user.Password, dbUser.Password)
// 	if !match {
// 		panic("Passwords do not match")
// 	}

// 	claims := token{Sub: user.ID, Exp: time.Now().Add(time.Hour * 24)}
// 	token, err := jwt.ParseWithClaims()

// 	var _user transformedUser
// 	_user.Email = user.Email
// 	_user.ID = user.ID
//  //	_user.Token =
// }

// user login password helper functions
// from https://gowebexamples.com/password-hashing/

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
