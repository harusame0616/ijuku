package contacts

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/response"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	nameMaxLength    = 100
	emailMaxLength   = 255
	phoneMaxLength   = 20
	contentMaxLength = 2000
)

type insertContactCommand interface {
	InsertContact(ctx context.Context, arg db.InsertContactParams) error
}

type PostContactHandler struct {
	command insertContactCommand
}

func NewPostContactHandler(command insertContactCommand) *PostContactHandler {
	return &PostContactHandler{command: command}
}

type postContactRequest struct {
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Phone   *string `json:"phone"`
	Content string  `json:"content"`
}

func (h *PostContactHandler) PostContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body postContactRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "body must be valid JSON")
		return
	}

	name := strings.TrimSpace(body.Name)
	email := strings.TrimSpace(body.Email)
	content := strings.TrimSpace(body.Content)
	var phone string
	if body.Phone != nil {
		phone = strings.TrimSpace(*body.Phone)
	}

	if name == "" {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "name is required")
		return
	}
	if len([]rune(name)) > nameMaxLength {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "name is too long")
		return
	}
	if email == "" {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "email is required")
		return
	}
	if len(email) > emailMaxLength {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "email is too long")
		return
	}
	if _, err := mail.ParseAddress(email); err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "email format is invalid")
		return
	}
	if len([]rune(phone)) > phoneMaxLength {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "phone is too long")
		return
	}
	if content == "" {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "content is required")
		return
	}
	if len([]rune(content)) > contentMaxLength {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "content is too long")
		return
	}

	contactID, err := uuid.NewRandom()
	if err != nil {
		log.Printf("contact uuid generation error: %v", err)
		response.WriteInternalServerErrorResponse(w)
		return
	}

	var contactIDPg pgtype.UUID
	if err := contactIDPg.Scan(contactID.String()); err != nil {
		log.Printf("contact uuid scan error: %v", err)
		response.WriteInternalServerErrorResponse(w)
		return
	}

	var phonePg pgtype.Text
	if phone != "" {
		phonePg = pgtype.Text{String: phone, Valid: true}
	}

	if err := h.command.InsertContact(r.Context(), db.InsertContactParams{
		ContactID: contactIDPg,
		Name:      name,
		Email:     email,
		Phone:     phonePg,
		Content:   content,
		IpAddress: clientIP(r),
		UserAgent: r.UserAgent(),
	}); err != nil {
		log.Printf("InsertContact error: %v", err)
		response.WriteInternalServerErrorResponse(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func clientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		if first := strings.TrimSpace(strings.SplitN(forwarded, ",", 2)[0]); first != "" {
			return first
		}
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}
