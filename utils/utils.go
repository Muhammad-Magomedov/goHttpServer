package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}

func WriteError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), code)
}
