package api

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	UseMemosHost  string
	UseMemosToken string
	ListenAddr    string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("USEMEMOS_TOKEN")
	if token == "" {
		log.Fatal("Environment variable USEMEMOS_TOKEN is not set")
	}
	host := os.Getenv("USEMEMOS_HOST")
	if host == "" {
		log.Fatal("Environment variable USEMEMOS_HOST is not set")
	}
	return &Config{
		UseMemosHost:  host,
		UseMemosToken: token,
		ListenAddr:    ":8090",
	}

}
