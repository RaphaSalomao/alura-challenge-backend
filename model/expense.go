package model

type Expense struct {
	Base
	Description string  `json:"description,omitempty"`
	Value       float64 `json:"value,omitempty"`
	Date        string  `json:"date,omitempty"`
}

type ExpenseResponse struct {
	Id          string  `json:"id,omitempty"`
	Description string  `json:"description,omitempty"`
	Value       float64 `json:"value,omitempty"`
	Date        string  `json:"date,omitempty"`
}
