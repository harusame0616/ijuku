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

	raw, err := repository.sqrc.GetCourseById(ctx, courseIdUuid)
	if err != nil {
		return Course{}, err
	}

	return buildCourse(courseRawFields{
		CourseID:      raw.CourseID,
		Title:         raw.Title,
		Description:   raw.Description,
		Slug:          raw.Slug,
		Tags:          raw.Tags,
		PublishStatus: raw.PublishStatus,
		CategoryID:    raw.CategoryID,
		CategoryName:  raw.CategoryName,
		PublishedAt:   raw.PublishedAt,
		AuthorID:      raw.AuthorID,
		AuthorName:    raw.AuthorName,
		Visibility:    raw.Visibility,
		Sections:      raw.Sections,
	}), nil
}

func (repository *SqrcCourseRepository) getCourseBySlug(ctx context.Context, authorSlug, courseSlug string) (Course, error) {
	raw, err := repository.sqrc.GetCourseBySlug(ctx, db.GetCourseBySlugParams{
		Authorslug: authorSlug,
		Courseslug: courseSlug,
	})
	if err != nil {
		return Course{}, err
	}

	return buildCourse(courseRawFields{
		CourseID:      raw.CourseID,
		Title:         raw.Title,
		Description:   raw.Description,
		Slug:          raw.Slug,
		Tags:          raw.Tags,
		PublishStatus: raw.PublishStatus,
		CategoryID:    raw.CategoryID,
		CategoryName:  raw.CategoryName,
		PublishedAt:   raw.PublishedAt,
		AuthorID:      raw.AuthorID,
		AuthorName:    raw.AuthorName,
		Visibility:    raw.Visibility,
		Sections:      raw.Sections,
	}), nil
}

type courseRawFields struct {
	CourseID      pgtype.UUID
	Title         string
	Description   string
	Slug          string
	Tags          []byte
	PublishStatus string
	CategoryID    pgtype.UUID
	CategoryName  string
	PublishedAt   pgtype.Timestamptz
	AuthorID      pgtype.UUID
	AuthorName    string
	Visibility    string
	Sections      []byte
}

func buildCourse(raw courseRawFields) Course {
	var tags []string
	json.Unmarshal(raw.Tags, &tags)

	type topicRaw struct {
		CourseSectionTopicId string `json:"course_section_topic_id"`
		Title                string `json:"title"`
		Description          string `json:"description"`
		Content              string `json:"content"`
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
				topicId:     t.CourseSectionTopicId,
				title:       t.Title,
				description: t.Description,
				content:     t.Content,
				number:      j,
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
		courseId:      uuid.UUID(raw.CourseID.Bytes),
		title:         raw.Title,
		description:   raw.Description,
		slug:          raw.Slug,
		tags:          tags,
		publishStatus: raw.PublishStatus,
		category: Category{
			categoryId: uuid.UUID(raw.CategoryID.Bytes),
			name:       raw.CategoryName,
		},
		publishedAt: publishedAt,
		author: Author{
			authorId: uuid.UUID(raw.AuthorID.Bytes),
			name:     raw.AuthorName,
		},
		visibility: raw.Visibility,
		sections:   sections,
	}
}
