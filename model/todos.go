package model

// FetchAll is the model function which interfaces with the DB and returns a []byte of the todos in json format.
func FetchAll() ([]TransformedTodo, error) {
	var todos []TodoModel
	var _todos []TransformedTodo

	db.Find(&todos)

	if len(todos) <= 0 {
		return nil, ErrorNotFound
	}

	for _, item := range todos {
		_todos = append(_todos, TransformedTodo{ID: item.ID, Completed: item.Completed, Title: item.Title})
	}

	return _todos, nil
}

// Create creates a new todo item and returns the []byte json object and an error.
func Create(todo TodoModel) TransformedTodo {
	db.Save(&todo)
	_todo := TransformedTodo{ID: todo.ID, Title: todo.Title, Completed: todo.Completed}
	return _todo
}

// FetchSingle gets a single todo based on param passed, returning []byte and error
func FetchSingle(id string) (TransformedTodo, error) {

	var todo TodoModel
	db.First(&todo, id)

	if todo.ID == 0 {
		return TransformedTodo{}, ErrorNotFound
	}

	_todo := TransformedTodo{ID: todo.ID, Title: todo.Title, Completed: todo.Completed}

	return _todo, nil
}

// Update is the model function for PUT
func Update(updatedTodo TodoModel, id string) (TransformedTodo, error) {

	var todo TodoModel
	db.Find(&todo, "id = ?", id)

	if todo.ID == 0 {
		return TransformedTodo{}, ErrorNotFound
	}

	db.Model(&todo).Update("title", updatedTodo.Title)
	db.Model(&todo).Update("completed", updatedTodo.Completed)

	return TransformedTodo{ID: todo.ID, Completed: todo.Completed, Title: todo.Title}, nil
}

// Delete deletes the todo from the database
func Delete(id string) (TransformedTodo, error) {

	var todo TodoModel
	db.First(&todo, id)

	if todo.ID == 0 {
		return TransformedTodo{}, ErrorNotFound
	}

	db.Delete(&todo)
	return TransformedTodo{ID: todo.ID, Completed: todo.Completed, Title: todo.Title}, nil
}
