package uslogs

import (
	"io"
	"log/slog"

	"github.com/Drathveloper/uslogs/internal/logutils"
)

// MaskPattern represents a pattern that should be masked.
type MaskPattern struct {
	Start      string
	Delimiters []byte
}

// LogWriterOption represents a function that configures a log writer.
type LogWriterOption = func(w *UnstructuredHandler)

// WithWriter sets the writer to be used by the handler.
func WithWriter(writer io.Writer) LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		logWriter.writer = writer
	}
}

// WithLevel sets the minimum log level to be logged.
func WithLevel(level slog.Level) LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		logWriter.level = level
	}
}

// WithSeparator sets the separator between attributes.
func WithSeparator(separator byte) LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		logWriter.separator = separator
	}
}

// WithTimestamp adds a timestamp in RFC3339 format to the log line.
func WithTimestamp() LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		logWriter.withTime = true
	}
}

// WithMaskedAttributes masks the given attributes.
func WithMaskedAttributes(attrs ...string) LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		logWriter.maskedAttrs = attrs
	}
}

// WithMaskedPatterns masks all attributes that start with the given pattern.
func WithMaskedPatterns(patterns ...MaskPattern) LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		dict := make([]string, 0, len(patterns))
		maskPatterns := make([]logutils.MaskPattern, 0, len(patterns))
		for _, pattern := range patterns {
			dict = append(dict, pattern.Start)
			maskPatterns = append(maskPatterns, logutils.NewMaskPattern(pattern.Start, '*', pattern.Delimiters...))
		}
		logWriter.partialMasker = logutils.NewMasker(dict...)
		logWriter.partialMaskPatterns = maskPatterns
	}
}

// WithResponsivePool enables the use of a responsive pool for the log writer.
func WithResponsivePool() LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		logWriter.isResponsivePool = true
	}
}
