package uslog_test

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/Drathveloper/uslogs"
)

type SimulatedDiskWriter struct {
	Latency time.Duration
}

func (s *SimulatedDiskWriter) Write(p []byte) (n int, err error) {
	time.Sleep(s.Latency)
	return len(p), nil
}

func BenchmarkSlogAsyncWriter_Simple(b *testing.B) {
	w := uslog.NewAsyncWriter(io.Discard, 10000)
	defer w.Close()

	logger := slog.New(slog.NewTextHandler(w, nil))

	msg := "benchmark log message"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info(msg)
	}
}

func BenchmarkSlogAsyncWriter_WithAttrs(b *testing.B) {
	w := uslog.NewAsyncWriter(io.Discard, 10000)
	defer w.Close()

	logger := slog.New(slog.NewTextHandler(w, nil))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info("benchmark log with attrs",
			slog.String("user", "alice"),
			slog.Int("id", i),
			slog.Bool("success", true),
		)
	}
}

func BenchmarkSlogAsyncWriter_Parallel(b *testing.B) {
	w := uslog.NewAsyncWriter(io.Discard, 1000)
	h := uslog.NewUnstructuredHandler(
		uslog.WithWriter(w),
		uslog.WithLevel(slog.LevelInfo),
		uslog.WithTimestamp(),
		uslog.WithSeparator('|'),
		uslog.WithResponsivePool())
	defer w.Close()

	logger := slog.New(h)
	args := []any{slog.Int("iteration", 288)}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("parallel benchmark log",
				args...,
			)
		}
	})
}

func BenchmarkSlogRegularWriter_Parallel(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	args := []any{slog.Int("iteration", 288)}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("parallel benchmark log",
				args...,
			)
		}
	})
}

func BenchmarkSlogAsyncWriter_ParallelSimulatedDisk(b *testing.B) {
	disk := &SimulatedDiskWriter{Latency: 50 * time.Microsecond}

	w := uslog.NewAsyncWriter(disk, 10000)
	h := uslog.NewUnstructuredHandler(
		uslog.WithWriter(w),
		uslog.WithLevel(slog.LevelInfo),
		uslog.WithTimestamp(),
		uslog.WithSeparator('|'),
		uslog.WithResponsivePool())
	defer w.Close()

	logger := slog.New(h)
	args := []any{slog.Int("iteration", 288)}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("parallel benchmark log",
				args...,
			)
		}
	})
}

func BenchmarkSlogRegularWriter_ParallelSimulatedDisk(b *testing.B) {
	disk := &SimulatedDiskWriter{Latency: 50 * time.Microsecond}

	logger := slog.New(slog.NewTextHandler(disk, &slog.HandlerOptions{}))
	args := []any{slog.Int("iteration", 288)}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("parallel benchmark log",
				args...,
			)
		}
	})
}
