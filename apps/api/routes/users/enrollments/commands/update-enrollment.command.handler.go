package commands

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/harusame0616/ijuku/apps/api/lib/response"
	"github.com/jackc/pgx/v5"
)

const (
	errorCodeEnrollmentForbidden = "ENROLLMENT_FORBIDDEN"
	errorCodeCourseNotFound      = "COURSE_NOT_FOUND"
	errorCodeTopicNotFound       = "TOPIC_NOT_FOUND"
	errorCodeNotEnrolled         = "NOT_ENROLLED"
)

type UpdateEnrollmentHandler struct {
	usecase UpdateEnrollmentUsecaseInterface
}

func NewUpdateEnrollmentHandler(usecase UpdateEnrollmentUsecaseInterface) *UpdateEnrollmentHandler {
	return &UpdateEnrollmentHandler{usecase: usecase}
}

type patchEnrollmentRequestBody struct {
	TopicId string `json:"topicId"`
	Status  string `json:"status"`
}

type patchEnrollmentResponse struct {
	TopicId string `json:"topicId"`
	Status  string `json:"status"`
}

func (h *UpdateEnrollmentHandler) PatchEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
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

	courseIdStr := r.PathValue("courseId")
	if courseIdStr == "" {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "courseId must be required")
		return
	}
	courseId, err := uuid.Parse(courseIdStr)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "courseId must be UUID format")
		return
	}

	var body patchEnrollmentRequestBody
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "body parameter is invalid json format")
		return
	}

	if body.TopicId == "" {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "topicId must be required")
		return
	}
	topicId, err := uuid.Parse(body.TopicId)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "topicId must be UUID format")
		return
	}

	if body.Status == "" {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "status must be required")
		return
	}
	status := ProgressStatus(body.Status)
	if status != ProgressStatusInProgress && status != ProgressStatusCompleted {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "status must be IN_PROGRESS or COMPLETED")
		return
	}

	result, err := h.usecase.execute(r.Context(), UpdateEnrollmentParams{
		userId:   userId,
		courseId: courseId,
		topicId:  topicId,
		status:   status,
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidProgressStatus), errors.Is(err, ErrInvalidStatusTransit):
			response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, err.Error())
		case errors.Is(err, pgx.ErrNoRows):
			response.WriteErrorResponse(w, http.StatusNotFound, errorCodeCourseNotFound, "course not found")
		case errors.Is(err, ErrTopicNotFoundInCourse):
			response.WriteErrorResponse(w, http.StatusNotFound, errorCodeTopicNotFound, "topic not found in course")
		case errors.Is(err, ErrNotEnrolled):
			response.WriteErrorResponse(w, http.StatusNotFound, errorCodeNotEnrolled, "not enrolled in this course")
		default:
			log.Printf("PatchEnrollmentHandler error: %v", err)
			response.WriteInternalServerErrorResponse(w)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(patchEnrollmentResponse{
		TopicId: result.TopicId,
		Status:  result.Status,
	})
}
