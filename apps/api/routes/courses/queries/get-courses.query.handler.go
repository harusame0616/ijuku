package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/response"
	"github.com/jackc/pgx/v5/pgtype"
)

const pageSize = 200

type GetCoursesQuery interface {
	GetCourses(ctx context.Context, arg db.GetCoursesParams) ([]db.GetCoursesRow, error)
}

type Handlers struct {
	query GetCoursesQuery
}

type GetCoursesResult struct {
	Courses []db.GetCoursesRow `json:"courses"`
	Cursor  *string            `json:"cursor"`
}

func NewCoursesHandlers(qs GetCoursesQuery) *Handlers {
	return &Handlers{query: qs}
}

func parseGetCoursesQuery(r *http.Request) (keyword string, cursorUuid pgtype.UUID, err error) {
	keyword = r.URL.Query().Get("keyword")
	cursor := r.URL.Query().Get("cursor")

	if len([]rune(keyword)) > 40 {
		return "", pgtype.UUID{}, fmt.Errorf("keyword must be 40 characters or less")
	}
	if cursor != "" {
		if err := cursorUuid.Scan(cursor); err != nil {
			return "", pgtype.UUID{}, fmt.Errorf("invalid cursor")
		}
	}

	return keyword, cursorUuid, nil
}

func (handler *Handlers) GetCoursesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	keyword, cursorUuid, err := parseGetCoursesQuery(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": response.InputValidationError, "message": err.Error()})
		return
	}

	rawCourses, err := handler.query.GetCourses(r.Context(), db.GetCoursesParams{Cursor: cursorUuid, Keyword: keyword, Size: pageSize + 1})
	if err != nil {
		log.Printf("GetCourses error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": "INTERNAL_SERVER_ERROR", "message": "internal server error"})
		return
	}

	courses := rawCourses[0:min(len(rawCourses), pageSize)]

	var nextCursor *string
	if len(rawCourses) > pageSize {
		s := rawCourses[pageSize-1].CourseId.String()
		nextCursor = &s
	}

	_ = json.NewEncoder(w).Encode(GetCoursesResult{Courses: courses, Cursor: nextCursor})
}
