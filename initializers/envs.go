package initializers

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Println(err)
		panic("Failed to load environmental variables")
	}
}
