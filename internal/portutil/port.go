package portutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func IsRunning(port int) bool {
	if port <= 0 {
		return false
	}
	// Probe both loopback stacks — many Node/Vite servers (e.g. Slidev) bind
	// [::1] only, while older services listen on 127.0.0.1.
	addrs := []string{
		fmt.Sprintf("127.0.0.1:%d", port),
		fmt.Sprintf("[::1]:%d", port),
	}
	for _, addr := range addrs {
		conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
		if err != nil {
			continue
		}
		_ = conn.Close()
		return true
	}
	return false
}

func FindPID(port int) int {
	if port <= 0 {
		return 0
	}
	portStr := strconv.Itoa(port)
	if runtime.GOOS == "darwin" {
		return findPIDLsof(portStr)
	}
	pid := findPIDLsof(portStr)
	if pid > 0 {
		return pid
	}
	return findPIDSS(portStr)
}

func findPIDLsof(portStr string) int {
	out, err := exec.Command("lsof", "-t", "-iTCP:"+portStr, "-sTCP:LISTEN").Output()
	if err != nil {
		return 0
	}
	line := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	pid, _ := strconv.Atoi(line)
	return pid
}

func findPIDSS(portStr string) int {
	out, err := exec.Command("ss", "-tlnp", "sport", ":"+portStr).Output()
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(out), "\n") {
		if idx := strings.Index(line, "pid="); idx != -1 {
			rest := line[idx+4:]
			end := strings.Index(rest, ",")
			if end == -1 {
				end = len(rest)
			}
			pidStr := strings.TrimSpace(rest[:end])
			pid, _ := strconv.Atoi(pidStr)
			if pid > 0 {
				return pid
			}
		}
	}
	return 0
}

func WaitForPort(port int, timeout time.Duration) bool {
	if port <= 0 {
		return true
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if IsRunning(port) {
			return true
		}
		time.Sleep(300 * time.Millisecond)
	}
	return false
}

func GetUptime(pid int) string {
	if pid == 0 {
		return "-"
	}
	if runtime.GOOS == "darwin" {
		out, err := exec.Command("ps", "-o", "etime=", "-p", strconv.Itoa(pid)).Output()
		if err != nil {
			return "?"
		}
		return strings.TrimSpace(string(out))
	}
	out, err := exec.Command("ps", "-o", "etimes=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return "?"
	}
	seconds, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	if seconds == 0 {
		return "0s"
	}
	d := time.Duration(seconds) * time.Second
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	secs := seconds % 60
	if days > 0 {
		return fmt.Sprintf("%d-%02d:%02d:%02d", days, hours, mins, secs)
	}
	return fmt.Sprintf("%02d:%02d:%02d", hours, mins, secs)
}

// FetchHealth GETs healthURL and reports whether the service says it is healthy, with a
// one-line reason when it isn't. It understands the JSON shape ast-context-cache serves
// (status, embed_error, startup_message) and treats any other 2xx body as healthy.
func FetchHealth(healthURL string) (bool, string) {
	if healthURL == "" {
		return true, ""
	}
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(healthURL)
	if err != nil {
		return false, fmt.Sprintf("health check failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, fmt.Sprintf("health returned HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var h struct {
		Status         string `json:"status"`
		Error          string `json:"error"`
		EmbedError     string `json:"embed_error"`
		StartupMessage string `json:"startup_message"`
	}
	if json.Unmarshal(body, &h) != nil {
		return true, "healthy"
	}
	switch strings.ToLower(h.Status) {
	case "", "ok", "healthy":
		return true, "healthy"
	case "starting":
		// Still warming up; not a fault to fix.
		return true, strings.TrimSpace("starting " + h.StartupMessage)
	}
	reason := h.Status
	for _, extra := range []string{h.Error, h.EmbedError, h.StartupMessage} {
		if extra != "" {
			reason += ": " + extra
			break
		}
	}
	return false, reason
}
