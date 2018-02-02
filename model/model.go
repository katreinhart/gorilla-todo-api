package model

import (
	"errors"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var db *gorm.DB

type (
	// TodoModel is the GORM model for todo
	TodoModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	// TransformedTodo is the version that is sent back in response
	TransformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	// UserModel is the GORM model for the user database
	UserModel struct {
		gorm.Model
		Email    string `json:"email"`
		Password string `json:"password"`
		Admin    bool   `json:"admin"`
	}

	// TransformedUser is the version sent back on register/login
	TransformedUser struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}

	// ListedUser is the format sent back for GET all user queries
	ListedUser struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Admin bool   `json:"admin"`
	}

	// CustomClaims is used in JWT
	CustomClaims struct {
		UID uint `json:"uid"`
		Rol bool `json:"rol"`
		jwt.StandardClaims
	}
)

// ErrorNotFound is used for a 404
var ErrorNotFound = errors.New("Not found")

// ErrorMalformedInput is used for a 400
var ErrorMalformedInput = errors.New("Unable to parse input")

// ErrorUserExists is used when user's email already is in database
var ErrorUserExists = errors.New("User exists in DB")

// ErrorInternalServer handles general 500 type errors
var ErrorInternalServer = errors.New("Something went wrong")

// ErrorNotAllowed is bad login data
var ErrorNotAllowed = errors.New("Forbidden")

func init() {
	_ = godotenv.Load()

	hostname := os.Getenv("HOST")
	dbname := os.Getenv("DBNAME")
	username := os.Getenv("USER")
	password := os.Getenv("PASSWORD")

	dbString := "host=" + hostname + " user=" + username + " dbname=" + dbname + " sslmode=disable password=" + password

	var err error
	db, err = gorm.Open("postgres", dbString)
	if err != nil {
		panic("Unable to connect to DB")
	}

	db.AutoMigrate(&TodoModel{})
	db.AutoMigrate(&UserModel{})
}
