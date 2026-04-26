package commands

import "time"

type ProgressPosition struct {
	sectionNumber int
	topicNumber   int
	status        string
}

type ProgressRecord struct {
	courseSectionTopicId string
	sectionNumber        int
	topicNumber          int
	status               string
	startedAt            string
	completedAt          *string
}

type Progress struct {
	courseId string
	userId   string
	records  []ProgressRecord
}

func (p *Progress) lastTopicPosition() ProgressPosition {
	pos := ProgressPosition{sectionNumber: -1, topicNumber: -1, status: "NOT_STARTED"}

	for _, record := range p.records {
		if record.sectionNumber > pos.sectionNumber || (record.sectionNumber == pos.sectionNumber && record.topicNumber > pos.topicNumber) {
			pos = ProgressPosition{sectionNumber: record.sectionNumber, topicNumber: record.topicNumber, status: record.status}
		}
	}

	return pos
}

func (p *Progress) start(section Section, topic Topic) string {
	for _, record := range p.records {
		if record.courseSectionTopicId == topic.topicId {
			return topic.topicId
		}
	}

	p.records = append(p.records, ProgressRecord{
		courseSectionTopicId: topic.topicId,
		sectionNumber:        section.number,
		topicNumber:          topic.number,
		status:               "IN_PROGRESS",
		startedAt:            time.Now().Format(time.RFC3339),
		completedAt:          nil,
	})
	return topic.topicId
}
