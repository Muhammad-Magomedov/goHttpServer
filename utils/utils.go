package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Utils struct {
}

func New() Utils {
	return Utils{}
}

const dateFormat = "2006-01-02"

func (u *Utils) ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(dateFormat, dateStr)
}

func (u *Utils) ParseUserID(userIDStr string) (int, error) {
	return strconv.Atoi(userIDStr)
}

func (u *Utils) WriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}

func (u *Utils) WriteError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), code)
}
