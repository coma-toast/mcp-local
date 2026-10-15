package logs

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type syncBuf struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (s *syncBuf) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuf) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}

func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func appendFile(t *testing.T, path, s string) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(s); err != nil {
		t.Fatal(err)
	}
}

type run struct {
	out, errOut *syncBuf
	opens       chan struct{}
	cancel      context.CancelFunc
	done        chan error
}

// follow starts a Follower on path and returns once the goroutine is running.
func follow(t *testing.T, path string) *run {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	r := &run{out: &syncBuf{}, errOut: &syncBuf{}, opens: make(chan struct{}, 10), cancel: cancel, done: make(chan error, 1)}
	fl := Follower{Path: path, Out: r.out, Err: r.errOut, Interval: 5 * time.Millisecond, onOpen: func() { r.opens <- struct{}{} }}
	go func() { r.done <- fl.Run(ctx) }()
	t.Cleanup(func() { r.stop(t) })
	return r
}

func (r *run) waitOpen(t *testing.T) {
	t.Helper()
	select {
	case <-r.opens:
	case <-time.After(3 * time.Second):
		t.Fatal("follower never opened the file")
	}
}

func (r *run) stop(t *testing.T) {
	t.Helper()
	r.cancel()
	select {
	case err := <-r.done:
		if err != nil {
			t.Errorf("Run: %v", err)
		}
		r.done <- nil
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not stop after cancel")
	}
}

func TestFollowPrintsAppendedLinesInOrder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "svc.log")
	appendFile(t, path, "old line\n")
	r := follow(t, path)
	r.waitOpen(t)
	var want strings.Builder
	for i := range 20 {
		line := "line " + string(rune('a'+i)) + "\n"
		want.WriteString(line)
		appendFile(t, path, line)
		if i%5 == 0 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	waitFor(t, "20 lines", func() bool { return r.out.String() == want.String() })
}

func TestFollowHandlesTruncation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "svc.log")
	appendFile(t, path, "some existing content\n")
	r := follow(t, path)
	r.waitOpen(t)
	appendFile(t, path, "before truncate\n")
	waitFor(t, "first line", func() bool { return r.out.String() == "before truncate\n" })
	if err := os.WriteFile(path, []byte("after\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "line after truncation", func() bool { return r.out.String() == "before truncate\nafter\n" })
}

func TestFollowHandlesRotation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "svc.log")
	appendFile(t, path, "")
	r := follow(t, path)
	r.waitOpen(t)
	appendFile(t, path, "one\n")
	waitFor(t, "first line", func() bool { return r.out.String() == "one\n" })
	if err := os.Rename(path, path+".1"); err != nil {
		t.Fatal(err)
	}
	appendFile(t, path+".1", "late write to old file\n")
	appendFile(t, path, "new file\n")
	r.waitOpen(t)
	waitFor(t, "lines from both files", func() bool { return r.out.String() == "one\nlate write to old file\nnew file\n" })
	appendFile(t, path, "two\n")
	waitFor(t, "line appended to new file", func() bool { return strings.HasSuffix(r.out.String(), "new file\ntwo\n") })
}

func TestFollowWaitsForMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "later.log")
	r := follow(t, path)
	waitFor(t, "waiting notice", func() bool { return r.errOut.String() != "" })
	time.Sleep(30 * time.Millisecond) // several polls; the notice must not repeat
	if got, want := r.errOut.String(), "waiting for "+path+"\n"; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
	appendFile(t, path, "hello\n")
	waitFor(t, "line from new file", func() bool { return r.out.String() == "hello\n" })
}

func TestFollowCompletesPartialLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), "svc.log")
	appendFile(t, path, "")
	r := follow(t, path)
	r.waitOpen(t)
	appendFile(t, path, "par")
	time.Sleep(30 * time.Millisecond)
	if got := r.out.String(); got != "" {
		t.Fatalf("partial line printed early: %q", got)
	}
	appendFile(t, path, "tial\nnext\n")
	waitFor(t, "joined line", func() bool { return r.out.String() == "partial\nnext\n" })
}

func TestFollowStopsOnCancelWhileWaiting(t *testing.T) {
	r := follow(t, filepath.Join(t.TempDir(), "never.log"))
	waitFor(t, "waiting notice", func() bool { return r.errOut.String() != "" })
	r.stop(t)
}

func TestStreamLogsPrefixesEachService(t *testing.T) {
	old := defaultInterval
	defaultInterval = 5 * time.Millisecond
	t.Cleanup(func() { defaultInterval = old })
	dir := t.TempDir()
	a, b := filepath.Join(dir, "a.log"), filepath.Join(dir, "b.log")
	ctx, cancel := context.WithCancel(context.Background())
	out, errOut := &syncBuf{}, &syncBuf{}
	done := make(chan struct{})
	go func() {
		StreamLogs(ctx, []Entry{{Name: "a", Path: a}, {Name: "b", Path: b}}, out, errOut)
		close(done)
	}()
	waitFor(t, "both waiting notices", func() bool { return strings.Count(errOut.String(), "waiting for") == 2 })
	appendFile(t, a, "from a\n")
	appendFile(t, b, "from b\n")
	waitFor(t, "prefixed lines", func() bool {
		s := out.String()
		return strings.Contains(s, "[a] from a\n") && strings.Contains(s, "[b] from b\n")
	})
	cancel()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("StreamLogs did not stop after cancel")
	}
}
