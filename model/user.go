package model

import (
	"strings"

	"gorm.io/gorm"
)

type User struct {
	Base
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) BeforCreate(tx *gorm.DB) (err error) {
	u.Email = strings.ToLower(u.Email)
	return nil
}
