package env

import (
	"os"
)

func Require(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("required environment variable is not set: " + key)
	}

	return v
}
