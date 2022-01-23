package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

func HandleResponse(w http.ResponseWriter, status int, i interface{}) {
	w.WriteHeader(status)
	if i != nil {
		json.NewEncoder(w).Encode(&i)
	}
}

func MonthInterval() (firstDay, lastDay time.Time) {
	firstDay = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	lastDay = time.Date(time.Now().Year(), time.Now().Month()+1, 1, 0, 0, 0, -1, time.Local)
	return
}
