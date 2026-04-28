package queries

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/response"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	progressStatusNotStarted = "NOT_STARTED"
	progressStatusInProgress = "IN_PROGRESS"
	progressStatusCompleted  = "COMPLETED"
	coursePublishStatusPub   = "published"

	errorCodeCourseNotFound     = "COURSE_NOT_FOUND"
	errorCodeEnrollmentForbidden = "ENROLLMENT_FORBIDDEN"
)

type GetEnrollmentQuery interface {
	GetCourseAuthorityById(ctx context.Context, courseid pgtype.UUID) (db.GetCourseAuthorityByIdRow, error)
	GetCourseStructureWithProgress(ctx context.Context, arg db.GetCourseStructureWithProgressParams) (db.GetCourseStructureWithProgressRow, error)
}

type GetEnrollmentHandler struct {
	query GetEnrollmentQuery
}

func NewGetEnrollmentHandler(q GetEnrollmentQuery) *GetEnrollmentHandler {
	return &GetEnrollmentHandler{query: q}
}

// SQL の jsonb 配列をそのまま受ける中間構造体。
// section_agg / jsonb_build_object のキーと一致させる。
type rawTopic struct {
	TopicId string `json:"topicId"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	Index   int    `json:"index"`
}

type rawSection struct {
	SectionId string     `json:"sectionId"`
	Title     string     `json:"title"`
	Index     int        `json:"index"`
	Topics    []rawTopic `json:"topics"`
}

type getEnrollmentTopicResponse struct {
	TopicId string `json:"topicId"`
	Title   string `json:"title"`
	Status  string `json:"status"`
}

type getEnrollmentSectionResponse struct {
	SectionId string                       `json:"sectionId"`
	Title     string                       `json:"title"`
	Topics    []getEnrollmentTopicResponse `json:"topics"`
}

type getEnrollmentNextTopicResponse struct {
	SectionId string `json:"sectionId"`
	TopicId   string `json:"topicId"`
}

type GetEnrollmentResponse struct {
	CourseId  string                          `json:"courseId"`
	Title     string                          `json:"title"`
	Sections  []getEnrollmentSectionResponse  `json:"sections"`
	NextTopic *getEnrollmentNextTopicResponse `json:"nextTopic"`
}

func (h *GetEnrollmentHandler) GetEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("userID")); err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "userID must be a valid UUID")
		return
	}

	var courseID pgtype.UUID
	if err := courseID.Scan(r.PathValue("courseId")); err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, response.InputValidationError, "courseId must be a valid UUID")
		return
	}

	authority, err := h.query.GetCourseAuthorityById(r.Context(), courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteErrorResponse(w, http.StatusNotFound, errorCodeCourseNotFound, "course not found")
			return
		}
		log.Printf("GetCourseAuthorityById error: %v", err)
		response.WriteInternalServerErrorResponse(w)
		return
	}

	if authority.PublishStatus != coursePublishStatusPub && authority.AuthorID.Bytes != userID.Bytes {
		response.WriteErrorResponse(w, http.StatusForbidden, errorCodeEnrollmentForbidden, "this course is not enrollable")
		return
	}

	row, err := h.query.GetCourseStructureWithProgress(r.Context(), db.GetCourseStructureWithProgressParams{
		Courseid: courseID,
		Userid:   userID,
	})
	if err != nil {
		log.Printf("GetCourseStructureWithProgress error: %v", err)
		response.WriteInternalServerErrorResponse(w)
		return
	}

	rawSections, err := unmarshalSections(row.Sections)
	if err != nil {
		log.Printf("unmarshal sections error: %v", err)
		response.WriteInternalServerErrorResponse(w)
		return
	}

	_ = json.NewEncoder(w).Encode(GetEnrollmentResponse{
		CourseId:  courseID.String(),
		Title:     row.Title,
		Sections:  buildSections(rawSections),
		NextTopic: decideNextTopic(rawSections),
	})
}

func unmarshalSections(raw []byte) ([]rawSection, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var sections []rawSection
	if err := json.Unmarshal(raw, &sections); err != nil {
		return nil, err
	}
	return sections, nil
}

func buildSections(rawSecs []rawSection) []getEnrollmentSectionResponse {
	sections := make([]getEnrollmentSectionResponse, 0, len(rawSecs))
	for _, rs := range rawSecs {
		topics := make([]getEnrollmentTopicResponse, 0, len(rs.Topics))
		for _, rt := range rs.Topics {
			topics = append(topics, getEnrollmentTopicResponse{
				TopicId: rt.TopicId,
				Title:   rt.Title,
				Status:  rt.Status,
			})
		}
		sections = append(sections, getEnrollmentSectionResponse{
			SectionId: rs.SectionId,
			Title:     rs.Title,
			Topics:    topics,
		})
	}
	return sections
}

// decideNextTopic は section/topic の index ASC 順に並んだ rawSection から、
// 次に取るべき topic を course.entity.go の findTopicToStart/nextTopic 踏襲で決定する。
//   - 全 NOT_STARTED                                  -> 先頭トピック
//   - (sec_idx, top_idx) 最大の非 NOT_STARTED が IN_PROGRESS -> その topic
//   - 同上が COMPLETED かつ次の topic がある           -> 次の topic
//   - 同上が COMPLETED かつ最終位置 (= 全完了)         -> nil
func decideNextTopic(rawSecs []rawSection) *getEnrollmentNextTopicResponse {
	if len(rawSecs) == 0 {
		return nil
	}

	lastSecIdx, lastTopIdx := -1, -1
	var lastStatus string
	for si, sec := range rawSecs {
		for ti, top := range sec.Topics {
			if top.Status == progressStatusNotStarted {
				continue
			}
			lastSecIdx, lastTopIdx = si, ti
			lastStatus = top.Status
		}
	}

	if lastSecIdx == -1 {
		return firstTopic(rawSecs)
	}

	switch lastStatus {
	case progressStatusInProgress:
		return topicAt(rawSecs, lastSecIdx, lastTopIdx)
	case progressStatusCompleted:
		if nextSec, nextTop, ok := nextPosition(rawSecs, lastSecIdx, lastTopIdx); ok {
			return topicAt(rawSecs, nextSec, nextTop)
		}
		return nil
	default:
		return firstTopic(rawSecs)
	}
}

// nextPosition は course.entity.go の nextTopic を踏襲。
// (secIdx, topIdx) が最終位置なら ok=false を返す。
func nextPosition(rawSecs []rawSection, secIdx, topIdx int) (int, int, bool) {
	if secIdx < 0 || secIdx >= len(rawSecs) {
		return 0, 0, false
	}
	sec := rawSecs[secIdx]
	if topIdx < 0 || topIdx >= len(sec.Topics) {
		return 0, 0, false
	}
	if topIdx == len(sec.Topics)-1 {
		if secIdx == len(rawSecs)-1 {
			return 0, 0, false
		}
		return secIdx + 1, 0, true
	}
	return secIdx, topIdx + 1, true
}

func firstTopic(rawSecs []rawSection) *getEnrollmentNextTopicResponse {
	if len(rawSecs) == 0 || len(rawSecs[0].Topics) == 0 {
		return nil
	}
	return &getEnrollmentNextTopicResponse{
		SectionId: rawSecs[0].SectionId,
		TopicId:   rawSecs[0].Topics[0].TopicId,
	}
}

func topicAt(rawSecs []rawSection, secIdx, topIdx int) *getEnrollmentNextTopicResponse {
	if secIdx < 0 || secIdx >= len(rawSecs) {
		return nil
	}
	sec := rawSecs[secIdx]
	if topIdx < 0 || topIdx >= len(sec.Topics) {
		return nil
	}
	return &getEnrollmentNextTopicResponse{
		SectionId: sec.SectionId,
		TopicId:   sec.Topics[topIdx].TopicId,
	}
}
