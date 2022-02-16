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

// Create Receipt
// @Summary Create a new receipt
// @Description Create a new receipt
// @Tags Receipts
// @Param receipt body model.ReceiptRequest true "Receipt"
// @Success 201 {object} uuid.UUID
// @Router /budget-control/api/v1/receipt [post]
func CreateReceipt(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	var receipt model.ReceiptRequest
	json.NewDecoder(r.Body).Decode(&receipt)
	id, err := service.ReceiptService.CreateReceipt(&receipt, userId)
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct {
			Error string
			Id    uuid.UUID
		}{err.Error(), id})
	} else {
		utils.HandleResponse(w, http.StatusCreated, struct{ Id uuid.UUID }{id})
	}
}

// Find All Receipts
// @Summary Find all receipts
// @Description Find all receipts
// @Tags Receipts
// @Success 200 {array} model.ReceiptResponse
// @Router /budget-control/api/v1/receipt [get]
func FindAllReceipts(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	var receipts []model.ReceiptResponse
	description := r.URL.Query().Get("description")
	service.ReceiptService.FindAllReceipts(&receipts, description, userId)
	utils.HandleResponse(w, http.StatusOK, receipts)
}

// Find Receipt By Id
// @Summary Find a receipt by id
// @Description Find a receipt by id
// @Tags Receipts
// @Param id path string true "Receipt id"
// @Success 200 {object} model.ReceiptResponse
// @Router /budget-control/api/v1/receipt/{id} [get]
func FindReceipt(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	var receipt model.ReceiptResponse
	id := uuid.MustParse(mux.Vars(r)["id"])
	err := service.ReceiptService.FindReceipt(&receipt, id, userId)
	if err != nil {
		utils.HandleResponse(w, http.StatusNotFound, struct {
			Error string
			Id    uuid.UUID
		}{err.Error(), id})
	} else {
		utils.HandleResponse(w, http.StatusOK, receipt)
	}
}

// Update Receipt
// @Summary Update a receipt
// @Description Update a receipt
// @Tags Receipts
// @Param id path string true "Receipt id"
// @Param receipt body model.ReceiptRequest true "Receipt"
// @Success 200 {object} uuid.UUID
// @Router /budget-control/api/v1/receipt/{id} [put]
func UpdateReceipt(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	var receipt model.ReceiptRequest
	id := uuid.MustParse(mux.Vars(r)["id"])
	json.NewDecoder(r.Body).Decode(&receipt)
	id, err := service.ReceiptService.UpdateReceipt(&receipt, id, userId)
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct {
			Error string
			Id    uuid.UUID
		}{err.Error(), id})
	} else {
		utils.HandleResponse(w, http.StatusNoContent, nil)
	}
}

// Delete Receipt
// @Summary Delete a receipt
// @Description Delete a receipt
// @Tags Receipts
// @Param id path string true "Receipt id"
// @Success 204
// @Router /budget-control/api/v1/receipt/{id} [delete]
func DeleteReceipt(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	id := uuid.MustParse(mux.Vars(r)["id"])
	service.ReceiptService.DeleteReceipt(id, userId)
	utils.HandleResponse(w, http.StatusNoContent, nil)
}

// Find All Receipts By Period
// @Summary Find all receipts by Period
// @Description Find all receipts by Period
// @Tags Receipts
// @Param year path int true "Year"
// @Param month path int true "Month"
// @Success 200 {array} model.ReceiptResponse
// @Router /budget-control/api/v1/receipt/{year}/{month} [get]
func ReceiptsByPeriod(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	var receipts []model.ReceiptResponse
	vars := mux.Vars(r)
	err := service.ReceiptService.ReceiptsByPeriod(&receipts, vars["year"], vars["month"], userId)
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct{ Error string }{err.Error()})
	} else {
		utils.HandleResponse(w, http.StatusOK, receipts)
	}
}
