package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/RaphaSalomao/alura-challenge-backend/database"
	"github.com/RaphaSalomao/alura-challenge-backend/model"
	"github.com/RaphaSalomao/alura-challenge-backend/model/enum"
	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type expenseService struct{}

var ExpenseService = expenseService{}

func (es *expenseService) CreateExpense(e *model.ExpenseRequest) (uuid.UUID, error) {
	var entity *model.Expense
	t1, t2, err := utils.MonthInterval(e.Date)
	if err != nil {
		return uuid.Nil, err
	}
	tx := database.DB.Where("description = ? AND date between ? AND ?", strings.ToUpper(e.Description), t1, t2).First(&entity)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		entity = &model.Expense{
			Description: e.Description,
			Value:       e.Value,
			Date:        e.Date,
			Category:    e.Category,
		}
		database.DB.Create(entity)
	} else {
		return entity.Id, errors.New("expense already created in current month")
	}
	return entity.Id, nil
}

func (es *expenseService) FindAllExpenses(e *[]model.ExpenseResponse, description string) {
	var expenses []model.Expense
	if description != "" {
		database.DB.Find(&expenses, "description = ?", description)
	} else {
		database.DB.Find(&expenses)
	}
	for _, v := range expenses {
		*e = append(*e, model.ExpenseResponse{
			Id:          v.Id.String(),
			Description: v.Description,
			Value:       v.Value,
			Date:        v.Date})
	}
}

func (es *expenseService) FindExpense(e *model.ExpenseResponse, id uuid.UUID) error {
	var expense model.Expense
	tx := database.DB.First(&expense, id)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		return errors.New("expense not found")
	}
	*e = model.ExpenseResponse{
		Description: expense.Description,
		Value:       expense.Value,
		Date:        expense.Date,
		Category:    expense.Category,
	}
	return nil
}

func (es *expenseService) UpdateExpense(e *model.ExpenseRequest, id uuid.UUID) (uuid.UUID, error) {
	var expense model.Expense
	tx := database.DB.First(&expense, id)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		return id, errors.New("expense not found")
	}
	if es.shouldCheckExpenseInCurrentMonth(e, &expense) {
		var entity model.Expense
		t1, t2, err := utils.MonthInterval(e.Date)
		if err != nil {
			return uuid.Nil, err
		}
		tx := database.DB.Where("description = ? AND date between ? AND ?", strings.ToUpper(e.Description), t1, t2).First(&entity)
		if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
			expense.Category = e.Category
			expense.Date = e.Date
			expense.Description = e.Description
			expense.Value = e.Value
			database.DB.Save(&expense)
		} else {
			return entity.Id, fmt.Errorf("expense %s already created in current month", entity.Description)
		}
	} else {
		expense.Category = e.Category
		expense.Date = e.Date
		expense.Description = e.Description
		expense.Value = e.Value
		database.DB.Save(&expense)
	}
	return id, nil
}

func (es *expenseService) DeleteExpense(id uuid.UUID) {
	var expense model.Expense
	database.DB.Delete(&expense, id)
}

func (es *expenseService) ExpensesByPeriod(e *[]model.ExpenseResponse, year string, month string) error {
	var expenses []model.Expense
	t1, t2, err := utils.MonthInterval(fmt.Sprintf("%s-%s", year, month))
	if err != nil {
		return err
	}
	database.DB.Find(&expenses, "date between ? AND ?", t1, t2)
	for _, v := range expenses {
		*e = append(*e, model.ExpenseResponse{
			Id:          v.Id.String(),
			Description: v.Description,
			Value:       v.Value,
			Date:        v.Date,
			Category:    v.Category,
		})
	}
	return nil
}

func (es *expenseService) TotalExpenseValueByPeriod(year, month string) (total float64, categoriesBalance map[enum.Category]float64, err error) {
	var expenses []model.Expense
	categoriesBalance = map[enum.Category]float64{}
	t1, t2, err := utils.MonthInterval(fmt.Sprintf("%s-%s", year, month))
	if err != nil {
		return 0, nil, err
	}
	database.DB.Find(&expenses, "date between ? AND ?", t1, t2)
	for _, v := range expenses {
		categoriesBalance[v.Category] += v.Value
		total += v.Value
	}
	return
}

func (es *expenseService) shouldCheckExpenseInCurrentMonth(expenseRequest *model.ExpenseRequest, expense *model.Expense) bool {
	if expenseRequest.Date != expense.Date && expenseRequest.Description != expense.Description {
		return true
	} else {
		return false
	}
}
