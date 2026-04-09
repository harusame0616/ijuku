package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

func WriteErrorResponse(w http.ResponseWriter, status int, errorCode, message string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{ErrorCode: errorCode, Message: message})
}

func WriteInternalServerErrorResponse(w http.ResponseWriter) {
	WriteErrorResponse(w, http.StatusInternalServerError, ServerInternalError, "An unexpected error occurred. Please try again later.")

}
