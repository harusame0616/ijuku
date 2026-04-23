package queries

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/response"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type GetUserQuery interface {
	GetUser(ctx context.Context, userid pgtype.UUID) (db.GetUserRow, error)
}

type GetUserHandler struct {
	query GetUserQuery
}

func NewGetUserHandler(q GetUserQuery) *GetUserHandler {
	return &GetUserHandler{query: q}
}

type GetUserResponse struct {
	Nickname  string `json:"nickname"`
	Introduce string `json:"introduce"`
}

func (h *GetUserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("userID")); err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "userID must be a valid UUID")
		return
	}

	user, err := h.query.GetUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteErrorResponse(w, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
			return
		}
		log.Printf("GetUser error: %v", err)
		response.WriteInternalServerErrorResponse(w)
		return
	}

	_ = json.NewEncoder(w).Encode(GetUserResponse{
		Nickname:  user.Nickname,
		Introduce: user.Introduce,
	})
}
