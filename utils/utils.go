package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RaphaSalomao/alura-challenge-backend/database"
	"github.com/RaphaSalomao/alura-challenge-backend/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	key = []byte("0sQPpmdBGjDHKXb18jNh")
)

func HandleResponse(w http.ResponseWriter, status int, i interface{}) {
	w.WriteHeader(status)
	if i != nil {
		json.NewEncoder(w).Encode(&i)
	}
}

func MonthInterval(date string) (firstDay, lastDay time.Time, err error) {
	year, month, err := GetYearMonthFromDateString(date)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	firstDay = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	lastDay = time.Date(year, time.Month(month+1), 1, 0, 0, 0, -1, time.Local)
	return
}

func GetYearMonthFromDateString(date string) (int, int, error) {
	splitDate := strings.Split(date, "-")
	year, err := strconv.Atoi(splitDate[0])
	if err != nil {
		return -1, -1, errors.New("unable to parse date")
	}
	month, err := strconv.Atoi(splitDate[1])
	if err != nil {
		return -1, -1, errors.New("unable to parse date")
	}
	return year, month, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidadeHashAndPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(email string, id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"id":    id,
		"exp":   time.Now().Add(time.Minute * 3).Unix(),
	})
	tknString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tknString, nil
}

func ParseToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, KeyFunc)
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}

func KeyFunc(t *jwt.Token) (interface{}, error) {
	id := t.Claims.(jwt.MapClaims)["id"]
	var user model.User
	tx := database.DB.First(&user, uuid.MustParse(id.(string)))
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid token")
		} else {
			return nil, tx.Error
		}
	}
	return []byte(key), nil
}
