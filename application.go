package main

import (
	"fmt"

	"github.com/RaphaSalomao/alura-challenge-backend/database"
	"github.com/RaphaSalomao/alura-challenge-backend/router"
)

func main() {
	err := database.Connect()
	if err != nil {
		panic(err)
	}
	router.HandleRequests()

	fmt.Println("Exiting...")
}
