package logutils_test

import (
	"testing"
	"time"

	"github.com/Drathveloper/uslogs/internal/logutils"
)

func BenchmarkAppendTimeRFC3339Optimized(b *testing.B) {
	tm := time.Date(2025, 11, 22, 17, 59, 25, 0, time.UTC)
	var buf []byte

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		buf = buf[:0]
		buf = logutils.AppendTimeRFC3339(buf, tm)
	}
}

func BenchmarkTimeFormatRFC3339(b *testing.B) {
	tm := time.Date(2025, 11, 22, 17, 59, 25, 0, time.UTC)

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = tm.Format(time.RFC3339)
	}
}
