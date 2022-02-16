package main

import (
	"fmt"

	"github.com/RaphaSalomao/alura-challenge-backend/database"
	"github.com/RaphaSalomao/alura-challenge-backend/router"
)

// @title     Alura Backend Challenge 2nd Edition API
// @version   1.0.0
// @host      http://alurachallengebackend2ndedition-env.eba-cmaxmrtx.us-east-2.elasticbeanstalk.com
// @BasePath  /
func main() {
	err := database.Connect()
	if err != nil {
		panic(err)
	}
	router.HandleRequests()

	fmt.Println("Exiting...")
}
