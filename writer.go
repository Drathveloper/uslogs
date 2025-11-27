package uslogs

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"

	"github.com/Drathveloper/uslogs/internal/logutils"
)

const allocatedLogSize = 64 * 1024

// AsyncWriter is a writer that asynchronously writes logs to an underlying writer.
type AsyncWriter struct {
	logChan  chan *[]byte
	freeChan chan *[]byte
	writer   io.Writer
	closed   atomic.Bool
	wg       sync.WaitGroup
}

// NewAsyncWriter creates a new AsyncWriter instance.
func NewAsyncWriter(writer io.Writer, bufSize int) *AsyncWriter {
	asyncWriter := &AsyncWriter{
		logChan:  make(chan *[]byte, bufSize),
		freeChan: make(chan *[]byte, bufSize),
		writer:   writer,
		wg:       sync.WaitGroup{},
		closed:   atomic.Bool{},
	}
	for range bufSize {
		b := make([]byte, 0, allocatedLogSize)
		asyncWriter.freeChan <- &b
	}
	asyncWriter.wg.Add(1)
	go func() {
		for buf := range asyncWriter.logChan {
			numBytes, err := writer.Write(*buf)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "asyncWriter error: %v\n", err)
			}
			select {
			case asyncWriter.freeChan <- buf:
			default:
				logutils.BytesPools.GetPool(numBytes).Put(buf)
			}
		}
		asyncWriter.wg.Done()
	}()
	return asyncWriter
}

// Write writes the given input to the underlying writer.
func (w *AsyncWriter) Write(input []byte) (int, error) {
	if w.closed.Load() {
		return 0, io.ErrClosedPipe
	}
	var buf *[]byte
	select {
	case buf = <-w.freeChan:
	default:
		buf = logutils.BytesPools.GetPool(len(input)).Get().(*[]byte) //nolint:forcetypeassert
	}
	*buf = append((*buf)[:0], input...)
	w.logChan <- buf
	return len(input), nil
}

// Close closes the underlying writer.
func (w *AsyncWriter) Close() error {
	if w.closed.Swap(true) {
		return nil
	}

	close(w.logChan)

	w.wg.Wait()
	if c, ok := w.writer.(io.Closer); ok {
		return c.Close() //nolint: wrapcheck
	}

	return nil
}
