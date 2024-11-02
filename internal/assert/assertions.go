package assert

import (
	"testing"
	"time"
)

func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("received unexpected error, %v", err)
	}
}

func NonZero(t *testing.T, value any, key string) {
	t.Helper()
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		if v == 0 {
			t.Errorf("%q is expected to be non-zero", key)
		}
	case time.Time:
		if v.IsZero() {
			t.Errorf("%q is expected to be non-zero", key)
		}
	default:
		t.Fatal("type is not present in NonZero switch, add it")
	}
}
