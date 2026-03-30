package commands

import (
	"context"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type SqrcUserTopicProgressRepository struct {
	sqrc db.Querier
}

func NewSqrcUserTopicProgressRepository(q db.Querier) *SqrcUserTopicProgressRepository {
	return &SqrcUserTopicProgressRepository{sqrc: q}
}

func (repository *SqrcUserTopicProgressRepository) getProgress(ctx context.Context, userId string, courseId string) (Progress, error) {
	var userIdUuid, courseIdUuid pgtype.UUID

	if err := userIdUuid.Scan(userId); err != nil {
		return Progress{}, err
	}
	if err := courseIdUuid.Scan(courseId); err != nil {
		return Progress{}, err
	}

	rows, err := repository.sqrc.GetProgressByUserIdAndCourseId(ctx, db.GetProgressByUserIdAndCourseIdParams{
		Userid:   userIdUuid,
		Courseid: courseIdUuid,
	})
	if err != nil {
		return Progress{}, err
	}

	records := make([]ProgressRecord, len(rows))
	for i, row := range rows {
		records[i] = ProgressRecord{
			courseSectionTopicId: row.CourseSectionTopicID.String(),
			sectionNumber:        int(row.SectionIndex),
			topicNumber:          int(row.TopicIndex),
			status:               row.Status,
		}
	}

	return Progress{courseId: courseId, userId: userId, records: records}, nil
}

func (repository *SqrcUserTopicProgressRepository) save(ctx context.Context, progress *Progress) error {
	for _, record := range progress.records {
		var topicIdUuid, userIdUuid pgtype.UUID
		if err := topicIdUuid.Scan(record.courseSectionTopicId); err != nil {
			return err
		}
		if err := userIdUuid.Scan(progress.userId); err != nil {
			return err
		}

		if err := repository.sqrc.UpsertProgress(ctx, db.UpsertProgressParams{
			Coursesectiontopicid: topicIdUuid,
			Userid:               userIdUuid,
			Status:               record.status,
		}); err != nil {
			return err
		}
	}
	return nil
}
