package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tapiaw38/practiq-be/internal/platform/config"
)

func loadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}
	config.InitConfigService()
}
