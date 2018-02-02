package controller

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katreinhart/gorilla-api/model"
)

// FetchAllTodos fetches from model and returns json
func FetchAllTodos(w http.ResponseWriter, r *http.Request) {
	var _todos []model.TransformedTodo

	_todos, err := model.FetchAll()

	js, err := json.Marshal(_todos)

	handleErrorAndRespond(js, err, w)
}

// CreateTodo takes request body and sends it to model, sending back success message or error on response
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	// read stuff from the request
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	var todo model.TodoModel
	var _todo model.TransformedTodo

	err := json.Unmarshal(b, &todo)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
	}

	_todo = model.Create(todo)

	js, err := json.Marshal(_todo)
	handleErrorAndRespond(js, err, w)

}

// FetchSingleTodo takes URL param and passes to model,
func FetchSingleTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_todo, err := model.FetchSingle(id)

	js, err := json.Marshal(_todo)
	handleErrorAndRespond(js, err, w)
}

// UpdateTodo modifies the content of Todo based on url param and body content.
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	var todo model.TodoModel
	err := json.Unmarshal(b, &todo)

	var _todo model.TransformedTodo
	_todo, err = model.Update(todo, id)

	js, err := json.Marshal(_todo)
	handleErrorAndRespond(js, err, w)
}

// DeleteTodo deletes a todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var deletedTodo model.TransformedTodo

	deletedTodo, err := model.Delete(id)

	js, err := json.Marshal(deletedTodo)
	handleErrorAndRespond(js, err, w)
}
