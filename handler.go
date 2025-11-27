package uslogs

import (
	"context"
	"io"
	"log/slog"
	"os"
	"slices"
	"sync"

	"github.com/Drathveloper/uslogs/internal/logutils"
)

const (
	maskedFieldValue = "<MASKED>"
)

//nolint:gochecknoglobals
var levelNames = map[slog.Level][]byte{
	slog.LevelDebug: []byte("DEBUG"),
	slog.LevelInfo:  []byte("INFO"),
	slog.LevelWarn:  []byte("WARN"),
	slog.LevelError: []byte("ERROR"),
}

// UnstructuredHandler writes log lines in plain text format.
type UnstructuredHandler struct {
	writer              io.Writer
	partialMasker       *logutils.Masker
	group               []byte
	attrs               []byte
	maskedAttrs         []string
	partialMaskPatterns []logutils.MaskPattern
	level               slog.Level
	withTime            bool
	isResponsivePool    bool
	separator           byte
	groupSeparator      byte
}

// NewUnstructuredHandler creates a new UnstructuredHandler instance.
func NewUnstructuredHandler(opts ...LogWriterOption) *UnstructuredHandler {
	//nolint:exhaustruct
	logWriter := &UnstructuredHandler{
		separator:      ' ',
		groupSeparator: '.',
		maskedAttrs:    make([]string, 0),
		writer:         os.Stdout,
		level:          slog.LevelInfo,
	}
	for _, opt := range opts {
		opt(logWriter)
	}
	return logWriter
}

// Enabled returns true if the log level is greater than or equal to the configured level.
func (l *UnstructuredHandler) Enabled(_ context.Context, level slog.Level) bool {
	return l.level <= level
}

// Handle writes the log line to the writer.
func (l *UnstructuredHandler) Handle(_ context.Context, record slog.Record) error {
	attrBuf := logutils.SimplePool.Get().(*[]byte) //nolint:forcetypeassert
	attrBytes := (*attrBuf)[:0]
	record.Attrs(func(attr slog.Attr) bool {
		attrBytes = l.appendAttr(attrBytes, attr)
		return true
	})
	var pool *sync.Pool
	if l.isResponsivePool {
		pool = logutils.BytesPools.GetPool(len(record.Message) + len(l.attrs) + len(attrBytes))
	} else {
		pool = logutils.SimplePool
	}
	buf := pool.Get().(*[]byte) //nolint:forcetypeassert
	bytes := (*buf)[:0]

	if l.withTime {
		bytes = logutils.AppendTimeRFC3339(bytes, record.Time.UTC())
		bytes = logutils.AppendSeparator(bytes, l.separator)
	}
	bytes = append(bytes, levelNames[record.Level]...)
	bytes = logutils.AppendSeparator(bytes, l.separator)
	bytes = append(bytes, record.Message...)
	bytes = append(bytes, l.attrs...)
	bytes = append(bytes, attrBytes...)
	bytes = append(bytes, '\n')

	if l.partialMasker != nil && len(l.partialMaskPatterns) > 0 {
		bytes = l.partialMasker.Mask(bytes, l.partialMaskPatterns)
	}

	if _, err := l.writer.Write(bytes); err != nil {
		logutils.PutPool(logutils.SimplePool, attrBuf)
		logutils.PutPool(pool, buf)
		return err //nolint:wrapcheck
	}
	logutils.PutPool(logutils.SimplePool, attrBuf)
	logutils.PutPool(pool, buf)
	return nil
}

// WithAttrs adds attributes to the log line.
func (l *UnstructuredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return l
	}
	clonedLogWriter := l.clone()
	b := make([]byte, 0, len(clonedLogWriter.attrs)+1024) //nolint:mnd
	b = append(b, clonedLogWriter.attrs...)
	for _, attr := range attrs {
		b = l.appendAttr(b, attr)
	}
	clonedLogWriter.attrs = b
	return clonedLogWriter
}

// WithGroup adds a group name to the log line.
func (l *UnstructuredHandler) WithGroup(name string) slog.Handler {
	if len(name) == 0 {
		return l
	}
	clonedLogWriter := l.clone()
	if len(clonedLogWriter.group) == 0 {
		clonedLogWriter.group = []byte(name)
		return clonedLogWriter
	}
	clonedLogWriter.group = append(clonedLogWriter.group, clonedLogWriter.groupSeparator)
	clonedLogWriter.group = append(clonedLogWriter.group, name...)
	return clonedLogWriter
}

func (l *UnstructuredHandler) clone() *UnstructuredHandler {
	if l == nil {
		return nil
	}
	clone := *l
	return &clone
}

func (l *UnstructuredHandler) appendAttr(input []byte, attr slog.Attr) []byte {
	input = logutils.AppendSeparator(input, l.separator)
	if len(l.group) != 0 {
		input = append(input, l.group...)
		input = append(input, l.groupSeparator)
	}
	input = append(input, attr.Key...)
	input = append(input, '=')
	if slices.Contains(l.maskedAttrs, attr.Key) {
		input = append(input, maskedFieldValue...)
	} else {
		input = logutils.AppendValue(input, attr.Value)
	}
	return input
}
