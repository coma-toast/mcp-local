package portutil

import (
	"net"
	"testing"
	"time"
)

func TestIsRunning_ZeroPort(t *testing.T) {
	if IsRunning(0) {
		t.Error("IsRunning(0) should be false")
	}
	if IsRunning(-1) {
		t.Error("IsRunning(-1) should be false")
	}
}

func TestIsRunning_ClosedPort(t *testing.T) {
	// Pick a port that's unlikely to be in use
	// Note: This may occasionally fail if the port happens to be in use
	if IsRunning(54321) {
		t.Error("IsRunning on unused port should be false")
	}
}

func TestWaitForPort_ZeroPort(t *testing.T) {
	if !WaitForPort(0, 100*time.Millisecond) {
		t.Error("WaitForPort(0, ...) should return true immediately")
	}
	if !WaitForPort(-1, 100*time.Millisecond) {
		t.Error("WaitForPort(-1, ...) should return true immediately")
	}
}

func TestWaitForPort_Timeout(t *testing.T) {
	// Use a port that's definitely not listening
	start := time.Now()
	result := WaitForPort(54322, 100*time.Millisecond)
	elapsed := time.Since(start)

	if result {
		t.Error("WaitForPort on unused port should return false (timeout)")
	}
	if elapsed < 80*time.Millisecond {
		t.Errorf("WaitForPort should wait for timeout, elapsed=%v", elapsed)
	}
}

func TestWaitForPort_Success(t *testing.T) {
	// Start a listener on a random port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port

	// WaitForPort should succeed quickly
	if !WaitForPort(port, 2*time.Second) {
		t.Errorf("WaitForPort should succeed for listening port %d", port)
	}
}

func TestIsRunning_IPv6LoopbackOnly(t *testing.T) {
	ln, err := net.Listen("tcp", "[::1]:0")
	if err != nil {
		t.Skipf("IPv6 loopback unavailable: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	if !IsRunning(port) {
		t.Errorf("IsRunning(%d) should be true for [::1]-only listener", port)
	}
}
