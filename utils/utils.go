package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
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
