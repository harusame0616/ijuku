package queries

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	libauth "github.com/harusame0616/ijuku/apps/api/lib/auth"
	"github.com/harusame0616/ijuku/apps/api/lib/response"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	coursePublishedStatus = "published"
	courseVisibilityPub   = "public"

	errorCodeCourseDetailNotFound = "COURSE_NOT_FOUND"
)

type GetCourseDetailQuery interface {
	GetCourseBySlug(ctx context.Context, arg db.GetCourseBySlugParams) (db.GetCourseBySlugRow, error)
	GetEnrollmentByUserIdAndCourseId(ctx context.Context, arg db.GetEnrollmentByUserIdAndCourseIdParams) (db.GetEnrollmentByUserIdAndCourseIdRow, error)
}

type GetCourseDetailHandler struct {
	query GetCourseDetailQuery
}

func NewGetCourseDetailHandler(q GetCourseDetailQuery) *GetCourseDetailHandler {
	return &GetCourseDetailHandler{query: q}
}

type courseDetailTopicResponse struct {
	TopicId     string `json:"topicId"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type courseDetailSectionResponse struct {
	SectionId   string                      `json:"sectionId"`
	Title       string                      `json:"title"`
	Description string                      `json:"description"`
	Topics      []courseDetailTopicResponse `json:"topics"`
}

type courseDetailAuthorResponse struct {
	AuthorId string `json:"authorId"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
}

type GetCourseDetailResponse struct {
	CourseId    string                        `json:"courseId"`
	Title       string                        `json:"title"`
	Description string                        `json:"description"`
	Slug        string                        `json:"slug"`
	Tags        []string                      `json:"tags"`
	Author      courseDetailAuthorResponse    `json:"author"`
	Sections    []courseDetailSectionResponse `json:"sections"`
	IsEnrolled  bool                          `json:"isEnrolled"`
}

func (h *GetCourseDetailHandler) GetCourseDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorSlug := r.PathValue("authorSlug")
	courseSlug := r.PathValue("courseSlug")
	if authorSlug == "" || courseSlug == "" {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "authorSlug and courseSlug are required")
		return
	}

	row, err := h.query.GetCourseBySlug(r.Context(), db.GetCourseBySlugParams{
		Authorslug: authorSlug,
		Courseslug: courseSlug,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteErrorResponse(w, http.StatusNotFound, errorCodeCourseDetailNotFound, "course not found")
			return
		}
		log.Printf("GetCourseBySlug error: %v", err)
		response.WriteInternalServerErrorResponse(w)
		return
	}

	if row.PublishStatus != coursePublishedStatus || row.Visibility != courseVisibilityPub {
		response.WriteErrorResponse(w, http.StatusNotFound, errorCodeCourseDetailNotFound, "course not found")
		return
	}

	sections, err := buildCourseDetailSections(row.Sections)
	if err != nil {
		log.Printf("unmarshal sections error: %v", err)
		response.WriteInternalServerErrorResponse(w)
		return
	}

	tags := buildCourseDetailTags(row.Tags)

	isEnrolled := false
	if userIdStr, ok := libauth.UserIDFromContext(r.Context()); ok {
		var userId pgtype.UUID
		if err := userId.Scan(userIdStr); err == nil {
			_, err := h.query.GetEnrollmentByUserIdAndCourseId(r.Context(), db.GetEnrollmentByUserIdAndCourseIdParams{
				Userid:   userId,
				Courseid: row.CourseID,
			})
			if err == nil {
				isEnrolled = true
			} else if !errors.Is(err, pgx.ErrNoRows) {
				log.Printf("GetEnrollmentByUserIdAndCourseId error: %v", err)
				response.WriteInternalServerErrorResponse(w)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetCourseDetailResponse{
		CourseId:    row.CourseID.String(),
		Title:       row.Title,
		Description: row.Description,
		Slug:        row.Slug,
		Tags:        tags,
		Author: courseDetailAuthorResponse{
			AuthorId: row.AuthorID.String(),
			Name:     row.AuthorName,
			Slug:     row.AuthorSlug,
		},
		Sections:   sections,
		IsEnrolled: isEnrolled,
	})
}

type courseDetailSectionRaw struct {
	CourseSectionId string                  `json:"course_section_id"`
	Title           string                  `json:"title"`
	Description     string                  `json:"description"`
	Topics          []courseDetailTopicRaw  `json:"topics"`
}

type courseDetailTopicRaw struct {
	CourseSectionTopicId string `json:"course_section_topic_id"`
	Title                string `json:"title"`
	Description          string `json:"description"`
}

// 既存レコードの tags カラムは JSON 配列以外（単純な文字列など）が入っているケースがあるため、
// unmarshal 失敗時は空配列にフォールバックしてレスポンス 200 を維持する。
func buildCourseDetailTags(raw []byte) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var tags []string
	if err := json.Unmarshal(raw, &tags); err != nil {
		log.Printf("tags unmarshal failed (treating as empty): %v", err)
		return []string{}
	}
	return tags
}

func buildCourseDetailSections(raw []byte) ([]courseDetailSectionResponse, error) {
	if len(raw) == 0 {
		return []courseDetailSectionResponse{}, nil
	}
	var rawSections []courseDetailSectionRaw
	if err := json.Unmarshal(raw, &rawSections); err != nil {
		return nil, err
	}
	sections := make([]courseDetailSectionResponse, 0, len(rawSections))
	for _, rs := range rawSections {
		topics := make([]courseDetailTopicResponse, 0, len(rs.Topics))
		for _, rt := range rs.Topics {
			topics = append(topics, courseDetailTopicResponse{
				TopicId:     rt.CourseSectionTopicId,
				Title:       rt.Title,
				Description: rt.Description,
			})
		}
		sections = append(sections, courseDetailSectionResponse{
			SectionId:   rs.CourseSectionId,
			Title:       rs.Title,
			Description: rs.Description,
			Topics:      topics,
		})
	}
	return sections, nil
}
