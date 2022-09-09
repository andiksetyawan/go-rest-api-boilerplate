package httputil

import (
	"encoding/json"
	"net/http"

	"go-rest-api-boilerplate/internal/model/reqres"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, reqres.ApiResponse{
		Error:   true,
		Message: message,
		Data:    nil,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
