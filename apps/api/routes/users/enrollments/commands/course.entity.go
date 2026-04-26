package commands

import "errors"

type Topic struct {
	topicId            string
	title              string
	description        string
	prerequisites      string
	knowledge          string
	flow               string
	quiz               string
	completionCriteria string
	number             int
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

var ErrEnrollmentNumberIsNotFound = errors.New("enrollment number is not found")
var ErrEnrollmentNotAllowed = errors.New("enrollment not allowed")
var ErrTopicNumberRequireSectionNumber = errors.New("topic number require section number")

func (course *Course) checkEnrollable(userId string) error {
	if course.publishStatus != "published" && course.author.authorId != userId {
		return ErrEnrollmentNotAllowed
	}
	return nil
}

func (course *Course) findTopicToStart(sectionNumber *int, topicNumber *int, lastPosition ProgressPosition) (Section, Topic, error) {
	if sectionNumber == nil && topicNumber != nil {
		return Section{}, Topic{}, ErrTopicNumberRequireSectionNumber
	}

	var secNum, topNum int

	if sectionNumber != nil {
		secNum = *sectionNumber
		if topicNumber != nil {
			topNum = *topicNumber
		}
	} else {
		switch lastPosition.status {
		case "COMPLETED":
			nextPosition, err := course.nextTopic(TopicPosition{sectionNumber: lastPosition.sectionNumber, topicNumber: lastPosition.topicNumber})
			if err != nil {
				return Section{}, Topic{}, ErrEnrollmentNumberIsNotFound
			}
			secNum, topNum = nextPosition.sectionNumber, nextPosition.topicNumber
		case "IN_PROGRESS":
			secNum = lastPosition.sectionNumber
			topNum = lastPosition.topicNumber
		default:
			secNum = 0
			topNum = 0
		}
	}

	section, topic, err := course.getTopic(secNum, topNum)
	if err != nil {
		return Section{}, Topic{}, ErrEnrollmentNumberIsNotFound
	}
	return section, topic, nil
}

type TopicPosition struct {
	sectionNumber int
	topicNumber   int
}

var ErrTopicPositionOutOfRange = errors.New("topic position is out of range")

func (course *Course) getTopic(sectionNumber int, topicNumber int) (Section, Topic, error) {
	if sectionNumber < 0 || sectionNumber >= len(course.sections) {
		return Section{}, Topic{}, ErrTopicPositionOutOfRange
	}
	section := course.sections[sectionNumber]
	if topicNumber < 0 || topicNumber >= len(section.topics) {
		return Section{}, Topic{}, ErrTopicPositionOutOfRange
	}
	return section, section.topics[topicNumber], nil
}

func (course *Course) nextTopic(currentTopicPosition TopicPosition) (topicPosition TopicPosition, err error) {
	if currentTopicPosition.sectionNumber < 0 || currentTopicPosition.sectionNumber >= len(course.sections) {
		return TopicPosition{}, ErrTopicPositionOutOfRange
	}

	section := course.sections[currentTopicPosition.sectionNumber]
	if currentTopicPosition.topicNumber < 0 || currentTopicPosition.topicNumber >= len(section.topics) {
		return TopicPosition{}, ErrTopicPositionOutOfRange
	}

	if currentTopicPosition.topicNumber == len(section.topics)-1 {
		if currentTopicPosition.sectionNumber == len(course.sections)-1 {
			return TopicPosition{}, ErrTopicPositionOutOfRange
		}
		return TopicPosition{sectionNumber: currentTopicPosition.sectionNumber + 1, topicNumber: 0}, nil
	}

	return TopicPosition{sectionNumber: currentTopicPosition.sectionNumber, topicNumber: currentTopicPosition.topicNumber + 1}, nil
}
