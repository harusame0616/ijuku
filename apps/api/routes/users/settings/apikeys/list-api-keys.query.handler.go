package apikeys

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/harusame0616/ijuku/apps/api/internal/db"
	libauth "github.com/harusame0616/ijuku/apps/api/lib/auth"
	"github.com/harusame0616/ijuku/apps/api/lib/response"
	"github.com/jackc/pgx/v5/pgtype"
)

type listApiKeysQuery interface {
	ListApiKeysByUserID(ctx context.Context, userid pgtype.UUID) ([]db.ListApiKeysByUserIDRow, error)
}

type ListApiKeysHandler struct {
	query    listApiKeysQuery
	verifier *libauth.Verifier
}

func NewListApiKeysHandler(q listApiKeysQuery, verifier *libauth.Verifier) *ListApiKeysHandler {
	return &ListApiKeysHandler{query: q, verifier: verifier}
}

type ApiKeyListItem struct {
	ApiKeyID  string  `json:"apiKeyId"`
	Suffix    string  `json:"suffix"`
	CreatedAt string  `json:"createdAt"`
	ExpiredAt *string `json:"expiredAt"`
}

type ListApiKeysResponse struct {
	ApiKeys []ApiKeyListItem `json:"apiKeys"`
}

func (h *ListApiKeysHandler) ListApiKeysHandler(w http.ResponseWriter, r *http.Request) {
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

	var userID pgtype.UUID
	if err := userID.Scan(pathUserID); err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "userID must be a valid UUID")
		return
	}

	rows, err := h.query.ListApiKeysByUserID(r.Context(), userID)
	if err != nil {
		response.WriteInternalServerErrorResponse(w)
		return
	}

	items := make([]ApiKeyListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, ApiKeyListItem{
			ApiKeyID:  uuid.UUID(row.ApikeyID.Bytes).String(),
			Suffix:    row.PlainSuffix,
			CreatedAt: row.CreatedAt.Time.UTC().Format(time.RFC3339),
			ExpiredAt: formatExpiredAt(row.ExpiredAt),
		})
	}

	_ = json.NewEncoder(w).Encode(ListApiKeysResponse{ApiKeys: items})
}

func formatExpiredAt(t pgtype.Timestamptz) *string {
	if !t.Valid || t.InfinityModifier != pgtype.Finite {
		return nil
	}
	formatted := t.Time.UTC().Format(time.RFC3339)
	return &formatted
}
