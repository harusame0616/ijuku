package queries

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/response"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type TopicDetailHandler struct {
	query GetTopicDetailQuery
}

type GetTopicDetailQuery interface {
	GetTopicDetail(ctx context.Context, arg db.GetTopicDetailParams) (db.GetTopicDetailRow, error)
}

type GetTopicDetailHandlerResponse struct {
	CourseId           string `json:"courseId"`
	SectionId          string `json:"sectionId"`
	TopicId            string `json:"topicId"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	Prerequisites      string `json:"prerequisites"`
	Knowledge          string `json:"knowledge"`
	Flow               string `json:"flow"`
	Quiz               string `json:"quiz"`
	CompletionCriteria string `json:"completionCriteria"`
}

func NewTopicDetailHandler(q GetTopicDetailQuery) *TopicDetailHandler {
	return &TopicDetailHandler{query: q}
}

func (handler *TopicDetailHandler) GetTopicDetailHandler(w http.ResponseWriter, r *http.Request) {
	var topicId pgtype.UUID = pgtype.UUID{}
	var sectionId pgtype.UUID = pgtype.UUID{}
	var courseId pgtype.UUID = pgtype.UUID{}
	var userId pgtype.UUID = pgtype.UUID{}

	if err := topicId.Scan(r.PathValue("topicId")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": response.InputValidationError, "message": "topicId must be a valid UUID"})
		return
	}
	if err := sectionId.Scan(r.PathValue("sectionId")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": response.InputValidationError, "message": "sectionId must be a valid UUID"})
		return
	}
	if err := courseId.Scan(r.PathValue("courseId")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": response.InputValidationError, "message": "courseId must be a valid UUID"})
		return
	}
	if userIdStr := r.URL.Query().Get("userId"); userIdStr != "" {
		if err := userId.Scan(userIdStr); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"code": response.InputValidationError, "message": "userId must be a valid UUID"})
			return
		}
	}

	response, err := handler.query.GetTopicDetail(r.Context(), db.GetTopicDetailParams{
		TopicID:   topicId,
		SectionID: sectionId,
		CourseID:  courseId,
		UserID:    userId,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": "TOPIC_DETAIL_NOT_FOUND", "message": "Topic detail is not found"})
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": "INTERNAL_ERROR", "message": "An internal error occurred"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetTopicDetailHandlerResponse{
		CourseId:           response.CourseId.String(),
		SectionId:          response.SectionId.String(),
		TopicId:            response.TopicId.String(),
		Title:              response.Title,
		Description:        response.Description,
		Prerequisites:      response.Prerequisites,
		Knowledge:          response.Knowledge,
		Flow:               response.Flow,
		Quiz:               response.Quiz,
		CompletionCriteria: response.CompletionCriteria,
	})
}
