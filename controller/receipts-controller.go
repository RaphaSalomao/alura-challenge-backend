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

func FindAllReceipts(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	var receipts []model.ReceiptResponse
	description := r.URL.Query().Get("description")
	service.ReceiptService.FindAllReceipts(&receipts, description, userId)
	utils.HandleResponse(w, http.StatusOK, receipts)
}

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

func DeleteReceipt(w http.ResponseWriter, r *http.Request) {
	userId := utils.UserIdFromContext(r.Context())
	id := uuid.MustParse(mux.Vars(r)["id"])
	service.ReceiptService.DeleteReceipt(id, userId)
	utils.HandleResponse(w, http.StatusNoContent, nil)
}

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
