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

var levelNames = map[slog.Level][]byte{
	slog.LevelDebug: []byte("DEBUG"),
	slog.LevelInfo:  []byte("INFO"),
	slog.LevelWarn:  []byte("WARN"),
	slog.LevelError: []byte("ERROR"),
}

type UnstructuredHandler struct {
	level               slog.Level
	withTime            bool
	isResponsivePool    bool
	separator           byte
	groupSeparator      byte
	group               []byte
	attrs               []byte
	maskedAttrs         []string
	partialMaskPatterns []logutils.MaskPattern
	partialMasker       *logutils.Masker
	writer              io.Writer
}

func NewUnstructuredHandler(opts ...LogWriterOption) *UnstructuredHandler {
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
		b = logutils.AppendSeparator(b, l.separator)
	}
	b = append(b, levelNames[record.Level]...)
	b = logutils.AppendSeparator(b, l.separator)
	b = append(b, record.Message...)
	b = append(b, l.attrs...)
	b = append(b, attrBytes...)
	b = append(b, '\n')

	if l.partialMasker != nil && len(l.partialMaskPatterns) > 0 {
		b = l.partialMasker.Mask(b, l.partialMaskPatterns)
	}

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
	clonedLogWriter := l.clone()
	b := make([]byte, 0, len(clonedLogWriter.attrs)+1024)
	b = append(b, clonedLogWriter.attrs...)
	for _, attr := range attrs {
		b = l.appendAttr(b, attr)
	}
	clonedLogWriter.attrs = b
	return clonedLogWriter
}

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

func (l *UnstructuredHandler) appendAttr(b []byte, attr slog.Attr) []byte {
	b = logutils.AppendSeparator(b, l.separator)
	if len(l.group) != 0 {
		b = append(b, l.group...)
		b = append(b, l.groupSeparator)
	}
	b = append(b, attr.Key...)
	b = append(b, '=')
	if slices.Contains(l.maskedAttrs, attr.Key) {
		b = append(b, maskedFieldValue...)
	} else {
		b = logutils.AppendValue(b, attr.Value)
	}
	return b
}
