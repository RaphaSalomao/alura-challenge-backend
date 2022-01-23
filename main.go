package main

import (
	"fmt"

	"github.com/RaphaSalomao/alura-challenge-backend/database"
	"github.com/RaphaSalomao/alura-challenge-backend/router"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Server is running")
	godotenv.Load()
	database.Connect()
	router.HandleRequests()
}
