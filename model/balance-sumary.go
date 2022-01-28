package model

import "github.com/RaphaSalomao/alura-challenge-backend/model/enum"

type BalanceSumaryResponse struct {
	TotalReceipt    float64                   `json:"totalReceipt"`
	TotalExpense    float64                   `json:"totalExpense"`
	MonthBalance    float64                   `json:"monthBalance"`
	CategoryBalance map[enum.Category]float64 `json:"categoryBalance"`
}
