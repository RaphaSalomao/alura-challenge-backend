package controller

import (
	"net/http"

	"github.com/RaphaSalomao/alura-challenge-backend/model"
	"github.com/RaphaSalomao/alura-challenge-backend/service"
	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/gorilla/mux"
)

func MonthBalanceSumary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var balanceSumary model.BalanceSumaryResponse
	err := service.BalanceSumary(&balanceSumary, vars["year"], vars["month"])
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct{ Error string }{Error: err.Error()})
	} else {
		utils.HandleResponse(w, http.StatusOK, balanceSumary)
	}
}
