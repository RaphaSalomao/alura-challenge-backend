package main

import (
	"fmt"
	"strings"

	"github.com/RaphaSalomao/alura-challenge-backend/database"
	"github.com/RaphaSalomao/alura-challenge-backend/router"
	"github.com/joho/godotenv"
)

// @title     Alura Backend Challenge 2nd Edition API
// @version   1.0.1
// @host      localhost
// @BasePath  /
func main() {
	err := godotenv.Load()
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		fmt.Println("Error loading .env file, using default environment variables")
	}
	err = database.Connect()
	if err != nil {
		panic(err)
	}
	router.HandleRequests()

	fmt.Println("Exiting...")
}
