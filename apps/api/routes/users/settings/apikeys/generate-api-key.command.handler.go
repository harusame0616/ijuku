package apikeys

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/harusame0616/ijuku/apps/api/lib/response"
	"github.com/harusame0616/ijuku/apps/api/lib/txrunner"
)

type generateApiKey struct {
	usecase generateApiKeyUsecase
}

func NewGenerateApiKeyHandler(usecase generateApiKeyUsecase) generateApiKey {
	return generateApiKey{
		usecase: usecase,
	}
}

func (generateApiKey generateApiKey) GenerateApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO 認証

	userID := r.PathValue("userID")
	if userID == "" {
		// パスで UserID がマッピングされているのでこのパスは通過しない前提
		// ここを通過する場合はパスのマッピングが間違っているなどの不具合の可能性
		response.WriteInternalServerErrorResponse(w)
		return
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "User ID must be valid UUID")
		return
	}

	var bodyParams struct {
		ExpiredAt *time.Time `json:"expiredAt"`
	}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&bodyParams); err != nil {
		var timeParseErr *time.ParseError

		switch {
		case errors.As(err, &timeParseErr):
			response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "expiredAt must be ISO 8601 format")
		default:
			response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "Body must be valid JSON")
		}
		return
	}

	result, err := generateApiKey.usecase.Execute(r.Context(), parsedUserID, bodyParams.ExpiredAt)

	if err != nil {
		switch {
		case errors.Is(err, ErrApiKeyCountExceedsLimit):
			response.WriteErrorResponse(w, http.StatusConflict, "APIKEY_QUOTA_EXCEEDS_LIMIT", "Api key quota exceeds limit. Api key quota limit is "+strconv.Itoa(apiKeyMaxCount))
		case errors.Is(err, txrunner.ErrLockTimeout):
			response.WriteErrorResponse(w, http.StatusServiceUnavailable, "APIKEY_LOCK_TIMEOUT", "Api key generation is temporarily unavailable. Please try again later.")
		default:
			fmt.Printf("err %v", err)
			response.WriteInternalServerErrorResponse(w)
		}

		return
	}

	json.NewEncoder(w).Encode(result)
}
