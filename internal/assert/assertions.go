package assert

import (
	"slices"
	"strings"
	"testing"
	"time"
)

func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func NormalizedStringsEqual(t *testing.T, got, want string) {
	t.Helper()

	normalize := func(s string) string {
		s = strings.TrimSpace(s)
		var builder strings.Builder
		for _, r := range s {
			if r != '\n' && r != '\t' && r != '\r' {
				builder.WriteRune(r)
			}
		}
		return builder.String()
	}

	normalizedGot := normalize(got)
	normalizedWant := normalize(want)

	if normalizedGot != normalizedWant {
		t.Errorf("normalized got %s should be equal to normalized %s", normalizedGot, normalizedWant)
	}
}

func NotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Errorf("got %v should not be equal %v", got, want)
	}
}

func StringContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Errorf("%q should contain %q", got, want)
	}
}

func SliceContains(t *testing.T, got []string, want string) {
	t.Helper()
	if !slices.Contains(got, want) {
		t.Errorf("%q should contain %q", got, want)
	}
}

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
	case string:
		if v == "" {
			t.Errorf("%q is expected to be non-zero", key)
		}
	default:
		t.Fatal("type is not present in NonZero switch, add it")
	}
}

func Zero(t *testing.T, value any, key string) {
	t.Helper()
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		if v != 0 {
			t.Errorf("%q is expected to be zero", key)
		}
	case time.Time:
		if !v.IsZero() {
			t.Errorf("%q is expected to be zero", key)
		}
	case string:
		if v != "" {
			t.Errorf("%q is expected to be zero", key)
		}
	default:
		t.Fatal("type is not present in Zero switch, add it")
	}
}

func NotNil(t *testing.T, ptr any) {
	t.Helper()
	if ptr == nil {
		t.Errorf("expected not to be nil")
	}
}

func True(t *testing.T, b bool) {
	t.Helper()
	if !b {
		t.Errorf("expected to be true")
	}
}

func False(t *testing.T, b bool) {
	t.Helper()
	if b {
		t.Errorf("expected to be true")
	}
}
