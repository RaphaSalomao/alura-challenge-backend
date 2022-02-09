package service

import (
	"errors"
	"fmt"

	"github.com/RaphaSalomao/alura-challenge-backend/database"
	"github.com/RaphaSalomao/alura-challenge-backend/model"
	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUser(u *model.UserRequest) error {
	hashPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	user := model.User{
		Email:    u.Email,
		Password: hashPassword,
	}
	database.DB.Create(&user)
	return nil
}

func BalanceSumary(bs *model.BalanceSumaryResponse, year string, month string, userId uuid.UUID) error {
	totalReceipt, err := ReceiptService.TotalReceiptValueByPeriod(year, month, userId)
	if err != nil {
		return err
	}
	totalExpense, categoryBalance, err := ExpenseService.TotalExpenseValueByPeriod(year, month, userId)
	if err != nil {
		return err
	}
	bs.CategoryBalance = categoryBalance
	bs.TotalExpense = totalExpense
	bs.TotalReceipt = totalReceipt
	bs.MonthBalance = totalReceipt - totalExpense
	return nil
}

func Authenticate(u *model.UserRequest) (string, error) {
	user := model.User{}
	tx := database.DB.Where("email = ?", u.Email).First(&user)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("user %s not found", u.Email)
		} else {
			return "", tx.Error
		}
	}
	if utils.ValidadeHashAndPassword(u.Password, user.Password) {
		return utils.GenerateJWT(user.Email, user.Id.String())
	} else {
		return "", errors.New("invalid user/password")
	}
}
