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

func CreateReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt model.Receipt
	json.NewDecoder(r.Body).Decode(&receipt)
	id, err := service.ReceiptService.CreateReceipt(&receipt)
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct {
			Error string
			Id    uuid.UUID
		}{err.Error(), id})
	} else {
		utils.HandleResponse(w, http.StatusCreated, struct{ Id uuid.UUID }{receipt.Id})
	}
}

func FindAllReceipts(w http.ResponseWriter, r *http.Request) {
	var receipts []model.ReceiptResponse
	description := strings.ToUpper(r.URL.Query().Get("description"))
	service.ReceiptService.FindAllReceipts(&receipts, description)
	utils.HandleResponse(w, http.StatusOK, receipts)
}

func FindReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt model.ReceiptResponse
	id := uuid.MustParse(mux.Vars(r)["id"])
	err := service.ReceiptService.FindReceipt(&receipt, id)
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
	var receipt model.ReceiptRequest
	id := uuid.MustParse(mux.Vars(r)["id"])
	json.NewDecoder(r.Body).Decode(&receipt)
	id, err := service.ReceiptService.UpdateReceipt(&receipt, id)
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
	id := uuid.MustParse(mux.Vars(r)["id"])
	service.ReceiptService.DeleteReceipt(id)
	utils.HandleResponse(w, http.StatusNoContent, nil)
}

func ReceiptsByPeriod(w http.ResponseWriter, r *http.Request) {
	var receipts []model.ReceiptResponse
	vars := mux.Vars(r)
	err := service.ReceiptService.ReceiptsByPeriod(&receipts, vars["year"], vars["month"])
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct{ Error string }{err.Error()})
	} else {
		utils.HandleResponse(w, http.StatusOK, receipts)
	}
}
