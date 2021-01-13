package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/oauth2"
)

var db *gorm.DB

func ConnectDB() {
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
	Name       string

	AccessToken  string
	RefreshToken string
	Expires      time.Time
	TokenType    string
	Scopes       []string

	UserID string
}

func (t *Token) IntoOauth() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expires,
		TokenType:    t.TokenType,
	}
}

func (t *Token) FromOauth(token *oauth2.Token) {
	t.AccessToken = token.AccessToken
	t.RefreshToken = token.RefreshToken
	t.Expires = token.Expiry
	t.TokenType = token.TokenType
}
