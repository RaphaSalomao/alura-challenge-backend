package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/RaphaSalomao/alura-challenge-backend/controller"
	"github.com/RaphaSalomao/alura-challenge-backend/model/types"
	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/gorilla/mux"
)

var (
	unauthenticated = map[string]bool{
		"/budget-control/api/v1/health":       true,
		"/budget-control/api/v1/user":         true,
		"/budget-control/api/v1/authenticate": true,
	}
)

func HandleRequests() {
	srvPort := fmt.Sprintf(":%s", os.Getenv("SRV_PORT"))

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

	go http.ListenAndServe(srvPort, router)

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
			var userId string
			var err error
			if !unauthenticated[r.URL.Path] {
				token := strings.Split(r.Header.Get("Authorization"), " ")
				if len(token) == 1 && token[0] == "" {
					utils.HandleResponse(w, http.StatusBadRequest, struct{ Error string }{Error: "missing token"})
					return
				} else if token[0] != "Bearer" || len(token) != 2 {
					utils.HandleResponse(w, http.StatusBadRequest, struct{ Error string }{Error: "invalid token"})
					return
				}
				userId, err = utils.ParseToken(token[1])
				if err != nil {
					utils.HandleResponse(w, http.StatusUnauthorized, struct{ Error string }{Error: err.Error()})
					return
				}
			}
			r = r.WithContext(context.WithValue(r.Context(), types.ContextKey("userId"), userId))
			next.ServeHTTP(w, r)
		},
	)
}
