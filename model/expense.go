package model

type Expense struct {
	Base
	Description string  `json:"description,omitempty"`
	Value       float64 `json:"value,omitempty"`
	Date        string  `json:"date,omitempty"`
}
