// Package logs tails service log files like `tail -F`.
package logs

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sync"
	"time"
)

// defaultInterval is how long a Follower sleeps at EOF before polling again.
var defaultInterval = 200 * time.Millisecond

// Entry is a service name and its log path (process.LogPath).
type Entry struct {
	Name string
	Path string
}

// Follower prints lines appended to Path until its context is cancelled.
// It starts at the end of an existing file, waits for a missing one, and
// follows truncation and rotation (rename + recreate).
type Follower struct {
	Path     string
	Prefix   string    // prepended to every printed line
	Out      io.Writer // log lines
	Err      io.Writer // "waiting for" notices; nil discards them
	Interval time.Duration
	onOpen   func() // test hook, called after each successful open
}

type tailer struct {
	Follower
	f       *os.File
	r       *bufio.Reader
	off     int64  // bytes of f consumed so far
	pending []byte // partial line still waiting for its newline
}

// Run follows the file until ctx is done (returns nil) or an I/O error occurs.
func (fl Follower) Run(ctx context.Context) error {
	t := &tailer{Follower: fl}
	if t.Interval <= 0 {
		t.Interval = defaultInterval
	}
	if t.Err == nil {
		t.Err = io.Discard
	}
	defer func() {
		if t.f != nil {
			_ = t.f.Close()
		}
	}()
	if err := t.open(ctx, true); err != nil || t.f == nil {
		return err
	}
	for {
		if err := t.drain(); err != nil {
			return err
		}
		if !t.sleep(ctx) {
			t.flushPartial()
			return nil
		}
		if err := t.checkFile(ctx); err != nil || t.f == nil {
			return err
		}
	}
}

// open opens Path, polling until it exists. Only a file that already exists on
// the first attempt is read from its end; one that appears later is all new.
func (t *tailer) open(ctx context.Context, fromEnd bool) error {
	warned := false
	for {
		f, err := os.Open(t.Path)
		if err == nil {
			var off int64
			if fromEnd && !warned {
				if off, err = f.Seek(0, io.SeekEnd); err != nil {
					_ = f.Close()
					return err
				}
			}
			t.f, t.r, t.off = f, bufio.NewReader(f), off
			if t.onOpen != nil {
				t.onOpen()
			}
			return nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		if !warned {
			fmt.Fprintf(t.Err, "%swaiting for %s\n", t.Prefix, t.Path)
			warned = true
		}
		if !t.sleep(ctx) {
			return nil
		}
	}
}

// drain prints every complete line up to EOF and keeps any partial tail.
func (t *tailer) drain() error {
	for {
		chunk, err := t.r.ReadSlice('\n')
		t.off += int64(len(chunk))
		t.pending = append(t.pending, chunk...)
		switch {
		case err == nil:
			t.write(t.pending)
			t.pending = t.pending[:0]
		case errors.Is(err, bufio.ErrBufferFull):
		case err == io.EOF:
			return nil
		default:
			return fmt.Errorf("read %s: %w", t.Path, err)
		}
	}
}

// checkFile reopens Path when it was replaced and rewinds when it was truncated.
// A missing Path (mid-rotation) keeps the current file open.
func (t *tailer) checkFile(ctx context.Context) error {
	cur, err := t.f.Stat()
	if err != nil {
		return err
	}
	fi, err := os.Stat(t.Path)
	switch {
	case err != nil:
		return nil
	case !os.SameFile(cur, fi):
		if err := t.drain(); err != nil {
			return err
		}
		t.flushPartial()
		_ = t.f.Close()
		t.f = nil
		return t.open(ctx, false)
	case fi.Size() < t.off:
		t.flushPartial()
		if _, err := t.f.Seek(0, io.SeekStart); err != nil {
			return err
		}
		t.r.Reset(t.f)
		t.off = 0
	}
	return nil
}

func (t *tailer) flushPartial() {
	if len(t.pending) > 0 {
		t.write(append(t.pending, '\n'))
		t.pending = t.pending[:0]
	}
}

// write emits one line in a single Write so concurrent followers don't interleave.
func (t *tailer) write(line []byte) {
	_, _ = t.Out.Write(append([]byte(t.Prefix), line...))
}

func (t *tailer) sleep(ctx context.Context) bool {
	timer := time.NewTimer(t.Interval)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

type syncWriter struct {
	mu *sync.Mutex
	w  io.Writer
}

func (s syncWriter) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.w.Write(p)
}

// StreamLogs follows every entry's log with a "[name] " prefix until ctx is done.
func StreamLogs(ctx context.Context, entries []Entry, out, errOut io.Writer) {
	mu := &sync.Mutex{}
	out, errOut = syncWriter{mu, out}, syncWriter{mu, errOut}
	var wg sync.WaitGroup
	for _, e := range entries {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fl := Follower{Path: e.Path, Prefix: "[" + e.Name + "] ", Out: out, Err: errOut}
			if err := fl.Run(ctx); err != nil {
				fmt.Fprintf(errOut, "[%s] %v\n", e.Name, err)
			}
		}()
	}
	wg.Wait()
}
