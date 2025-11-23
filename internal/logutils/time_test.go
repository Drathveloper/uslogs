package logutils_test

import (
	"testing"
	"time"

	"github.com/Drathveloper/uslogs/internal/logutils"
)

func TestAppendTimeRFC3339(t *testing.T) {
	tests := []struct {
		time   time.Time
		expect string
	}{
		{time.Date(2025, 11, 22, 7, 5, 3, 0, time.UTC), "2025-11-22T07:05:03Z"},
		{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), "0001-01-01T00:00:00Z"},
		{time.Date(1999, 12, 31, 23, 59, 59, 0, time.UTC), "1999-12-31T23:59:59Z"},
	}

	for _, tt := range tests {
		var buf []byte
		buf = logutils.AppendTimeRFC3339(buf, tt.time)
		got := string(buf)
		if got != tt.expect {
			t.Errorf("AppendTimeRFC3339(%v) = %q, want %q", tt.time, got, tt.expect)
		}
	}
}
