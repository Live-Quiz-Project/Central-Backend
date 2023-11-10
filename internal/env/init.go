package env

import (
	"log"

	"github.com/joho/godotenv"
)

func Initialize() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
