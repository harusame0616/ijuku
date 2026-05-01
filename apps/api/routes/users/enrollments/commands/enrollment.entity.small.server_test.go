package commands

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestEnrollment_UpdateTopicProgress_NewTopic(t *testing.T) {
	t.Run("未進捗トピックをIN_PROGRESSにすると追加される", func(t *testing.T) {
		userId := uuid.New()
		courseId := uuid.New()
		topicId := uuid.New()
		e := NewEnrollment(userId, courseId, time.Now())

		progress, err := e.UpdateTopicProgress(topicId, ProgressStatusInProgress)
		if err != nil {
			t.Fatalf("エラーが発生しないこと: %v", err)
		}
		if progress.TopicId() != topicId || progress.Status() != ProgressStatusInProgress {
			t.Errorf("更新後のprogressが期待値と一致すること: got %+v", progress)
		}
		if len(e.TopicProgresses()) != 1 {
			t.Errorf("topicProgressesが1件であること: got %d", len(e.TopicProgresses()))
		}
	})

	t.Run("未進捗トピックをCOMPLETEDで開始することもできる", func(t *testing.T) {
		e := NewEnrollment(uuid.New(), uuid.New(), time.Now())
		topicId := uuid.New()

		progress, err := e.UpdateTopicProgress(topicId, ProgressStatusCompleted)
		if err != nil {
			t.Fatalf("エラーが発生しないこと: %v", err)
		}
		if progress.Status() != ProgressStatusCompleted {
			t.Errorf("statusがCOMPLETEDであること: got %q", progress.Status())
		}
	})

	t.Run("不正なstatusはErrInvalidProgressStatusを返す", func(t *testing.T) {
		e := NewEnrollment(uuid.New(), uuid.New(), time.Now())
		_, err := e.UpdateTopicProgress(uuid.New(), ProgressStatus("FOO"))
		if !errors.Is(err, ErrInvalidProgressStatus) {
			t.Errorf("ErrInvalidProgressStatusが返ること: got %v", err)
		}
	})
}

func TestEnrollment_UpdateTopicProgress_Existing(t *testing.T) {
	topicId := uuid.New()
	newEnrollment := func(initial ProgressStatus) Enrollment {
		return ReconstructEnrollment(
			uuid.New(), uuid.New(), time.Now(),
			[]TopicProgress{ReconstructTopicProgress(topicId, initial)},
		)
	}

	t.Run("IN_PROGRESS から COMPLETED への遷移は許可される", func(t *testing.T) {
		e := newEnrollment(ProgressStatusInProgress)
		progress, err := e.UpdateTopicProgress(topicId, ProgressStatusCompleted)
		if err != nil {
			t.Fatalf("エラーが発生しないこと: %v", err)
		}
		if progress.Status() != ProgressStatusCompleted {
			t.Errorf("statusがCOMPLETEDであること: got %q", progress.Status())
		}
	})

	t.Run("COMPLETED から IN_PROGRESS への巻き戻しはErrInvalidStatusTransitを返す", func(t *testing.T) {
		e := newEnrollment(ProgressStatusCompleted)
		_, err := e.UpdateTopicProgress(topicId, ProgressStatusInProgress)
		if !errors.Is(err, ErrInvalidStatusTransit) {
			t.Errorf("ErrInvalidStatusTransitが返ること: got %v", err)
		}
		if e.TopicProgresses()[0].Status() != ProgressStatusCompleted {
			t.Errorf("statusはCOMPLETEDのままであること: got %q", e.TopicProgresses()[0].Status())
		}
	})

	t.Run("COMPLETED から COMPLETED への再設定は許可される", func(t *testing.T) {
		e := newEnrollment(ProgressStatusCompleted)
		progress, err := e.UpdateTopicProgress(topicId, ProgressStatusCompleted)
		if err != nil {
			t.Fatalf("エラーが発生しないこと: %v", err)
		}
		if progress.Status() != ProgressStatusCompleted {
			t.Errorf("statusがCOMPLETEDであること: got %q", progress.Status())
		}
	})

	t.Run("IN_PROGRESS から IN_PROGRESS への再設定も許可される", func(t *testing.T) {
		e := newEnrollment(ProgressStatusInProgress)
		_, err := e.UpdateTopicProgress(topicId, ProgressStatusInProgress)
		if err != nil {
			t.Errorf("エラーが発生しないこと: %v", err)
		}
	})
}

func TestEnrollment_Accessors(t *testing.T) {
	t.Run("コンストラクタで指定した値が取得できる", func(t *testing.T) {
		userId := uuid.New()
		courseId := uuid.New()
		now := time.Now()
		e := NewEnrollment(userId, courseId, now)

		if e.UserId() != userId {
			t.Errorf("UserIdが一致すること: got %v", e.UserId())
		}
		if e.CourseId() != courseId {
			t.Errorf("CourseIdが一致すること: got %v", e.CourseId())
		}
		if !e.EnrolledAt().Equal(now) {
			t.Errorf("EnrolledAtが一致すること: got %v", e.EnrolledAt())
		}
	})
}
