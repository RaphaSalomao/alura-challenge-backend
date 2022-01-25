package model

import (
	"strings"

	"github.com/RaphaSalomao/alura-challenge-backend/model/enum"
	"gorm.io/gorm"
)

type Expense struct {
	Base
	Description string        `json:"description,omitempty"`
	Value       float64       `json:"value,omitempty"`
	Date        string        `json:"date,omitempty"`
	Category    enum.Category `json:"category,omitempty"`
}

type ExpenseRequest struct {
	Description string        `json:"description,omitempty"`
	Value       float64       `json:"value,omitempty"`
	Date        string        `json:"date,omitempty"`
	Category    enum.Category `json:"category,omitempty"`
}

type ExpenseResponse struct {
	Id          string        `json:"id,omitempty"`
	Description string        `json:"description,omitempty"`
	Value       float64       `json:"value,omitempty"`
	Date        string        `json:"date,omitempty"`
	Category    enum.Category `json:"category,omitempty"`
}

func (e *Expense) BeforeCreate(tx *gorm.DB) (err error) {
	e.Base.BeforeCreate(tx)
	if e.Category == enum.CategoryUndefined {
		e.Category = enum.CategoryOthers
	}
	e.Description = strings.ToUpper(e.Description)
	return
}

func (e *Expense) BeforeSave(tx *gorm.DB) (err error) {
	e.Base.BeforeSave(tx)
	e.Description = strings.ToUpper(e.Description)
	return
}
