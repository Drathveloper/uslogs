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

type AsyncWriter struct {
	size     int
	logChan  chan *[]byte
	freeChan chan *[]byte
	writer   io.Writer
	closed   atomic.Bool
	wg       sync.WaitGroup
}

func NewAsyncWriter(writer io.Writer, bufSize int) *AsyncWriter {
	asyncWriter := &AsyncWriter{
		logChan:  make(chan *[]byte, bufSize),
		freeChan: make(chan *[]byte, bufSize),
		writer:   writer,
		wg:       sync.WaitGroup{},
		closed:   atomic.Bool{},
	}
	for i := 0; i < bufSize; i++ {
		b := make([]byte, 0, allocatedLogSize)
		asyncWriter.freeChan <- &b
	}
	asyncWriter.wg.Add(1)
	go func() {
		for buf := range asyncWriter.logChan {
			n, err := writer.Write(*buf)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "asyncWriter error: %v\n", err)
			}
			select {
			case asyncWriter.freeChan <- buf:
			default:
				logutils.BytesPools.GetPool(n).Put(buf)
			}
		}
		asyncWriter.wg.Done()
	}()
	return asyncWriter
}

func (w *AsyncWriter) Write(input []byte) (int, error) {
	if w.closed.Load() {
		return 0, io.ErrClosedPipe
	}
	var buf *[]byte
	select {
	case buf = <-w.freeChan:
	default:
		buf = logutils.BytesPools.GetPool(len(input)).Get().(*[]byte)
	}
	*buf = append((*buf)[:0], input...)
	w.logChan <- buf
	return len(input), nil
}

func (w *AsyncWriter) Close() error {
	if w.closed.Swap(true) {
		return nil
	}

	close(w.logChan)

	w.wg.Wait()
	if c, ok := w.writer.(io.Closer); ok {
		return c.Close()
	}

	return nil
}
