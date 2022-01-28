package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RaphaSalomao/alura-challenge-backend/model"
	"github.com/RaphaSalomao/alura-challenge-backend/service"
	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func CreateExpense(w http.ResponseWriter, r *http.Request) {
	var expense model.ExpenseRequest
	json.NewDecoder(r.Body).Decode(&expense)
	id, err := service.ExpenseService.CreateExpense(&expense)
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct {
			Error string
			Id    uuid.UUID
		}{err.Error(), id})
	} else {
		utils.HandleResponse(w, http.StatusCreated, struct{ Id uuid.UUID }{id})
	}
}

func FindAllExpenses(w http.ResponseWriter, r *http.Request) {
	var expenses []model.ExpenseResponse
	desc := strings.ToUpper(r.URL.Query().Get("description"))
	service.ExpenseService.FindAllExpenses(&expenses, desc)
	utils.HandleResponse(w, http.StatusOK, expenses)
}

func FindExpense(w http.ResponseWriter, r *http.Request) {
	var expense model.ExpenseResponse
	id := uuid.MustParse(mux.Vars(r)["id"])
	err := service.ExpenseService.FindExpense(&expense, id)
	if err != nil {
		utils.HandleResponse(w, http.StatusNotFound, struct {
			Error string
			Id    uuid.UUID
		}{err.Error(), id})
	} else {
		utils.HandleResponse(w, http.StatusOK, expense)
	}
}

func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	var expense model.Expense
	id := uuid.MustParse(mux.Vars(r)["id"])
	json.NewDecoder(r.Body).Decode(&expense)
	id, err := service.ExpenseService.UpdateExpense(&expense, id)
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct {
			Error string
			Id    uuid.UUID
		}{err.Error(), id})
	} else {
		utils.HandleResponse(w, http.StatusOK, struct{ Id uuid.UUID }{expense.Id})
	}
}

func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(mux.Vars(r)["id"])
	service.ExpenseService.DeleteExpense(id)
	utils.HandleResponse(w, http.StatusNoContent, nil)
}

func ExpensesByPeriod(w http.ResponseWriter, r *http.Request) {
	var expenses []model.ExpenseResponse
	vars := mux.Vars(r)
	err := service.ExpenseService.ExpensesByPeriod(&expenses, vars["year"], vars["month"])
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct{ Error string }{err.Error()})
	} else {
		utils.HandleResponse(w, http.StatusOK, expenses)
	}
}
