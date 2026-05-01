package commands

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAlreadyEnrolled      = errors.New("already enrolled")
	ErrNotEnrolled          = errors.New("not enrolled")
	ErrTopicNotInEnrollment = errors.New("topic not in enrollment")
	ErrInvalidStatusTransit = errors.New("invalid progress status transition")
)

// Enrollment は (userId, courseId) を自然キーとする集約ルート。
type Enrollment struct {
	userId          uuid.UUID
	courseId        uuid.UUID
	enrolledAt      time.Time
	topicProgresses []TopicProgress
}

type TopicProgress struct {
	topicId uuid.UUID
	status  ProgressStatus
}

func NewEnrollment(userId, courseId uuid.UUID, enrolledAt time.Time) Enrollment {
	return Enrollment{
		userId:          userId,
		courseId:        courseId,
		enrolledAt:      enrolledAt,
		topicProgresses: []TopicProgress{},
	}
}

// ReconstructEnrollment は永続化層から取得した値で Enrollment を再構築する。
// ドメイン外（リポジトリ）からの利用のみを想定。
func ReconstructEnrollment(
	userId, courseId uuid.UUID,
	enrolledAt time.Time,
	topicProgresses []TopicProgress,
) Enrollment {
	if topicProgresses == nil {
		topicProgresses = []TopicProgress{}
	}
	return Enrollment{
		userId:          userId,
		courseId:        courseId,
		enrolledAt:      enrolledAt,
		topicProgresses: topicProgresses,
	}
}

func ReconstructTopicProgress(topicId uuid.UUID, status ProgressStatus) TopicProgress {
	return TopicProgress{topicId: topicId, status: status}
}

func (e Enrollment) UserId() uuid.UUID     { return e.userId }
func (e Enrollment) CourseId() uuid.UUID   { return e.courseId }
func (e Enrollment) EnrolledAt() time.Time { return e.enrolledAt }
func (e Enrollment) TopicProgresses() []TopicProgress {
	out := make([]TopicProgress, len(e.topicProgresses))
	copy(out, e.topicProgresses)
	return out
}

func (p TopicProgress) TopicId() uuid.UUID     { return p.topicId }
func (p TopicProgress) Status() ProgressStatus { return p.status }

// UpdateTopicProgress は集約内のトピック進捗を更新する。
// 既存進捗が無ければ新規追加、あれば状態遷移ルールを検証して更新する。
// 戻り値は更新後の TopicProgress。
func (e *Enrollment) UpdateTopicProgress(topicId uuid.UUID, status ProgressStatus) (TopicProgress, error) {
	if !status.isValid() {
		return TopicProgress{}, ErrInvalidProgressStatus
	}

	for i, tp := range e.topicProgresses {
		if tp.topicId != topicId {
			continue
		}
		if !tp.status.canTransitTo(status) {
			return TopicProgress{}, ErrInvalidStatusTransit
		}
		e.topicProgresses[i].status = status
		return e.topicProgresses[i], nil
	}

	tp := TopicProgress{topicId: topicId, status: status}
	e.topicProgresses = append(e.topicProgresses, tp)
	return tp, nil
}
