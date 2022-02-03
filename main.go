package main

import (
	"fmt"

	"github.com/RaphaSalomao/alura-challenge-backend/database"
	"github.com/RaphaSalomao/alura-challenge-backend/router"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	err := database.Connect()
	if err != nil {
		panic(err)
	}
	router.HandleRequests(false)

	fmt.Println("Exiting...")
}
