package commands

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type SqrcCourseRepository struct {
	sqrc db.Querier
}

func NewSqrcCourseRepository(q db.Querier) *SqrcCourseRepository {
	return &SqrcCourseRepository{sqrc: q}
}

func (repository *SqrcCourseRepository) getCourseByCourseId(ctx context.Context, courseId uuid.UUID) (Course, error) {
	var courseIdUuid pgtype.UUID

	if err := courseIdUuid.Scan(courseId.String()); err != nil {
		return Course{}, err
	}

	courseRaw, err := repository.sqrc.GetCourseById(ctx, courseIdUuid)
	if err != nil {
		return Course{}, err
	}

	return toCourse(courseRaw), nil
}

func toCourse(raw db.GetCourseByIdRow) Course {
	var tags []string
	json.Unmarshal(raw.Tags, &tags)

	type topicRaw struct {
		CourseSectionTopicId string `json:"course_section_topic_id"`
		Title                string `json:"title"`
		Description          string `json:"description"`
		Prerequisites        string `json:"prerequisites"`
		Knowledge            string `json:"knowledge"`
		Flow                 string `json:"flow"`
		Quiz                 string `json:"quiz"`
		CompletionCriteria   string `json:"completion_criteria"`
	}
	type sectionRaw struct {
		CourseSectionId string     `json:"course_section_id"`
		Title           string     `json:"title"`
		Description     string     `json:"description"`
		Topics          []topicRaw `json:"topics"`
	}

	var sectionsRaw []sectionRaw
	json.Unmarshal(raw.Sections, &sectionsRaw)

	sections := make([]Section, len(sectionsRaw))
	for i, s := range sectionsRaw {
		topics := make([]Topic, len(s.Topics))
		for j, t := range s.Topics {
			topics[j] = Topic{
				topicId:            t.CourseSectionTopicId,
				title:              t.Title,
				description:        t.Description,
				prerequisites:      t.Prerequisites,
				knowledge:          t.Knowledge,
				flow:               t.Flow,
				quiz:               t.Quiz,
				completionCriteria: t.CompletionCriteria,
				number:             j,
			}
		}
		sections[i] = Section{
			sectionId:   s.CourseSectionId,
			title:       s.Title,
			description: s.Description,
			number:      i,
			topics:      topics,
		}
	}

	var publishedAt *string
	if raw.PublishedAt.Valid {
		t := raw.PublishedAt.Time.Format(time.RFC3339)
		publishedAt = &t
	}

	return Course{
		courseId:      raw.CourseID.String(),
		title:         raw.Title,
		description:   raw.Description,
		slug:          raw.Slug,
		tags:          tags,
		publishStatus: raw.PublishStatus,
		category: Category{
			categoryId: raw.CategoryID.String(),
			name:       raw.CategoryName,
		},
		publishedAt: publishedAt,
		author: Author{
			authorId: raw.AuthorID.String(),
			name:     raw.AuthorName,
		},
		visibility: raw.Visibility,
		sections:   sections,
	}
}
