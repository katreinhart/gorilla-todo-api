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

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		if err.Error() == "Unable to parse input" {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Please check your inputs and try again."))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong."))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully created!"))
}

// LoginUser function handles request/response of login function
func LoginUser(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	js, err := model.LoginUser(b)

	if err != nil {
		// handle errors
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
