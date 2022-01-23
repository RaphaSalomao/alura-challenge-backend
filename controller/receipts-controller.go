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

var (
	receiptService = service.ReceiptService{}
)

func CreateReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt model.Receipt
	json.NewDecoder(r.Body).Decode(&receipt)
	id, err := receiptService.CreateReceipt(&receipt)
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
	receiptService.FindAllReceipts(&receipts)
	utils.HandleResponse(w, http.StatusOK, receipts)
}

func FindReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt model.ReceiptResponse
	id := uuid.MustParse(mux.Vars(r)["id"])
	receiptService.FindReceipt(&receipt, id)
	utils.HandleResponse(w, http.StatusOK, receipt)
}

func UpdateReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt model.Receipt
	id := uuid.MustParse(mux.Vars(r)["id"])
	json.NewDecoder(r.Body).Decode(&receipt)
	id, err := receiptService.UpdateReceipt(&receipt, id)
	if err != nil {
		utils.HandleResponse(w, http.StatusUnprocessableEntity, struct {
			Error string
			Id    uuid.UUID
		}{err.Error(), id})
	} else {
		utils.HandleResponse(w, http.StatusOK, struct{ Id uuid.UUID }{receipt.Id})
	}
}
