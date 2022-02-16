package controller

import (
	"encoding/json"
	"net/http"

	"github.com/RaphaSalomao/alura-challenge-backend/model"
	"github.com/RaphaSalomao/alura-challenge-backend/service"
	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/gorilla/mux"
)

// Health Check
// @Description return server status
// @Tags Health
// @Success 200
// @Failure 404
// @Router /budget-control/api/v1/health [get]
func Health(w http.ResponseWriter, r *http.Request) {
	utils.HandleResponse(w, http.StatusOK, struct{ Online bool }{true})
}

// Create User
// @Description create a new user
// @Tags User
// @Param user body model.UserRequest true "User"
// @Success 201 {object} model.UserRequest
// @Router /budget-control/api/v1/user [post]
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

// Authenticate
// @Description authenticate user
// @Tags User
// @Param user body model.UserRequest true "User"
// @Success 201 {string} string "token"
// @Router /budget-control/api/v1/authenticate [post]
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

// Month Balance Sumary
// @Description get month balance sumary
// @Tags Balance
// @Param year path string true "Year"
// @Param month path string true "Month"
// @Success 200 {object} model.BalanceSumaryResponse
// @Router /budget-control/api/v1/balance/{year}/{month} [get]
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
