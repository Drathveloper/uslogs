package uslog_test

import (
	"io"
	"testing"

	"github.com/Drathveloper/uslogs"
)

func BenchmarkAsyncWriter(b *testing.B) {
	w := uslog.NewAsyncWriter(io.Discard, 100)

	data := []byte("log line example\n")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = w.Write(data)
	}

	_ = w.Close()
}
