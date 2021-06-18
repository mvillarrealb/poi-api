package common

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func RespondWithErrorDetails(w http.ResponseWriter, code int, message string, details string) {
	RespondWithJSON(w, code, map[string]string{"code": strconv.Itoa(code), "error": message, "details": details})
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"code": strconv.Itoa(code), "error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
