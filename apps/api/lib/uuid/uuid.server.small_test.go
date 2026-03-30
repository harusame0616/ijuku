package uuid

import "testing"

func TestUuidSmallTest(t *testing.T) {
	t.Run("有効な UUID 文字列の場合、 true を返す", func(t *testing.T) {
		validUuid := "123e4567-e89b-12d3-a456-426614174000"
		result := IsValidUuid(validUuid)

		if result != true {
			t.Errorf("戻り値が true であること: got %v", result)
		}
	})

	t.Run("空文字の場合、 false を返す", func(t *testing.T) {
		validUuid := ""
		result := IsValidUuid(validUuid)

		if result != false {
			t.Errorf("戻り値が false であること: got %v", result)
		}
	})

	t.Run("不正な UUID 文字列の場合、 false を返す", func(t *testing.T) {
		invalidUuid := "invalid-uuid-string"
		result := IsValidUuid(invalidUuid)

		if result != false {
			t.Errorf("戻り値が false であること: got %v", result)
		}
	})
}
