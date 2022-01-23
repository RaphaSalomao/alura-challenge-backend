package service

import (
	"errors"

	"github.com/RaphaSalomao/alura-challenge-backend/database"
	"github.com/RaphaSalomao/alura-challenge-backend/model"
	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReceiptService struct{}

func (rs *ReceiptService) CreateReceipt(r *model.Receipt) (uuid.UUID, error) {
	var entity *model.Receipt
	t1, t2 := utils.MonthInterval()
	tx := database.DB.Where("description = ? AND date between ? AND ?", r.Description, t1, t2).First(&entity)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		database.DB.Create(r)
	} else {
		return entity.Id, errors.New("receipt already created in current month")
	}
	return r.Id, nil
}

func (rs *ReceiptService) FindAllReceipts(r *[]model.ReceiptResponse) {
	var receipts []model.Receipt
	database.DB.Find(&receipts)
	for _, v := range receipts {
		*r = append(*r, model.ReceiptResponse{
			Id:          v.Id.String(),
			Description: v.Description,
			Value:       v.Value,
			Date:        v.Date})
	}
}

func (rs *ReceiptService) FindReceipt(r *model.ReceiptResponse, id uuid.UUID) {
	var receipt model.Receipt
	database.DB.First(&receipt, id)
	*r = model.ReceiptResponse{
		Description: receipt.Description,
		Value:       receipt.Value,
		Date:        receipt.Date,
	}
}

func (rs *ReceiptService) UpdateReceipt(r *model.Receipt, id uuid.UUID) (uuid.UUID, error) {
	var receipt model.Receipt
	tx := database.DB.First(&receipt, id)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		return id, errors.New("receipt not found")
	}
	r.CreatedAt = receipt.CreatedAt
	if receipt.Description == r.Description {
		r.Id = id
		receipt = *r
		database.DB.Save(&receipt)
	} else {
		var entity model.Receipt
		t1, t2 := utils.MonthInterval()
		tx := database.DB.Where("description = ? AND date between ? AND ?", r.Description, t1, t2).First(&entity)
		if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
			r.Id = id
			receipt = *r
			database.DB.Save(&receipt)
		} else {
			return entity.Id, errors.New("receipt already created in current month")
		}
	}
	return id, nil
}
