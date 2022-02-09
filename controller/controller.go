package controller

import (
	"encoding/json"
	"net/http"

	"github.com/RaphaSalomao/alura-challenge-backend/model"
	"github.com/RaphaSalomao/alura-challenge-backend/service"
	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/gorilla/mux"
)

func Health(w http.ResponseWriter, r *http.Request) {
	utils.HandleResponse(w, http.StatusOK, struct{ Online bool }{true})
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.UserRequest
	json.NewDecoder(r.Body).Decode(&user)
	err := service.CreateUser(&user)
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct{ Error string }{Error: err.Error()})
	} else {
		utils.HandleResponse(w, http.StatusCreated, user)
	}
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	var user model.UserRequest
	json.NewDecoder(r.Body).Decode(&user)
	token, err := service.Authenticate(&user)
	if err != nil {
		utils.HandleResponse(w, http.StatusUnauthorized, struct{ Error string }{Error: err.Error()})
	} else {
		utils.HandleResponse(w, http.StatusOK, struct{ Token string }{Token: token})
	}
}

func MonthBalanceSumary(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	vars := mux.Vars(r)
	var balanceSumary model.BalanceSumaryResponse
	err := service.BalanceSumary(&balanceSumary, vars["year"], vars["month"], userId)
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct{ Error string }{Error: err.Error()})
	} else {
		utils.HandleResponse(w, http.StatusOK, balanceSumary)
	}
}
