package main

import (
	env "github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	ClientID          string `env:"SPOTIFY_CLIENT_ID"`
	ClientSecret      string `env:"SPOTIFY_CLIENT_SECRET"`
	ClientRedirectURL string `env:"SPOTIFY_CLIENT_REDIRECT_URL" envDefault:"https://spotify.conradludgate.com/spotify_callback"`

	FrontendDir string `env:"SPOTIFY_API_FRONTEND" envDefault:"./client"`

	Database string `env:"SPOTIFY_API_DB_CONN"`

	SessionKey string `env:"SPOTIFY_API_SESSION_KEY"`
	HashCost   int    `env:"SPOTIFY_API_HASH_COST" envDefault:"10"`
}

var config *Config

func init() {
	config = new(Config)
	if err := env.Parse(config); err != nil {
		panic(err)
	}
}
