package uuid

import "regexp"

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func IsValidUuid(u string) bool {
	return uuidRegex.MatchString(u)
}
