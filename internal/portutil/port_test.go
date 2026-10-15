package portutil

import (
	"net"
	"net/http"
	"net/http/httptest"
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

// freePort returns a port that was just free; nothing listens on it afterwards.
func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	return port
}

func TestIsRunning_ClosedPort(t *testing.T) {
	if IsRunning(freePort(t)) {
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
	result := WaitForPort(freePort(t), 100*time.Millisecond)
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

func TestFetchHealth(t *testing.T) {
	cases := []struct {
		code    int
		body    string
		healthy bool
		reason  string
	}{
		{200, `{"status":"healthy"}`, true, "healthy"},
		{200, `ok`, true, "healthy"},
		{200, `{"status":"starting","startup_message":"Opening databases…"}`, true, "starting Opening databases…"},
		{200, `{"status":"failed","embed_error":"data directory /Volumes/USB unavailable"}`, false, "failed: data directory /Volumes/USB unavailable"},
		{200, `{"status":"degraded"}`, false, "degraded"},
		{503, `down`, false, "health returned HTTP 503: down"},
	}
	for _, c := range cases {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(c.code)
			_, _ = w.Write([]byte(c.body))
		}))
		healthy, reason := FetchHealth(srv.URL)
		srv.Close()
		if healthy != c.healthy || reason != c.reason {
			t.Errorf("%d %s: got (%v, %q), want (%v, %q)", c.code, c.body, healthy, reason, c.healthy, c.reason)
		}
	}
	if ok, _ := FetchHealth(""); !ok {
		t.Error("no health_url should count as healthy")
	}
}
