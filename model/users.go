package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser handles registration of new user
func CreateUser(b []byte) ([]byte, error) {

	var user, dbUser userModel

	err := json.Unmarshal(b, &user)

	if err != nil {
		return []byte(""), errors.New("Unable to parse input")
	}
	email := user.Email
	db.First(&dbUser, "email = ?", email)
	if dbUser.ID != 0 {
		return []byte(""), errors.New("User exists")
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return []byte(""), errors.New("Something went wrong")
	}

	user.Password = hash

	db.Save(&user)

	if err != nil {
		return []byte(""), errors.New("User exists")
	}

	js, err := json.Marshal(user)

	return js, nil
}

// LoginUser function takes in request body of login post and checks it against database.
// Successful login returns user with token; unsuccessful login returns an http status error with message.
func LoginUser(b []byte) ([]byte, error) {

	// usermodel is a struct of email and password values
	var user userModel

	err := json.Unmarshal(b, &user)

	if err != nil {
		return []byte(""), errors.New("unable to unmarshal input into todoModel")
	}

	userEmail := user.Email

	var dbUser userModel

	db.First(&dbUser, "email = ?", userEmail)
	if dbUser.ID == 0 {
		return []byte(""), errors.New("user with that email already exists in database")
	}

	match := checkPasswordHash(user.Password, dbUser.Password)
	if !match {
		return []byte(""), errors.New("passwords do not match")
	}

	// jwt stuff
	exp := time.Now().Add(time.Hour * 24).Unix()
	claim := jwt.StandardClaims{Id: string(dbUser.ID), ExpiresAt: exp}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	secret := []byte(os.Getenv("SECRET"))

	t, err := token.SignedString(secret)

	if err != nil {
		fmt.Println(err.Error())
		return []byte("Something went wrong with JWT"), err
	}

	fmt.Println("token is", t)

	var _user transformedUser
	_user.Email = user.Email
	_user.ID = user.ID
	_user.Token = t

	js, err := json.Marshal(_user)

	if err != nil {
		fmt.Println(err.Error())
		return []byte("Error parsing user into JSON"), err
	}

	return js, nil
}

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

// FetchAllUsers asdfasdfasdf
func FetchAllUsers() ([]byte, error) {
	var users []userModel
	var _users []listedUser

	db.Find(&users)

	if len(users) <= 0 {
		err := errors.New("Not found")
		{
			return []byte("Users not found"), err
		}
	}

	for _, item := range users {
		_users = append(_users, listedUser{ID: item.ID, Email: item.Email, Admin: false})
	}

	js, err := json.Marshal(_users)

	{
		return js, err
	}
}
