package uslog_test

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Drathveloper/uslogs"
)

func TestAsyncWriter_BasicWrite(t *testing.T) {
	var buf bytes.Buffer
	w := uslog.NewAsyncWriter(&buf, 10)

	msg := []byte("hello\n")

	n, err := w.Write(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(msg) {
		t.Fatalf("written %d bytes, expected %d", n, len(msg))
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("close error: %v", err)
	}

	if buf.String() != "hello\n" {
		t.Fatalf("got %q, expected %q", buf.String(), "hello\n")
	}
}

func TestAsyncWriter_ConcurrentWrites(t *testing.T) {
	var buf bytes.Buffer
	w := uslog.NewAsyncWriter(&buf, 1000)

	var wg sync.WaitGroup
	total := 1000

	for i := 0; i < total; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			msg := []byte(fmt.Sprintf("msg-%d\n", i))

			_, err := w.Write(msg)
			if err != nil {
				t.Errorf("write failed: %v", err)
			}
		}(i)
	}

	wg.Wait()
	_ = w.Close()

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != total {
		t.Fatalf("got %d lines, expected %d", len(lines), total)
	}
}

func TestAsyncWriter_CopiesInput(t *testing.T) {
	var buf bytes.Buffer
	w := uslog.NewAsyncWriter(&buf, 10)

	data := []byte("first\n")
	_, _ = w.Write(data)

	copy(data, "XXXXX")

	_, _ = w.Write([]byte("second\n"))

	_ = w.Close()

	result := buf.String()

	if strings.Contains(result, "XXXXX") {
		t.Fatal("buffer was not copied; data was corrupted")
	}

	if !strings.Contains(result, "first\n") {
		t.Fatal("first message missing or corrupted")
	}
}

func TestAsyncWriter_WriteAfterClose(t *testing.T) {
	var buf bytes.Buffer
	w := uslog.NewAsyncWriter(&buf, 2)

	_ = w.Close()

	_, err := w.Write([]byte("nope"))
	if err == nil {
		t.Fatal("expected error after close, got nil")
	}
}

type slowWriter struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *slowWriter) Write(p []byte) (int, error) {
	time.Sleep(5 * time.Millisecond)
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *slowWriter) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

func TestAsyncWriter_CloseWaitsForDrain(t *testing.T) {
	sw := &slowWriter{}
	w := uslog.NewAsyncWriter(sw, 1)

	total := 200

	for i := 0; i < total; i++ {
		_, _ = w.Write([]byte("x\n"))
	}

	_ = w.Close()

	if len(strings.Split(strings.TrimSpace(sw.String()), "\n")) != total {
		t.Fatalf("not all messages were written")
	}
}

func TestAsyncWriter_HTTPServerShutdown(t *testing.T) {
	var buf bytes.Buffer
	w := uslog.NewAsyncWriter(&buf, 1000)

	handler := func(wr http.ResponseWriter, r *http.Request) {
		for i := 0; i < 10; i++ {
			_, err := w.Write([]byte("log line\n"))
			if err != nil {
				t.Errorf("write failed: %v", err)
			}
			time.Sleep(5 * time.Millisecond)
		}
		wr.WriteHeader(http.StatusOK)
	}

	srv := &http.Server{
		Addr:    "127.0.0.1:0",
		Handler: http.HandlerFunc(handler),
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}

	go srv.Serve(ln)

	client := &http.Client{}

	for i := 0; i < 20; i++ {
		go client.Get("http://" + ln.Addr().String())
	}

	time.Sleep(20 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}

	if err = w.Close(); err != nil {
		t.Fatal(err)
	}

	lines := strings.Count(buf.String(), "log line\n")

	if lines != 200 {
		t.Fatalf("expected 200 lines, got %d", lines)
	}
}

func TestAsyncWriter_WithCanceledContext(t *testing.T) {
	var buf bytes.Buffer
	w := uslog.NewAsyncWriter(&buf, 100)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for i := 0; i < 100; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				_, _ = w.Write([]byte("log\n"))
			}
		}
	}()

	time.Sleep(10 * time.Millisecond)
	cancel()

	_ = w.Close()

	if buf.Len() == 0 {
		t.Fatal("no logs written")
	}
}
