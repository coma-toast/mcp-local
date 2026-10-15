package process

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
)

// Stop timing; variables so tests can shorten them.
var (
	TermGrace = 3 * time.Second
	KillGrace = 2 * time.Second
	pollEvery = 100 * time.Millisecond
)

// Alive reports whether pid exists (a negative pid means process group -pid). EPERM
// means it exists but belongs to another user.
func Alive(pid int) bool {
	if pid == 0 {
		return false
	}
	err := syscall.Kill(pid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}

// GroupOf returns pid's process group, or 0 when pid is gone.
func GroupOf(pid int) int {
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		return 0
	}
	return pgid
}

// SignalTarget is -pid when pid leads its own process group, as services started by
// mcp-local (Setsid) do, so signals also reach wrapper children; otherwise pid.
func SignalTarget(pid int) int {
	if pid > 1 && GroupOf(pid) == pid {
		return -pid
	}
	return pid
}

// Terminate sends SIGTERM to target (a PID, or a process group when negative) and waits
// until it is gone and done returns nil, escalating to SIGKILL after TermGrace. killed
// reports the escalation. The error is non-nil unless both conditions were reached.
func Terminate(target int, done func() error) (killed bool, err error) {
	abs := target
	if abs < 0 {
		abs = -abs
	}
	if abs <= 1 || abs == os.Getpid() || abs == syscall.Getpgrp() {
		return false, fmt.Errorf("refusing to signal %d", target)
	}
	if done == nil {
		done = func() error { return nil }
	}
	finished := func() error {
		if Alive(target) {
			return fmt.Errorf("pid %d still running after SIGKILL", abs)
		}
		return done()
	}
	if err := signal(target, syscall.SIGTERM, "SIGTERM"); err != nil {
		return false, err
	}
	if waitFor(finished, TermGrace) == nil {
		return false, nil
	}
	if err := signal(target, syscall.SIGKILL, "SIGKILL"); err != nil {
		return true, err
	}
	return true, waitFor(finished, KillGrace)
}

func signal(target int, sig syscall.Signal, name string) error {
	if err := syscall.Kill(target, sig); err != nil && !errors.Is(err, syscall.ESRCH) {
		return fmt.Errorf("%s %d: %w", name, target, err)
	}
	return nil
}

func waitFor(check func() error, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		err := check()
		if err == nil || time.Now().After(deadline) {
			return err
		}
		time.Sleep(pollEvery)
	}
}
