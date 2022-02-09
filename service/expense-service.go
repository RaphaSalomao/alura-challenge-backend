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

func (es *expenseService) CreateExpense(e *model.ExpenseRequest, userId uuid.UUID) (uuid.UUID, error) {
	var entity *model.Expense
	isTwice, entityId, err := es.isMonthDuplicated(e.Date, e.Description, userId)
	if isTwice {
		return entityId, errors.New("expense already created in current month")
	} else if err != nil {
		return uuid.Nil, err
	} else {
		entity = &model.Expense{
			Description: e.Description,
			Value:       e.Value,
			Date:        e.Date,
			Category:    e.Category,
			UserId:      userId,
		}
		database.DB.Create(entity)
	}
	return entity.Id, nil
}

func (es *expenseService) FindAllExpenses(e *[]model.ExpenseResponse, description string, userId uuid.UUID) error {
	var expenses []model.Expense
	if description != "" {
		database.DB.Where("user_id = ? AND description = ?", userId, strings.ToUpper(description)).Find(&expenses)
	} else {
		database.DB.Where("user_id = ?", userId).Find(&expenses)
	}
	for _, v := range expenses {
		*e = append(*e, model.ExpenseResponse{
			Id:          v.Id.String(),
			Description: v.Description,
			Value:       v.Value,
			Date:        v.Date,
			Category:    v.Category,
			UserId:      v.UserId.String(),
		},
		)
	}
	return nil
}

func (es *expenseService) FindExpense(e *model.ExpenseResponse, id uuid.UUID, userId uuid.UUID) error {
	var expense model.Expense
	tx := database.DB.Where("id = ? AND user_id = ?", id, userId).First(&expense)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		return errors.New("expense not found")
	}
	*e = model.ExpenseResponse{
		Id:          expense.Id.String(),
		Description: expense.Description,
		Value:       expense.Value,
		Date:        expense.Date,
		Category:    expense.Category,
		UserId:      expense.UserId.String(),
	}
	return nil
}

func (es *expenseService) UpdateExpense(e *model.ExpenseRequest, id uuid.UUID, userId uuid.UUID) (uuid.UUID, error) {
	var expense model.Expense
	tx := database.DB.First(&expense, id)
	if tx.Error != nil && tx.Error == gorm.ErrRecordNotFound {
		return id, errors.New("expense not found")
	}
	if es.shouldCheckExpenseInCurrentMonth(e, &expense) {
		isTwice, entityId, err := es.isMonthDuplicated(e.Date, e.Description, userId)
		if isTwice {
			return entityId, fmt.Errorf("expense %s already created in current month", strings.ToUpper(e.Description))
		} else if err != nil {
			return uuid.Nil, err
		} else {
			expense.Category = e.Category
			expense.Date = e.Date
			expense.Description = e.Description
			expense.Value = e.Value
			database.DB.Save(&expense)
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

func (es *expenseService) DeleteExpense(id uuid.UUID, userId uuid.UUID) {
	var expense model.Expense
	database.DB.Where("id = ? AND user_id = ?", id, userId).Delete(&expense)
}

func (es *expenseService) ExpensesByPeriod(e *[]model.ExpenseResponse, year string, month string, userId uuid.UUID) error {
	var expenses []model.Expense
	t1, t2, err := utils.MonthInterval(fmt.Sprintf("%s-%s", year, month))
	if err != nil {
		return err
	}
	database.DB.Where("user_id = ? AND date between ? AND ?", userId, t1, t2).Find(&expenses)
	for _, v := range expenses {
		*e = append(*e, model.ExpenseResponse{
			Id:          v.Id.String(),
			Description: v.Description,
			Value:       v.Value,
			Date:        v.Date,
			Category:    v.Category,
			UserId:      v.UserId.String(),
		})
	}
	return nil
}

func (es *expenseService) TotalExpenseValueByPeriod(year string, month string, userId uuid.UUID) (total float64, categoriesBalance map[enum.Category]float64, err error) {
	var expenses []model.Expense
	categoriesBalance = map[enum.Category]float64{}
	t1, t2, err := utils.MonthInterval(fmt.Sprintf("%s-%s", year, month))
	if err != nil {
		return 0, nil, err
	}
	database.DB.Where("user_id = ? AND date between ? AND ?", userId, t1, t2).Find(&expenses)
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

func (ex *expenseService) isMonthDuplicated(date string, description string, userId uuid.UUID) (bool, uuid.UUID, error) {
	var entity model.Expense
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
