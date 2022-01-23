package service

import (
	"errors"

	"github.com/RaphaSalomao/alura-challenge-backend/database"
	"github.com/RaphaSalomao/alura-challenge-backend/model"
	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpenseService struct{}

func (rs *ExpenseService) CreateExpense(e *model.Expense) (uuid.UUID, error) {
	var entity *model.Expense
	t1, t2 := utils.MonthInterval()
	tx := database.DB.Where("description = ? AND date between ? AND ?", e.Description, t1, t2).First(&entity)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		database.DB.Create(e)
	} else {
		return entity.Id, errors.New("expense already created in current month")
	}
	return e.Id, nil
}

func (rs *ExpenseService) FindAllExpenses(e *[]model.ExpenseResponse) {
	var expenses []model.Expense
	database.DB.Find(&expenses)
	for _, v := range expenses {
		*e = append(*e, model.ExpenseResponse{
			Id:          v.Id.String(),
			Description: v.Description,
			Value:       v.Value,
			Date:        v.Date})
	}
}

func (rs *ExpenseService) FindExpense(e *model.ExpenseResponse, id uuid.UUID) error {
	var expense model.Expense
	tx := database.DB.First(&expense, id)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		return errors.New("expense not found")
	}
	*e = model.ExpenseResponse{
		Description: expense.Description,
		Value:       expense.Value,
		Date:        expense.Date,
	}
	return nil
}

func (rs *ExpenseService) UpdateExpense(e *model.Expense, id uuid.UUID) (uuid.UUID, error) {
	var expense model.Expense
	tx := database.DB.First(&expense, id)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		return id, errors.New("expense not found")
	}
	e.CreatedAt = expense.CreatedAt
	if expense.Description == e.Description {
		e.Id = id
		expense = *e
		database.DB.Save(&expense)
	} else {
		var entity model.Expense
		t1, t2 := utils.MonthInterval()
		tx := database.DB.Where("description = ? AND date between ? AND ?", e.Description, t1, t2).First(&entity)
		if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
			e.Id = id
			expense = *e
			database.DB.Save(&expense)
		} else {
			return entity.Id, errors.New("expense already created in current month")
		}
	}
	return id, nil
}

func (rs *ExpenseService) DeleteExpense(id uuid.UUID) {
	var expense model.Expense
	database.DB.Delete(&expense, id)
}
