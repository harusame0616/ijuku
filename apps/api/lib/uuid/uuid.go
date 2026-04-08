package uuid

import (
	"regexp"

	"github.com/google/uuid"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func IsValidUuid(u string) bool {
	return uuidRegex.MatchString(u)
}

func MustNewUuidString() string {
	id, err := uuid.NewV7()
	if err != nil {
		// crypto/rand の読み取り失敗したケースなど通常環境では基本的に発生しないため panic で終了する
		panic(err)
	}

	return id.String()
}
