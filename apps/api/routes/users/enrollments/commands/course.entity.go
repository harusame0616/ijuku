package commands

import (
	"errors"

	"github.com/google/uuid"
)

type Topic struct {
	topicId     string
	title       string
	description string
	content     string
	number      int
}

type Section struct {
	sectionId   string
	title       string
	description string
	number      int
	topics      []Topic
}

type Author struct {
	authorId string
	name     string
}

type Category struct {
	categoryId string
	name       string
}

type Course struct {
	courseId      string
	title         string
	description   string
	slug          string
	tags          []string
	publishStatus string
	category      Category
	publishedAt   *string
	author        Author
	visibility    string
	sections      []Section
}

var ErrEnrollmentNotAllowed = errors.New("enrollment not allowed")
var ErrTopicNotFoundInCourse = errors.New("topic not found in course")

func (course *Course) checkEnrollable(userId uuid.UUID) error {
	if course.publishStatus != "published" && course.author.authorId != userId.String() {
		return ErrEnrollmentNotAllowed
	}
	return nil
}

func (course *Course) findTopicById(topicId uuid.UUID) (Topic, error) {
	for _, section := range course.sections {
		for _, topic := range section.topics {
			if topic.topicId == topicId.String() {
				return topic, nil
			}
		}
	}
	return Topic{}, ErrTopicNotFoundInCourse
}
