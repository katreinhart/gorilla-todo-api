package controller

import (
	"bytes"
	"net/http"

	"github.com/katreinhart/gorilla-api/model"
)

// CreateUser handles registration of new user
func CreateUser(w http.ResponseWriter, r *http.Request) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	_, err := model.CreateUser(b)

	if err != nil {
		if err.Error() == "Unable to parse input" {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Please check your inputs and try again."))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong."))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully created!"))
}
