package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/RaphaSalomao/alura-challenge-backend/database"
	"github.com/RaphaSalomao/alura-challenge-backend/model"
	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReceiptService struct{}

func (rs *ReceiptService) CreateReceipt(r *model.Receipt) (uuid.UUID, error) {
	var entity *model.Receipt
	t1, t2, err := utils.MonthInterval(r.Date)
	if err != nil {
		return uuid.Nil, err
	}
	tx := database.DB.Where("description = ? AND date between ? AND ?", strings.ToUpper(r.Description), t1, t2).First(&entity)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		database.DB.Create(r)
	} else {
		return entity.Id, errors.New("receipt already created in current month")
	}
	return r.Id, nil
}

func (rs *ReceiptService) FindAllReceipts(r *[]model.ReceiptResponse, description string) {
	var receipts []model.Receipt
	if description != "" {
		database.DB.Find(&receipts, "description = ?", description)
	} else {
		database.DB.Find(&receipts)
	}
	for _, v := range receipts {
		*r = append(*r, model.ReceiptResponse{
			Id:          v.Id.String(),
			Description: v.Description,
			Value:       v.Value,
			Date:        v.Date})
	}
}

func (rs *ReceiptService) FindReceipt(r *model.ReceiptResponse, id uuid.UUID) error {
	var receipt model.Receipt
	tx := database.DB.First(&receipt, id)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		return errors.New("receipt not found")
	}
	*r = model.ReceiptResponse{
		Description: receipt.Description,
		Value:       receipt.Value,
		Date:        receipt.Date,
	}
	return nil
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
		t1, t2, err := utils.MonthInterval(r.Date)
		if err != nil {
			return uuid.Nil, err
		}
		tx := database.DB.Where("description = ? AND date between ? AND ?", strings.ToUpper(r.Description), t1, t2).First(&entity)
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

func (rs *ReceiptService) DeleteReceipt(id uuid.UUID) {
	var receipt model.Receipt
	database.DB.Delete(&receipt, id)
}

func (rs *ReceiptService) ReceiptsByPeriod(r *[]model.ReceiptResponse, year string, month string) error {
	var receipts []model.Receipt
	t1, t2, err := utils.MonthInterval(fmt.Sprintf("%s-%s", year, month))
	if err != nil {
		return err
	}
	database.DB.Find(&receipts, "date between ? AND ?", t1, t2)
	for _, v := range receipts {
		*r = append(*r, model.ReceiptResponse{
			Id:          v.Id.String(),
			Description: v.Description,
			Value:       v.Value,
			Date:        v.Date,
			Category:    v.Category,
		})
	}
	return nil
}

func (rs *ReceiptService) TotalReceiptValueByPeriod(year, month string) (float64, error) {
	var receipts []model.Receipt
	t1, t2, err := utils.MonthInterval(fmt.Sprintf("%s-%s", year, month))
	if err != nil {
		return 0, err
	}
	database.DB.Find(&receipts, "date between ? AND ?", t1, t2)
	var total float64
	for _, v := range receipts {
		total += v.Value
	}
	return total, nil
}

