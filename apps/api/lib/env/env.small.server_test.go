package env_test

import (
	"os"
	"testing"

	"github.com/harusame0616/ijuku/apps/api/lib/env"
	"github.com/stretchr/testify/assert"
)

func TestRequire(t *testing.T) {
	t.Run("環境変数が設定されている場合、値を返す", func(t *testing.T) {
		t.Setenv("TEST_REQUIRE_KEY", "test_value")

		got := env.Require("TEST_REQUIRE_KEY")

		assert.Equal(t, "test_value", got)
	})

	t.Run("環境変数が設定されていない場合、panic する", func(t *testing.T) {
		os.Unsetenv("TEST_REQUIRE_KEY_NOT_SET")

		assert.Panics(t, func() {
			env.Require("TEST_REQUIRE_KEY_NOT_SET")
		})
	})
}
