package uslogs_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/Drathveloper/uslogs"
)

func TestEnabled(t *testing.T) {
	handler := uslogs.NewUnstructuredHandler(
		uslogs.WithWriter(io.Discard),
		uslogs.WithLevel(slog.LevelInfo))

	tests := []struct {
		level    slog.Level
		expected bool
	}{
		{slog.LevelDebug, false},
		{slog.LevelInfo, true},
		{slog.LevelWarn, true},
	}

	for _, tt := range tests {
		got := handler.Enabled(context.Background(), tt.level)
		if got != tt.expected {
			t.Errorf("Enabled(%v) = %v, want %v", tt.level, got, tt.expected)
		}
	}
}

func TestHandleWritesMessage(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := uslogs.NewUnstructuredHandler(uslogs.WithWriter(buf))

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 1)
	if err := handler.Handle(context.Background(), record); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "INFO") || !strings.Contains(out, "test message") {
		t.Errorf("output = %q, want it to contain level and message", out)
	}
}

func TestHandleWithAttrsAndAllValueKinds(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := uslogs.NewUnstructuredHandler(uslogs.WithWriter(buf))

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 1)
	record.AddAttrs(
		slog.String("s", "string"),
		slog.Int64("i", 42),
		slog.Uint64("u", 100),
		slog.Float64("f", 3.14),
		slog.Bool("b", true),
	)

	if err := handler.Handle(context.Background(), record); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	out := buf.String()
	for _, substr := range []string{"s=string", "i=42", "u=100", "f=3.14", "b=true"} {
		if !strings.Contains(out, substr) {
			t.Errorf("output = %q, want it to contain %q", out, substr)
		}
	}
}

func TestWithAttrsAndWithGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := uslogs.NewUnstructuredHandler(uslogs.WithWriter(buf))
	h1 := handler.WithAttrs([]slog.Attr{slog.String("a", "val")})
	h2 := h1.WithGroup("grp")

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 1)
	record.AddAttrs(slog.String("b", "val2"))

	if err := h2.Handle(context.Background(), record); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "a=val") {
		t.Errorf("output = %q, want it to contain attribute a=val", out)
	}
	if !strings.Contains(out, "grp.b=val2") {
		t.Errorf("output = %q, want attribute b prefixed by group", out)
	}
}

func TestMaskedFields(t *testing.T) {
	buf := &bytes.Buffer{}

	handler := uslogs.NewUnstructuredHandler(
		uslogs.WithWriter(buf),
		uslogs.WithMaskedAttributes("secret", "dontshow"))

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 1)
	record.AddAttrs(slog.String("secret", "dontshow"))

	if err := handler.Handle(context.Background(), record); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "<MASKED>") {
		t.Errorf("output = %q, want masked value", out)
	}
}

func TestWithTime(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := uslogs.NewUnstructuredHandler(uslogs.WithWriter(buf), uslogs.WithTimestamp())

	tm := time.Now()
	record := slog.NewRecord(tm, slog.LevelInfo, "msg", 1)

	if err := handler.Handle(context.Background(), record); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	out := buf.String()
	timeStr := tm.UTC().Format(time.RFC3339)
	if !strings.Contains(out, timeStr) {
		t.Errorf("output = %q, want it to contain time %q", out, timeStr)
	}
}

func TestIsResponsivePool(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := uslogs.NewUnstructuredHandler(uslogs.WithWriter(buf), uslogs.WithResponsivePool())

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 1)
	record.AddAttrs(slog.String("a", "val"))

	if err := handler.Handle(context.Background(), record); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "a=val") || !strings.Contains(out, "msg") {
		t.Errorf("output = %q, want it to contain message and attribute", out)
	}
}

func TestWithPatternMasking(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := uslogs.NewUnstructuredHandler(
		uslogs.WithWriter(buf),
		uslogs.WithMaskedPatterns(uslogs.MaskPattern{Start: "Authorization:[", Delimiters: []byte{']'}}))

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "Authorization:[Bearer asdfqwerty]", 1)
	record.AddAttrs(slog.String("a", "val"))

	if err := handler.Handle(context.Background(), record); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Authorization:[*****************]") {
		t.Errorf("output = %q, want it to contain message and attribute", out)
	}
}
