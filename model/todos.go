package model

import (
	"encoding/json"
	"errors"
)

// FetchAll is the model function which interfaces with the DB and returns a []byte of the todos in json format.
func FetchAll() ([]byte, error) {

	var todos []todoModel
	var _todos []transformedTodo

	db.Find(&todos)

	if len(todos) <= 0 {
		err := errors.New("Not found")
		{
			return []byte("Todos not found"), err
		}
	}

	for _, item := range todos {
		_todos = append(_todos, transformedTodo{ID: item.ID, Completed: item.Completed, Title: item.Title})
	}

	js, err := json.Marshal(_todos)

	{
		return js, err
	}
}

// Create creates a new todo item and returns the []byte json object and an error.
func Create(b []byte) ([]byte, error) {

	var todo todoModel

	err := json.Unmarshal(b, &todo)

	if err != nil {
		return []byte("Something went wrong"), err
	}

	db.Save(&todo)

	return []byte("Todo successfully created"), nil
}

// FetchSingle gets a single todo based on param passed, returning []byte and error
func FetchSingle(id string) ([]byte, error) {

	var todo todoModel
	db.First(&todo, id)

	if todo.ID == 0 {
		err := errors.New("Not found")
		return []byte("Todo not found"), err
	}

	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: todo.Completed}

	js, err := json.Marshal(_todo)
	if err != nil {
		js = []byte("Unable to convert todo to JSON format")
	}

	return js, err
}

// Update is the model function for PUT
func Update(b []byte, id string) ([]byte, error) {

	var todo, updatedTodo todoModel
	db.First(&todo, id)

	if todo.ID == 0 {
		err := errors.New("Not found")
		return []byte("Todo not found"), err
	}

	err := json.Unmarshal(b, &updatedTodo)
	if err != nil {
		return []byte("Malformed input"), err
	}

	db.Model(&todo).Update("title", updatedTodo.Title)
	db.Model(&todo).Update("completed", updatedTodo.Completed)

	js, err := json.Marshal(&todo)
	if err != nil {
		return []byte("Unable to marshal json"), err
	}

	return js, nil
}

// Delete deletes the todo from the database
func Delete(id string) ([]byte, error) {

	var todo todoModel
	db.First(&todo, id)

	if todo.ID == 0 {
		// w.WriteHeader(http.StatusNotFound)
		// w.Write([]byte("Todo not found"))
		// return
	}

	db.Delete(&todo)

	js, err := json.Marshal(&todo)
	if err != nil {
		panic("Unable to marshal todo into json")
	}

	return js, nil
}
