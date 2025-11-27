package uslogs_test

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/Drathveloper/uslogs"
)

var output = io.Discard

func BenchmarkLogWriter_Handle(b *testing.B) {
	writer := uslogs.NewUnstructuredHandler(
		uslogs.WithLevel(slog.LevelInfo),
		uslogs.WithWriter(output),
		uslogs.WithSeparator('|'))
	msg := strings.Repeat("X", 63*1024)
	record := slog.NewRecord(time.Now(), slog.LevelInfo, msg, 0)
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = writer.Handle(context.Background(), record)
		}
	})
}

func BenchmarkLogWriter_HandleWithTime(b *testing.B) {
	writer := uslogs.NewUnstructuredHandler(
		uslogs.WithLevel(slog.LevelInfo),
		uslogs.WithWriter(output),
		uslogs.WithSeparator('|'),
		uslogs.WithTimestamp())
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "It was a simple tip of the hat. Grace didn't think that anyone else besides her had even noticed it", 0)
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = writer.Handle(context.Background(), record)
		}
	})
}

func BenchmarkLogWriter_HandleWithTimeAndAttrs(b *testing.B) {
	writer := uslogs.NewUnstructuredHandler(
		uslogs.WithLevel(slog.LevelInfo),
		uslogs.WithWriter(output),
		uslogs.WithSeparator('|'),
		uslogs.WithTimestamp())
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "It was a simple tip of the hat. Grace didn't think that anyone else besides her had even noticed it", 0)
	record.Add(slog.String("foo", "bar"), slog.Int("baz", 25))

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = writer.Handle(context.Background(), record)
		}
	})
}

func BenchmarkSlogWriter_Handle(b *testing.B) {
	writer := uslogs.NewUnstructuredHandler(
		uslogs.WithLevel(slog.LevelInfo),
		uslogs.WithWriter(output),
		uslogs.WithSeparator('|'))
	logger := slog.New(writer)
	msg := strings.Repeat("X", 63*1024)
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(msg)
		}
	})
}

func BenchmarkSlogWriter_HandleWithTime(b *testing.B) {
	writer := uslogs.NewUnstructuredHandler(
		uslogs.WithLevel(slog.LevelInfo),
		uslogs.WithWriter(output),
		uslogs.WithSeparator('|'),
		uslogs.WithTimestamp())
	logger := slog.New(writer)
	msg := strings.Repeat("X", 63*1024)
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(msg)
		}
	})
}

func BenchmarkSlogWriter_HandleWithTimeAndAttrs(b *testing.B) {
	writer := uslogs.NewUnstructuredHandler(
		uslogs.WithLevel(slog.LevelInfo),
		uslogs.WithWriter(output),
		uslogs.WithSeparator('|'),
		uslogs.WithTimestamp())
	logger := slog.New(writer)
	msg := strings.Repeat("X", 50*1024)
	attr1 := slog.String("foo", "bar")
	attr2 := slog.Int("baz", 25)
	args := []any{attr1, attr2}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(msg, args...)
		}
	})
}

func BenchmarkSlogWriter_HandleWithAll_64(b *testing.B) {
	writer := uslogs.NewUnstructuredHandler(
		uslogs.WithLevel(slog.LevelInfo),
		uslogs.WithWriter(output),
		uslogs.WithSeparator('|'),
		uslogs.WithMaskedAttributes("foo"),
		uslogs.WithTimestamp())
	logger := slog.New(writer)
	msg := strings.Repeat("X", 64*1024-36)
	attr1 := slog.String("foo", "bar")
	attr2 := slog.Int("baz", 25)
	args := []any{attr1, attr2}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(msg, args...)
		}
	})
}

func BenchmarkSlogWriter_HandleWithAll_32(b *testing.B) {
	writer := uslogs.NewUnstructuredHandler(
		uslogs.WithLevel(slog.LevelInfo),
		uslogs.WithWriter(output),
		uslogs.WithSeparator('|'),
		uslogs.WithMaskedAttributes("foo"),
		uslogs.WithTimestamp())
	logger := slog.New(writer)
	msg := strings.Repeat("X", 32*1024-36)
	attr1 := slog.String("foo", "bar")
	attr2 := slog.Int("baz", 25)
	args := []any{attr1, attr2}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(msg, args...)
		}
	})
}

func BenchmarkSlogWriter_HandleWithAll_64Responsive(b *testing.B) {
	writer := uslogs.NewUnstructuredHandler(
		uslogs.WithLevel(slog.LevelInfo),
		uslogs.WithWriter(output),
		uslogs.WithSeparator('|'),
		uslogs.WithMaskedAttributes("foo"),
		uslogs.WithTimestamp(),
		uslogs.WithResponsivePool())
	logger := slog.New(writer)
	msg := strings.Repeat("X", 64*1024-36)
	attr1 := slog.String("foo", "bar")
	attr2 := slog.Int("baz", 25)
	args := []any{attr1, attr2}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(msg, args...)
		}
	})
}

func BenchmarkSlogWriter_HandleWithAll_32Responsive(b *testing.B) {
	writer := uslogs.NewUnstructuredHandler(
		uslogs.WithLevel(slog.LevelInfo),
		uslogs.WithWriter(output),
		uslogs.WithSeparator('|'),
		uslogs.WithMaskedAttributes("foo"),
		uslogs.WithTimestamp(),
		uslogs.WithResponsivePool())
	logger := slog.New(writer)
	msg := strings.Repeat("X", 32*1024-36)
	attr1 := slog.String("foo", "bar")
	attr2 := slog.Int("baz", 25)
	args := []any{attr1, attr2}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(msg, args...)
		}
	})
}
