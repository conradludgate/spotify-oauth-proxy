package main

import (
	"encoding/base64"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("postgres", config.Database)
	if err != nil {
		log.Fatalln("Could not connect to database:", err)
	}

	db.AutoMigrate(&Session{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Token{})
}

type Session struct {
	ID     string
	UserID string `gorm:"UNIQUE"`

	AccessToken  string
	RefreshToken string
	Expires      time.Time
	TokenType    string
}

type User struct {
	ID      string
	Tokens  []Token
	Session Session
}

type Token struct {
	ID         string
	APIKeyHash []byte

	AccessToken  string
	RefreshToken string
	Expires      time.Time
	TokenType    string

	UserID string
}

func GetTokenIfValid(id, apiKey string) *Token {
	token := new(Token)
	db.Where("id = ?", id).First(token)
	if token.ID == "" {
		return nil
	}

	apiKeyBytes, err := base64.StdEncoding.DecodeString(apiKey)
	if err != nil {
		return nil
	}

	if bcrypt.CompareHashAndPassword(token.APIKeyHash, apiKeyBytes) != nil {
		return nil
	}
	return token
}
