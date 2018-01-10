package model

import (
	"encoding/json"
	"errors"
	"fmt"

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
