package model

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var db *gorm.DB

type (
	todoModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	transformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	userModel struct {
		gorm.Model
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	transformedUser struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}

	// claims is for JWT authentication
	// claims struct {
	// 	ID  uint      `json:"id"`
	// 	Exp time.Time `json:"exp"`
	// }
)

func init() {
	_ = godotenv.Load()

	hostname := os.Getenv("HOST")
	dbname := os.Getenv("DBNAME")
	username := os.Getenv("USER")
	password := os.Getenv("PASSWORD")

	dbString := "host=" + hostname + " user=" + username + " dbname=" + dbname + " sslmode=disable password=" + password
	fmt.Println(dbString)
	var err error
	db, err = gorm.Open("postgres", dbString)
	if err != nil {
		panic("Unable to connect to DB")
	}

	db.AutoMigrate(&todoModel{})
	db.AutoMigrate(&userModel{})
}
