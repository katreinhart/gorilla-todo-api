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

	var user userModel

	err := json.Unmarshal(b, &user)

	if err != nil {
		return []byte("Unable to parse input"), errors.New("Unable to parse input")
	}

	fmt.Println(user)
	hash, err := hashPassword(user.Password)
	if err != nil {
		return []byte("Something went wrong"), errors.New("Something went wrong")
	}

	user.Password = hash
	db.Save(&user)

	return []byte("User created successfully"), nil
}

// LoginUser function takes in request body of login post and checks it against database.
// Successful login returns user with token; unsuccessful login returns an http status error with message.
func LoginUser(b []byte) ([]byte, error) {
	fmt.Println("Login User function in model")

	// usermodel is a struct of email and password values
	var user userModel

	err := json.Unmarshal(b, &user)

	if err != nil {
		panic("unable to marshal input into todoModel")
	}

	userEmail := user.Email

	var dbUser userModel
	db.First(&dbUser, "email = ?", userEmail)
	if dbUser.ID == 0 {
		panic("Unable to find user in db")
	}

	match := checkPasswordHash(user.Password, dbUser.Password)
	if !match {
		panic("Passwords do not match")
	}
	exp := time.Now().Add(time.Hour * 24).Unix()

	claim := jwt.StandardClaims{Id: string(dbUser.ID), ExpiresAt: exp}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	t, err := token.SignedString(os.Getenv("SECRET"))

	var _user transformedUser
	_user.Email = user.Email
	_user.ID = user.ID
	_user.Token = t

	js, err := json.Marshal(_user)

	if err != nil {
		// handle error
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
