package commands

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/harusame0616/ijuku/apps/api/lib/response"
	"github.com/jackc/pgx/v5"
)

const (
	errorCodeAlreadyEnrolled = "ALREADY_ENROLLED"
)

type EnrollHandler struct {
	usecase EnrollUsecaseInterface
}

func NewEnrollHandler(usecase EnrollUsecaseInterface) *EnrollHandler {
	return &EnrollHandler{usecase: usecase}
}

type postEnrollmentRequestBody struct {
	CourseId string `json:"courseId"`
}

type postEnrollmentResponse struct {
	CourseId   string `json:"courseId"`
	EnrolledAt string `json:"enrolledAt"`
}

func (h *EnrollHandler) PostEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIdStr := r.PathValue("userID")
	if userIdStr == "" {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "userID must be required")
		return
	}
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "userID must be UUID format")
		return
	}

	var body postEnrollmentRequestBody
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "body parameter is invalid json format")
		return
	}

	if body.CourseId == "" {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "courseId must be required")
		return
	}
	courseId, err := uuid.Parse(body.CourseId)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "courseId must be UUID format")
		return
	}

	result, err := h.usecase.execute(r.Context(), EnrollParams{
		userId:   userId,
		courseId: courseId,
	})
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			response.WriteErrorResponse(w, http.StatusNotFound, errorCodeCourseNotFound, "course not found")
		case errors.Is(err, ErrEnrollmentNotAllowed):
			response.WriteErrorResponse(w, http.StatusForbidden, errorCodeEnrollmentForbidden, "this course is not enrollable")
		case errors.Is(err, ErrAlreadyEnrolled):
			response.WriteErrorResponse(w, http.StatusConflict, errorCodeAlreadyEnrolled, "already enrolled in this course")
		default:
			log.Printf("PostEnrollmentHandler error: %v", err)
			response.WriteInternalServerErrorResponse(w)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(postEnrollmentResponse{
		CourseId:   result.CourseId,
		EnrolledAt: result.EnrolledAt.Format(time.RFC3339),
	})
}
