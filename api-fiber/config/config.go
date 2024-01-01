package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config retrieves the value of an environment variable by its key.
func Config(key string) string {
	// Load environment variables from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		// Print an error message if loading the .env file fails
		fmt.Print("Error loading .env file")
	}

	// Return the value of the specified environment variable
	return os.Getenv(key)
}
