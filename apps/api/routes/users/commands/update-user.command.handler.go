package commands

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	libauth "github.com/harusame0616/ijuku/apps/api/lib/auth"
	"github.com/harusame0616/ijuku/apps/api/lib/response"
)

type updateUserUsecase interface {
	Execute(ctx context.Context, userID uuid.UUID, nickname, introduce string) error
}

type UpdateUserHandler struct {
	usecase  updateUserUsecase
	verifier *libauth.Verifier
}

func NewUpdateUserHandler(usecase updateUserUsecase, verifier *libauth.Verifier) *UpdateUserHandler {
	return &UpdateUserHandler{usecase: usecase, verifier: verifier}
}

func (h *UpdateUserHandler) PatchUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, err := libauth.ExtractBearerToken(r)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	jwtUserID, err := h.verifier.GetUserID(token)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	pathUserID := r.PathValue("userID")
	if jwtUserID != pathUserID {
		response.WriteErrorResponse(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}

	userID, err := uuid.Parse(pathUserID)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "userID must be a valid UUID")
		return
	}

	var body struct {
		Nickname  string `json:"nickname"`
		Introduce string `json:"introduce"`
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "body must be valid JSON")
		return
	}

	if err := h.usecase.Execute(r.Context(), userID, body.Nickname, body.Introduce); err != nil {
		switch {
		case errors.Is(err, ErrValidation):
			response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, err.Error())
		case errors.Is(err, ErrUserNotFound):
			response.WriteErrorResponse(w, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
		default:
			response.WriteInternalServerErrorResponse(w)
		}
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"nickname":  body.Nickname,
		"introduce": body.Introduce,
	})
}
