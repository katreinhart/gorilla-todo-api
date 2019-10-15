package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/katreinhart/gorilla-todo-api/model"
)

// CreateUser handles registration of new user
func CreateUser(w http.ResponseWriter, r *http.Request) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	var user model.UserModel

	err := json.Unmarshal(b, &user)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
	}

	_user, err := model.CreateUser(user)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
	}

	js, err := json.Marshal(_user)
	handleErrorAndRespond(js, err, w)

}

// LoginUser function handles request/response of login function
func LoginUser(w http.ResponseWriter, r *http.Request) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	var u model.UserModel
	err := json.Unmarshal(b, &u)

	var _u model.TransformedUser
	_u, err = model.LoginUser(u)

	js, err := json.Marshal(_u)

	handleErrorAndRespond(js, err, w)
}

// FetchAllUsers does what it says on the tin
func FetchAllUsers(w http.ResponseWriter, r *http.Request) {

	// testing context setting from jwt middleware
	user := r.Context().Value("user")
	tok := user.(*jwt.Token)

	for k, v := range tok.Claims.(jwt.MapClaims) {
		fmt.Println(k, v)
	}
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

// FetchMyInfo gets info about the current user based on contents of token
func FetchMyInfo(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	tok := user.(*jwt.Token)

	if tok == nil {
		panic("token not found")
	}

	fmt.Fprintf(os.Stderr, "Header:\n%v\n", tok.Header)
	fmt.Fprintf(os.Stderr, "Claims:\n%v\n", tok.Claims)
	claims := tok.Claims.(jwt.MapClaims)
	fmt.Fprintf(os.Stderr, "UID:\n%v\n", claims["uid"])
	uid, ok := claims["uid"].(float64)

	if !ok {
		panic("what the fuck")
	}
	js, err := model.FetchMyInfo(uid)

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
