package uslogs

import (
	"context"
	"io"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"sync"

	"github.com/Drathveloper/uslogs/internal/logutils"
)

const (
	maskedFieldValue = "<MASKED>"
)

var levelNames = map[slog.Level][]byte{
	slog.LevelDebug: []byte("DEBUG"),
	slog.LevelInfo:  []byte("INFO"),
	slog.LevelWarn:  []byte("WARN"),
	slog.LevelError: []byte("ERROR"),
}

type UnstructuredHandler struct {
	level            slog.Level
	withTime         bool
	isResponsivePool bool
	separator        byte
	groupSeparator   byte
	group            []byte
	attrs            []byte
	maskedFields     []string
	writer           io.Writer
}

func NewUnstructuredHandler(opts ...LogWriterOption) *UnstructuredHandler {
	logWriter := &UnstructuredHandler{
		separator:      ' ',
		groupSeparator: '.',
		maskedFields:   make([]string, 0),
		writer:         os.Stdout,
		level:          slog.LevelInfo,
	}
	for _, opt := range opts {
		opt(logWriter)
	}
	return logWriter
}

func (l *UnstructuredHandler) Enabled(_ context.Context, level slog.Level) bool {
	return l.level <= level
}

func (l *UnstructuredHandler) Handle(_ context.Context, record slog.Record) error {
	attrBuf := logutils.SimplePool.Get().(*[]byte)
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
	buf := pool.Get().(*[]byte)
	b := (*buf)[:0]

	if l.withTime {
		b = logutils.AppendTimeRFC3339(b, record.Time.UTC())
		b = append(b, ' ', l.separator, ' ')
	}
	b = append(b, levelNames[record.Level]...)
	b = append(b, ' ', l.separator, ' ')
	b = append(b, record.Message...)
	b = append(b, l.attrs...)
	b = append(b, attrBytes...)
	b = append(b, '\n')

	if _, err := l.writer.Write(b); err != nil {
		logutils.SimplePool.Put(attrBuf)
		pool.Put(buf)
		return err
	}
	logutils.SimplePool.Put(attrBuf)
	pool.Put(buf)
	return nil
}

func (l *UnstructuredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return l
	}
	copiedLogWriter := copyLogWriter(l)
	b := make([]byte, 0, len(copiedLogWriter.attrs)+1024)
	b = append(b, copiedLogWriter.attrs...)
	for _, attr := range attrs {
		b = l.appendAttr(b, attr)
	}
	copiedLogWriter.attrs = b
	return copiedLogWriter
}

func (l *UnstructuredHandler) WithGroup(name string) slog.Handler {
	if len(name) == 0 {
		return l
	}
	newLogWriter := copyLogWriter(l)
	if len(newLogWriter.group) == 0 {
		newLogWriter.group = []byte(name)
		return newLogWriter
	}
	newLogWriter.group = append(newLogWriter.group, newLogWriter.groupSeparator)
	newLogWriter.group = append(newLogWriter.group, name...)
	return newLogWriter
}

func appendValue(b []byte, v slog.Value) []byte {
	switch v.Kind() {
	case slog.KindString:
		return append(b, v.String()...)
	case slog.KindInt64:
		return strconv.AppendInt(b, v.Int64(), 10)
	case slog.KindUint64:
		return strconv.AppendUint(b, v.Uint64(), 10)
	case slog.KindFloat64:
		return strconv.AppendFloat(b, v.Float64(), 'f', -1, 64)
	case slog.KindBool:
		return strconv.AppendBool(b, v.Bool())
	default:
		return append(b, v.String()...)
	}
}

func (l *UnstructuredHandler) appendAttr(b []byte, attr slog.Attr) []byte {
	b = append(b, ' ', l.separator, ' ')
	if len(l.group) != 0 {
		b = append(b, l.group...)
		b = append(b, l.groupSeparator)
	}
	b = append(b, attr.Key...)
	b = append(b, '=')
	if slices.Contains(l.maskedFields, attr.Key) {
		b = append(b, maskedFieldValue...)
	} else {
		b = appendValue(b, attr.Value)
	}
	return b
}

func copyLogWriter(logWriter *UnstructuredHandler) *UnstructuredHandler {
	return &UnstructuredHandler{
		level:            logWriter.level,
		withTime:         logWriter.withTime,
		separator:        logWriter.separator,
		groupSeparator:   logWriter.groupSeparator,
		group:            logWriter.group,
		attrs:            logWriter.attrs,
		maskedFields:     logWriter.maskedFields,
		writer:           logWriter.writer,
		isResponsivePool: logWriter.isResponsivePool,
	}
}
