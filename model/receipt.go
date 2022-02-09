package model

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Receipt struct {
	Base
	Description string    `json:"description,omitempty"`
	Value       float64   `json:"value,omitempty"`
	Date        string    `json:"date,omitempty"`
	UserId      uuid.UUID `json:"userId,omitempty"`
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
	UserId      string  `json:"userId,omitempty"`
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
