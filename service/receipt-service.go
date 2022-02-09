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

type receiptService struct{}

var ReceiptService = receiptService{}

func (rs *receiptService) CreateReceipt(r *model.ReceiptRequest, userId uuid.UUID) (uuid.UUID, error) {
	var entity *model.Receipt
	isTwice, entityId, err := rs.isMonthDuplicated(r.Date, r.Description, userId)
	if isTwice {
		return entityId, errors.New("receipt already created in current month")
	} else if err != nil {
		return uuid.Nil, err
	} else {
		entity = &model.Receipt{
			Description: r.Description,
			Value:       r.Value,
			Date:        r.Date,
			UserId:      userId,
		}
		database.DB.Create(entity)
	}
	return entity.Id, nil
}

func (rs *receiptService) FindAllReceipts(r *[]model.ReceiptResponse, description string, userId uuid.UUID) error {
	var receipts []model.Receipt
	if description != "" {
		database.DB.Where("user_id = ? AND description = ?", userId, strings.ToUpper(description)).Find(&receipts)
	} else {
		database.DB.Where("user_id = ?", userId).Find(&receipts)
	}
	for _, v := range receipts {
		*r = append(*r, model.ReceiptResponse{
			Id:          v.Id.String(),
			Description: v.Description,
			Value:       v.Value,
			Date:        v.Date,
			UserId:      v.UserId.String(),
		},
		)
	}
	return nil
}

func (rs *receiptService) FindReceipt(r *model.ReceiptResponse, receiptId uuid.UUID, userId uuid.UUID) error {
	var receipt model.Receipt
	tx := database.DB.Where("id = ? AND user_id = ?", receiptId, userId).First(&receipt)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		return errors.New("receipt not found")
	}
	*r = model.ReceiptResponse{
		Id:          receipt.Id.String(),
		Description: receipt.Description,
		Value:       receipt.Value,
		Date:        receipt.Date,
		UserId:      receipt.UserId.String(),
	}
	return nil
}

func (rs *receiptService) UpdateReceipt(r *model.ReceiptRequest, receiptId uuid.UUID, userId uuid.UUID) (uuid.UUID, error) {
	var receipt model.Receipt
	tx := database.DB.Where("id = ? AND user_id = ?", receiptId, userId).First(&receipt)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		return receiptId, errors.New("receipt not found")
	}
	if rs.shouldCheckReceiptInCurrentMonth(r, &receipt) {
		isTwice, entityId, err := rs.isMonthDuplicated(r.Date, r.Description, userId)
		if isTwice {
			return entityId, fmt.Errorf("receipt %s already created in current month", strings.ToUpper(r.Description))
		} else if err != nil {
			return uuid.Nil, err
		} else {
			receipt.Date = r.Date
			receipt.Description = r.Description
			receipt.Value = r.Value
			database.DB.Save(&receipt)
		}
	} else {
		receipt.Date = r.Date
		receipt.Description = r.Description
		receipt.Value = r.Value
		database.DB.Save(&receipt)
	}
	return receiptId, nil
}

func (rs *receiptService) DeleteReceipt(id uuid.UUID, userId uuid.UUID) {
	var receipt model.Receipt
	database.DB.Where("id = ? AND user_id = ?", id, userId).Delete(&receipt)
}

func (rs *receiptService) ReceiptsByPeriod(r *[]model.ReceiptResponse, year string, month string, userId uuid.UUID) error {
	var receipts []model.Receipt
	t1, t2, err := utils.MonthInterval(fmt.Sprintf("%s-%s", year, month))
	if err != nil {
		return err
	}
	database.DB.Where("user_id = ? AND date between ? AND ?", userId, t1, t2).Find(&receipts)
	for _, v := range receipts {
		*r = append(*r, model.ReceiptResponse{
			Id:          v.Id.String(),
			Description: v.Description,
			Value:       v.Value,
			Date:        v.Date,
			UserId:      v.UserId.String(),
		},
		)
	}
	return nil
}

func (rs *receiptService) TotalReceiptValueByPeriod(year string, month string, userId uuid.UUID) (float64, error) {
	var receipts []model.Receipt
	t1, t2, err := utils.MonthInterval(fmt.Sprintf("%s-%s", year, month))
	if err != nil {
		return 0, err
	}
	database.DB.Where("user_id = ? AND date between ? AND ?", userId, t1, t2).Find(&receipts)
	var total float64
	for _, v := range receipts {
		total += v.Value
	}
	return total, nil
}

func (rs *receiptService) shouldCheckReceiptInCurrentMonth(receiptRequest *model.ReceiptRequest, receipt *model.Receipt) bool {
	if (receiptRequest.Date != receipt.Date) || (receiptRequest.Description != receipt.Description) {
		return true
	} else {
		return false
	}
}

func (rs *receiptService) isMonthDuplicated(date string, description string, userId uuid.UUID) (bool, uuid.UUID, error) {
	var entity model.Receipt
	t1, t2, err := utils.MonthInterval(date)
	if err != nil {
		return false, uuid.Nil, err
	}
	tx := database.DB.Where("user_id = ? AND description = ? AND date between ? AND ?", userId, strings.ToUpper(description), t1, t2).First(&entity)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		return false, uuid.Nil, nil
	} else if tx.Error != nil {
		return true, uuid.Nil, tx.Error
	} else {
		return true, entity.Id, nil
	}
}
