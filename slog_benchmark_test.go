package uslogs_test

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

func (s *SimulatedDiskWriter) Write(p []byte) (int, error) {
	time.Sleep(s.Latency)
	return len(p), nil
}

func BenchmarkSlogAsyncWriter_Simple(b *testing.B) {
	w := uslogs.NewAsyncWriter(io.Discard, 10000)
	defer w.Close() //nolint:errcheck

	logger := slog.New(slog.NewTextHandler(w, nil))

	msg := "benchmark log message"

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		logger.Info(msg)
	}
}

func BenchmarkSlogAsyncWriter_WithAttrs(b *testing.B) {
	w := uslogs.NewAsyncWriter(io.Discard, 10000)
	defer w.Close() //nolint:errcheck

	logger := slog.New(slog.NewTextHandler(w, nil))
	attrs := []any{slog.String("user", "alice"),
		slog.Int("id", 288),
		slog.Bool("success", true)}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		logger.Info("benchmark log with attrs", attrs...)
	}
}

func BenchmarkSlogAsyncWriter_Parallel(b *testing.B) {
	w := uslogs.NewAsyncWriter(io.Discard, 1000)
	h := uslogs.NewUnstructuredHandler(
		uslogs.WithWriter(w),
		uslogs.WithLevel(slog.LevelInfo),
		uslogs.WithTimestamp(),
		uslogs.WithSeparator('|'),
		uslogs.WithResponsivePool())
	defer w.Close() //nolint:errcheck

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
	logger := slog.New(slog.DiscardHandler)
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

	w := uslogs.NewAsyncWriter(disk, 10000)
	h := uslogs.NewUnstructuredHandler(
		uslogs.WithWriter(w),
		uslogs.WithLevel(slog.LevelInfo),
		uslogs.WithTimestamp(),
		uslogs.WithSeparator('|'),
		uslogs.WithResponsivePool())
	defer w.Close() //nolint:errcheck

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

	//nolint:exhaustruct
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
