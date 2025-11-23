package uslog

import (
	"io"
	"log/slog"
)

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

func WithMaskedFields(fields ...string) LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		logWriter.maskedFields = fields
	}
}

func WithResponsivePool() LogWriterOption {
	return func(logWriter *UnstructuredHandler) {
		logWriter.isResponsivePool = true
	}
}
