package queries

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/harusame0616/ijuku/apps/api/lib/uuid"
	"github.com/harusame0616/ijuku/apps/api/lib/validation"
)

type Handlers struct {
	query CourseQueryService
}

func NewCoursesHandlers(qs CourseQueryService) *Handlers {
	return &Handlers{query: qs}
}

func parseGetCoursesQuery(r *http.Request) (keyword, cursor string, err error) {
	keyword = r.URL.Query().Get("keyword")
	cursor = r.URL.Query().Get("cursor")

	if len([]rune(keyword)) > 40 {
		return "", "", fmt.Errorf("keyword must be 40 characters or less")
	}
	if cursor != "" && !uuid.IsValidUuid(cursor) {
		return "", "", fmt.Errorf("invalid cursor")
	}

	return keyword, cursor, nil
}

func (handler *Handlers) GetCoursesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	keyword, cursor, err := parseGetCoursesQuery(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": validation.InputValidationError, "message": err.Error()})
		return
	}

	courses, err := handler.query.FindCourses(r.Context(), keyword, cursor)
	if err != nil {
		log.Printf("FindCourses error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": "INTERNAL_SERVER_ERROR", "message": "internal server error"})
		return
	}

	_ = json.NewEncoder(w).Encode(courses)
}
