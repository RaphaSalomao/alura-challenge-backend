package model

import (
	"strings"

	"gorm.io/gorm"
)

type Receipt struct {
	Base
	Description string  `json:"description,omitempty"`
	Value       float64 `json:"value,omitempty"`
	Date        string  `json:"date,omitempty"`
}

type ReceiptRequest struct {
	Description string  `json:"description,omitempty"`
	Value       float64 `json:"value,omitempty"`
	Date        string  `json:"date,omitempty"`
}

type ReceiptResponse struct {
	Id          string  `json:"id,omitempty"`
	Description string  `json:"description,omitempty"`
	Value       float64 `json:"value,omitempty"`
	Date        string  `json:"date,omitempty"`
}

func (r *Receipt) BeforeCreate(tx *gorm.DB) (err error) {
	r.Base.BeforeCreate(tx)
	r.Description = strings.ToUpper(r.Description)
	return
}

func (r *Receipt) BeforeSave(tx *gorm.DB) (err error) {
	r.Base.BeforeSave(tx)
	r.Description = strings.ToUpper(r.Description)
	return
}
