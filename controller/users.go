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

	js, err := model.CreateUser(b)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if err.Error() == "Unable to parse input" {
			w.WriteHeader(http.StatusBadRequest)
			// w.Write([]byte("Please check your inputs and try again."))
			return
		} else if err.Error() == "User exists" {
			w.WriteHeader(http.StatusBadRequest)
			// w.Write([]byte("User already exists in db."))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		// w.Write([]byte("Something went wrong."))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// LoginUser function handles request/response of login function
func LoginUser(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	js, err := model.LoginUser(b)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if err.Error() == "Something went wrong with JWT" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"message\": \"something went wrong\"}"))
		} else if err.Error() == "Passwords do not match" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("{\"message\": \"Not allowed\"}"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"message\": \"something went wrong\"}"))
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// FetchAllUsers does what it says on the tin
func FetchAllUsers(w http.ResponseWriter, r *http.Request) {
	js, err := model.FetchAllUsers()

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if err.Error() == "Not found" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(js)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong"))
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
