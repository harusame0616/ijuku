package commands

import "errors"

type ProgressStatus string

const (
	ProgressStatusInProgress ProgressStatus = "IN_PROGRESS"
	ProgressStatusCompleted  ProgressStatus = "COMPLETED"
)

var ErrInvalidProgressStatus = errors.New("invalid progress status")

func (s ProgressStatus) isValid() bool {
	return s == ProgressStatusInProgress || s == ProgressStatusCompleted
}

// canTransitTo は現在状態から next への遷移が許容されるかを返す。
// COMPLETED から IN_PROGRESS への巻き戻しを禁止する。
func (s ProgressStatus) canTransitTo(next ProgressStatus) bool {
	if s == ProgressStatusCompleted && next == ProgressStatusInProgress {
		return false
	}
	return true
}
