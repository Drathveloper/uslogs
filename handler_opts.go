package uslogs

import (
	"io"
	"log/slog"

	"github.com/Drathveloper/uslogs/internal/logutils"
)

type MaskPattern struct {
	Start      string
	Delimiters []byte
}

type LogWriterOption = func(w *UnstructuredHandler)

func WithWriter(writer io.Writer) LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		logWriter.writer = writer
	}
}

func WithLevel(level slog.Level) LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		logWriter.level = level
	}
}

func WithSeparator(separator byte) LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		logWriter.separator = separator
	}
}

func WithTimestamp() LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		logWriter.withTime = true
	}
}

func WithMaskedAttributes(attrs ...string) LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		logWriter.maskedAttrs = attrs
	}
}

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

func WithResponsivePool() LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		logWriter.isResponsivePool = true
	}
}
