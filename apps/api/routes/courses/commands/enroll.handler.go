package commands

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/harusame0616/ijuku/apps/api/lib/uuid"
	"github.com/harusame0616/ijuku/apps/api/lib/validation"
)

type Handler struct {
	usecase EnrollCourseUsecaseInterface
}

func NewHandler(usecase EnrollCourseUsecaseInterface) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) PostEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	courseId := r.PathValue("courseId")

	if courseId == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": validation.InputValidationError, "message": "courseId must be required"})
		return
	}

	if !uuid.IsValidUuid(courseId) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": validation.InputValidationError, "message": "courseId must be UUID format"})
		return
	}

	// TODO: 本来は API キーから userId を解決するべきだが、認証機能未実装のため暫定的に body から取得している
	var enrollmentBodyParams struct {
		UserId        string `json:"userId"`
		SectionNumber *int   `json:"sectionNumber"`
		TopicNumber   *int   `json:"topicNumber"`
	}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&enrollmentBodyParams); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": validation.InputValidationError, "message": "body parameter is invalid json format"})
		return
	}

	if enrollmentBodyParams.UserId == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": validation.InputValidationError, "message": "userId must be required"})
		return
	}

	if !uuid.IsValidUuid(enrollmentBodyParams.UserId) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": validation.InputValidationError, "message": "userId must be UUID format"})
		return
	}

	topicId, err := h.usecase.execute(r.Context(), EnrollCourseUsecaseParams{
		userId:        enrollmentBodyParams.UserId,
		courseId:      courseId,
		sectionNumber: enrollmentBodyParams.SectionNumber,
		topicNumber:   enrollmentBodyParams.TopicNumber,
	})

	if err == ErrTopicNumberRequireSectionNumber {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": validation.InputValidationError, "message": "topic number require section number"})
		return
	}

	if err == ErrEnrollmentNumberIsNotFound {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": validation.InputValidationError, "message": "enrollment number is not found"})
		return
	}

	if err == ErrEnrollmentNotAllowed {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": "FORBIDDEN", "message": "enrollment not allowed"})
		return
	}

	if err != nil {
		log.Printf("PostEnrollmentHandler error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": "INTERNAL_SERVER_ERROR", "message": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"topicId": topicId})
}
