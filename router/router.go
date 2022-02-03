package router

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/RaphaSalomao/alura-challenge-backend/controller"
	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/gorilla/mux"
)

var (
	PermitAll       bool
	unauthenticated = map[string]bool{
		"/budget-control/api/v1/health": true,
	}
)

func HandleRequests(permitAll bool) {
	PermitAll = permitAll
	router := mux.NewRouter()
	router.Use(middleware)

	router.HandleFunc("/budget-control/api/v1/health", controller.Health).Methods("GET")
	router.HandleFunc("/budget-control/api/v1/user", controller.CreateUser).Methods("POST")
	router.HandleFunc("/budget-control/api/v1/authenticate", controller.Authenticate).Methods("POST")

	router.HandleFunc("/budget-control/api/v1/receipt", controller.CreateReceipt).Methods("POST")
	router.HandleFunc("/budget-control/api/v1/receipt", controller.FindAllReceipts).Methods("GET")
	router.HandleFunc("/budget-control/api/v1/receipt/{id}", controller.FindReceipt).Methods("GET")
	router.HandleFunc("/budget-control/api/v1/receipt/{id}", controller.UpdateReceipt).Methods("PUT")
	router.HandleFunc("/budget-control/api/v1/receipt/{id}", controller.DeleteReceipt).Methods("DELETE")
	router.HandleFunc("/budget-control/api/v1/receipt/{year}/{month}", controller.ReceiptsByPeriod).Methods("GET")

	router.HandleFunc("/budget-control/api/v1/expense", controller.CreateExpense).Methods("POST")
	router.HandleFunc("/budget-control/api/v1/expense", controller.FindAllExpenses).Methods("GET")
	router.HandleFunc("/budget-control/api/v1/expense/{id}", controller.FindExpense).Methods("GET")
	router.HandleFunc("/budget-control/api/v1/expense/{id}", controller.UpdateExpense).Methods("PUT")
	router.HandleFunc("/budget-control/api/v1/expense/{id}", controller.DeleteExpense).Methods("DELETE")
	router.HandleFunc("/budget-control/api/v1/expense/{year}/{month}", controller.ExpensesByPeriod).Methods("GET")

	router.HandleFunc("/budget-control/api/v1/summary/{year}/{month}", controller.MonthBalanceSumary).Methods("GET")

	go http.ListenAndServe(":8080", router)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Server is running")
	<-quit

	fmt.Println("Server down.")
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")

			if !unauthenticated[r.URL.Path] && !PermitAll {
				token := r.Header.Get("Authorization")
				if token == "" {
					utils.HandleResponse(w, http.StatusUnauthorized, struct{ Error string }{Error: "Missing token"})
					return
				}
			}
			next.ServeHTTP(w, r)
		},
	)
}
