package uslogs_test

import (
	"io"
	"testing"

	"github.com/Drathveloper/uslogs"
)

func BenchmarkAsyncWriter(b *testing.B) {
	w := uslogs.NewAsyncWriter(io.Discard, 100)

	data := []byte("log line example\n")

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, _ = w.Write(data)
	}

	_ = w.Close()
}
