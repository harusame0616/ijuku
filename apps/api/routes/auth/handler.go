package auth

import (
	"net/http"

	"github.com/harusame0616/ijuku/apps/api/lib/auth"
	"github.com/harusame0616/ijuku/apps/api/lib/response"
)

type Handler struct {
	verifier *auth.Verifier
}

func NewHandler(verifier *auth.Verifier) *Handler {
	return &Handler{verifier: verifier}
}

// TODO テスト用のエンドポイント。削除する
func (h *Handler) CheckHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.ExtractBearerToken(r)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	if err := h.verifier.Verify(token); err != nil {
		response.WriteErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	w.WriteHeader(http.StatusOK)
}
