package controller

import (
	"encoding/json"
	"net/http"

	"github.com/RaphaSalomao/alura-challenge-backend/model"
	"github.com/RaphaSalomao/alura-challenge-backend/service"
	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Create Expense
// @Summary Create a new expense
// @Description Create a new expense
// @Tags Expenses
// @Param expense body model.ExpenseRequest true "Expense"
// @Success 201 {object} uuid.UUID
// @Router /budget-control/api/v1/expense [post]
func CreateExpense(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	var expense model.ExpenseRequest
	json.NewDecoder(r.Body).Decode(&expense)
	id, err := service.ExpenseService.CreateExpense(&expense, userId)
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct {
			Error string
			Id    uuid.UUID
		}{err.Error(), id})
	} else {
		utils.HandleResponse(w, http.StatusCreated, struct{ Id uuid.UUID }{id})
	}
}

// Find All Expenses
// @Summary Find all expenses
// @Description Find all expenses
// @Tags Expenses
// @Success 200 {array} model.ExpenseResponse
// @Router /budget-control/api/v1/expense [get]
func FindAllExpenses(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	var expenses []model.ExpenseResponse
	description := r.URL.Query().Get("description")
	service.ExpenseService.FindAllExpenses(&expenses, description, userId)
	utils.HandleResponse(w, http.StatusOK, expenses)
}

// Find Expense By Id
// @Summary Find expense by id
// @Description Find expense by id
// @Tags Expenses
// @Param id path string true "Expense ID"
// @Success 200 {object} model.ExpenseResponse
// @Router /budget-control/api/v1/expense/{id} [get]
func FindExpense(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	var expense model.ExpenseResponse
	id := uuid.MustParse(mux.Vars(r)["id"])
	err := service.ExpenseService.FindExpense(&expense, id, userId)
	if err != nil {
		utils.HandleResponse(w, http.StatusNotFound, struct {
			Error string
			Id    uuid.UUID
		}{err.Error(), id})
	} else {
		utils.HandleResponse(w, http.StatusOK, expense)
	}
}

// Update Expense
// @Summary Update an expense
// @Description Update an expense
// @Tags Expenses
// @Param id path string true "Expense ID"
// @Param expense body model.ExpenseRequest true "Expense"
// @Success 204
// @Router /budget-control/api/v1/expense/{id} [put]
func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	var expense model.ExpenseRequest
	id := uuid.MustParse(mux.Vars(r)["id"])
	json.NewDecoder(r.Body).Decode(&expense)
	id, err := service.ExpenseService.UpdateExpense(&expense, id, userId)
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct {
			Error string
			Id    uuid.UUID
		}{err.Error(), id})
	} else {
		utils.HandleResponse(w, http.StatusNoContent, nil)
	}
}

// Delete Expense
// @Summary Delete an expense
// @Description Delete an expense
// @Tags Expenses
// @Param id path string true "Expense ID"
// @Success 204
// @Router /budget-control/api/v1/expense/{id} [delete]
func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	id := uuid.MustParse(mux.Vars(r)["id"])
	service.ExpenseService.DeleteExpense(id, userId)
	utils.HandleResponse(w, http.StatusNoContent, nil)
}

// Find All Expenses By Period
// @Summary Find all expenses by period
// @Description Find all expenses by period
// @Tags Expenses
// @Param year path int true "Year"
// @Param month path int true "Month"
// @Success 200 {array} model.ExpenseResponse
// @Router /budget-control/api/v1/expense/{year}/{month} [get]
func ExpensesByPeriod(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	var expenses []model.ExpenseResponse
	vars := mux.Vars(r)
	err := service.ExpenseService.ExpensesByPeriod(&expenses, vars["year"], vars["month"], userId)
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct{ Error string }{err.Error()})
	} else {
		utils.HandleResponse(w, http.StatusOK, expenses)
	}
}
