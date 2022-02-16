package main

import (
	"fmt"

	"github.com/RaphaSalomao/alura-challenge-backend/database"
	"github.com/RaphaSalomao/alura-challenge-backend/router"
)

// @title     Alura Backend Challenge 2nd Edition API
// @version   1.0.0
// @host      localhost:5000
// @BasePath  /
func main() {
	err := database.Connect()
	if err != nil {
		panic(err)
	}
	router.HandleRequests()

	fmt.Println("Exiting...")
}
