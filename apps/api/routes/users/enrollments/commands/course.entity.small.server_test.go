package commands

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestCourse_findTopicById(t *testing.T) {
	topic00 := uuid.New()
	topic01 := uuid.New()
	topic10 := uuid.New()

	course := Course{
		sections: []Section{
			{
				topics: []Topic{
					{topicId: topic00.String()},
					{topicId: topic01.String()},
				},
			},
			{
				topics: []Topic{
					{topicId: topic10.String()},
				},
			},
		},
	}

	t.Run("一致するtopicIdがある場合は該当Topicを返す", func(t *testing.T) {
		topic, err := course.findTopicById(topic10)
		if err != nil {
			t.Fatalf("エラーが発生しないこと: %v", err)
		}
		if topic.topicId != topic10.String() {
			t.Errorf("topicIdが一致すること: got %q", topic.topicId)
		}
	})

	t.Run("一致するtopicIdがない場合はErrTopicNotFoundInCourseを返す", func(t *testing.T) {
		_, err := course.findTopicById(uuid.New())
		if !errors.Is(err, ErrTopicNotFoundInCourse) {
			t.Errorf("ErrTopicNotFoundInCourseが返ること: got %v", err)
		}
	})

	t.Run("sectionsが空の場合はErrTopicNotFoundInCourseを返す", func(t *testing.T) {
		empty := Course{}
		_, err := empty.findTopicById(topic00)
		if !errors.Is(err, ErrTopicNotFoundInCourse) {
			t.Errorf("ErrTopicNotFoundInCourseが返ること: got %v", err)
		}
	})
}

func TestCourse_checkEnrollable(t *testing.T) {
	author := uuid.New()
	other := uuid.New()

	t.Run("publishedの場合は誰でも受講可能", func(t *testing.T) {
		course := Course{publishStatus: "published", author: Author{authorId: author.String()}}
		if err := course.checkEnrollable(other); err != nil {
			t.Errorf("エラーが発生しないこと: %v", err)
		}
	})

	t.Run("draftで著者本人は受講可能", func(t *testing.T) {
		course := Course{publishStatus: "draft", author: Author{authorId: author.String()}}
		if err := course.checkEnrollable(author); err != nil {
			t.Errorf("エラーが発生しないこと: %v", err)
		}
	})

	t.Run("draftで非著者はErrEnrollmentNotAllowedを返す", func(t *testing.T) {
		course := Course{publishStatus: "draft", author: Author{authorId: author.String()}}
		if err := course.checkEnrollable(other); !errors.Is(err, ErrEnrollmentNotAllowed) {
			t.Errorf("ErrEnrollmentNotAllowedが返ること: got %v", err)
		}
	})
}
