package model

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser handles registration of new user
func CreateUser(u UserModel) (TransformedUser, error) {

	var dbUser UserModel
	var _user TransformedUser

	db.First(&dbUser, "email = ?", u.Email)

	if dbUser.ID != 0 {
		return TransformedUser{}, ErrorUserExists
	}

	hash, err := hashPassword(u.Password)
	if err != nil {
		return TransformedUser{}, ErrorInternalServer
	}

	u.Password = hash

	db.Save(&u)

	t, err := createTokenAndSign(u)

	if err != nil {
		return TransformedUser{}, err
	}
	_user = TransformedUser{ID: u.ID, Email: u.Email, Token: t}

	return _user, nil
}

// LoginUser function takes in request body of login post and checks it against database.
// Successful login returns user with token; unsuccessful login returns an http status error with message.
func LoginUser(u UserModel) (TransformedUser, error) {

	var dbUser UserModel

	db.First(&dbUser, "email = ?", u.Email)
	if dbUser.ID == 0 {
		return TransformedUser{}, ErrorUserExists
	}

	match := checkPasswordHash(u.Password, dbUser.Password)
	if !match {
		return TransformedUser{}, ErrorNotAllowed
	}

	t, err := createTokenAndSign(dbUser)

	if err != nil {
		return TransformedUser{}, err
	}

	_user := TransformedUser{Email: u.Email, ID: u.ID, Token: t}

	return _user, nil
}

// FetchAllUsers asdfasdfasdf
func FetchAllUsers() ([]byte, error) {
	var users []UserModel
	var _users []ListedUser

	db.Find(&users)

	if len(users) <= 0 {
		err := errors.New("Not found")
		{
			return []byte("Users not found"), err
		}
	}

	for _, item := range users {
		_users = append(_users, ListedUser{ID: item.ID, Email: item.Email, Admin: false})
	}

	js, err := json.Marshal(_users)

	{
		return js, err
	}
}

// FetchMyInfo finds the given user in the db and returns info about them
func FetchMyInfo(uid float64) ([]byte, error) {
	var user UserModel
	var _user ListedUser
	struid := strconv.FormatFloat(uid, 'f', -1, 64)

	db.First(&user, "id = ?", struid)
	if user.ID == 0 {
		return []byte(""), errors.New("User not found")
	}

	_user = ListedUser{ID: user.ID, Email: user.Email, Admin: user.Admin}

	js, err := json.Marshal(_user)
	return js, err
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

func createTokenAndSign(u UserModel) (string, error) {
	// jwt stuff
	exp := time.Now().Add(time.Hour * 24).Unix()
	claim := CustomClaims{
		u.ID,
		u.Admin,
		jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	secret := []byte(os.Getenv("SECRET"))

	t, err := token.SignedString(secret)

	return t, err
}
