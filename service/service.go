package service

import "github.com/RaphaSalomao/alura-challenge-backend/model"

func BalanceSumary(bs *model.BalanceSumaryResponse, year string, month string) error {
	totalReceipt, err := ReceiptService.TotalReceiptValueByPeriod(year, month)
	if err != nil {
		return err
	}
	totalExpense, categoryBalance, err := ExpenseService.TotalExpenseValueByPeriod(year, month)
	if err != nil {
		return err
	}
	bs.CategoryBalance = categoryBalance
	bs.TotalExpense = totalExpense
	bs.TotalReceipt = totalReceipt
	bs.MonthBalance = totalReceipt - totalExpense
	return nil
}
