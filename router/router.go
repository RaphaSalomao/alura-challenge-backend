package router

import (
	"net/http"

	"github.com/RaphaSalomao/alura-challenge-backend/controller"
	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/gorilla/mux"
)

func HandleRequests() {
	router := mux.NewRouter()
	router.Use(middleware)

	router.HandleFunc("/budget-control/api/v1/health", health).Methods("GET")

	router.HandleFunc("/budget-control/api/v1/receipt", controller.CreateReceipt).Methods("POST")
	router.HandleFunc("/budget-control/api/v1/receipt", controller.FindAllReceipts).Methods("GET")
	router.HandleFunc("/budget-control/api/v1/receipt/{id}", controller.FindReceipt).Methods("GET")
	router.HandleFunc("/budget-control/api/v1/receipt/{id}", controller.UpdateReceipt).Methods("PUT")
	http.ListenAndServe(":8080", router)
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		},
	)
}

func health(w http.ResponseWriter, r *http.Request) {
	utils.HandleResponse(w, http.StatusOK, struct{ Online bool }{true})
}
