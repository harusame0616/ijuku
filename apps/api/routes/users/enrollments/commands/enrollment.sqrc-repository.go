package commands

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type SqrcEnrollmentRepository struct {
	sqrc db.Querier
}

func NewSqrcEnrollmentRepository(q db.Querier) *SqrcEnrollmentRepository {
	return &SqrcEnrollmentRepository{sqrc: q}
}

func (r *SqrcEnrollmentRepository) findByUserAndCourse(ctx context.Context, userId, courseId uuid.UUID) (Enrollment, error) {
	userIdUuid, err := toPgUUID(userId)
	if err != nil {
		return Enrollment{}, err
	}
	courseIdUuid, err := toPgUUID(courseId)
	if err != nil {
		return Enrollment{}, err
	}

	row, err := r.sqrc.GetEnrollmentByUserIdAndCourseId(ctx, db.GetEnrollmentByUserIdAndCourseIdParams{
		Userid:   userIdUuid,
		Courseid: courseIdUuid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Enrollment{}, ErrNotEnrolled
		}
		return Enrollment{}, err
	}

	progresses, err := r.sqrc.GetTopicProgressesByUserIdAndCourseId(ctx, db.GetTopicProgressesByUserIdAndCourseIdParams{
		Userid:   userIdUuid,
		Courseid: courseIdUuid,
	})
	if err != nil {
		return Enrollment{}, err
	}

	topicProgresses := make([]TopicProgress, 0, len(progresses))
	for _, p := range progresses {
		topicId, err := uuid.FromBytes(p.CourseSectionTopicID.Bytes[:])
		if err != nil {
			return Enrollment{}, err
		}
		topicProgresses = append(topicProgresses, ReconstructTopicProgress(topicId, ProgressStatus(p.Status)))
	}

	return ReconstructEnrollment(userId, courseId, row.EnrolledAt.Time, topicProgresses), nil
}

func (r *SqrcEnrollmentRepository) create(ctx context.Context, enrollment Enrollment) error {
	userIdUuid, err := toPgUUID(enrollment.UserId())
	if err != nil {
		return err
	}
	courseIdUuid, err := toPgUUID(enrollment.CourseId())
	if err != nil {
		return err
	}

	var enrolledAt pgtype.Timestamptz
	if err := enrolledAt.Scan(enrollment.EnrolledAt()); err != nil {
		return err
	}

	return r.sqrc.InsertEnrollment(ctx, db.InsertEnrollmentParams{
		Userid:     userIdUuid,
		Courseid:   courseIdUuid,
		Enrolledat: enrolledAt,
	})
}

func (r *SqrcEnrollmentRepository) upsertTopicProgress(ctx context.Context, userId uuid.UUID, progress TopicProgress) error {
	userIdUuid, err := toPgUUID(userId)
	if err != nil {
		return err
	}
	topicIdUuid, err := toPgUUID(progress.TopicId())
	if err != nil {
		return err
	}

	return r.sqrc.UpsertTopicProgress(ctx, db.UpsertTopicProgressParams{
		Userid:               userIdUuid,
		Coursesectiontopicid: topicIdUuid,
		Status:               string(progress.Status()),
	})
}

func toPgUUID(id uuid.UUID) (pgtype.UUID, error) {
	var pg pgtype.UUID
	if err := pg.Scan(id.String()); err != nil {
		return pgtype.UUID{}, err
	}
	return pg, nil
}
