package portutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func IsRunning(port int) bool {
	if port <= 0 {
		return false
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 500*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
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

// GetCommandForPID returns the executable path of the process with given PID.
func GetCommandForPID(pid int) string {
	if pid <= 0 {
		return ""
	}
	if runtime.GOOS == "darwin" {
		out, err := exec.Command("ps", "-o", "comm=", "-p", strconv.Itoa(pid)).Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(out))
	}
	out, err := exec.Command("ps", "-o", "args=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return ""
	}
	fields := strings.Fields(strings.TrimSpace(string(out)))
	if len(fields) > 0 {
		return fields[0]
	}
	return ""
}

// IsSameBinary checks if the running process command matches the expected service command/binary.
func IsSameBinary(pid int, expectedCommand string) bool {
	if pid <= 0 || expectedCommand == "" {
		return false
	}
	// Get full process command or executable
	var cmdStr string
	if runtime.GOOS == "darwin" {
		out, err := exec.Command("ps", "-o", "command=", "-p", strconv.Itoa(pid)).Output()
		if err == nil {
			cmdStr = strings.TrimSpace(string(out))
		}
	} else {
		out, err := exec.Command("ps", "-o", "args=", "-p", strconv.Itoa(pid)).Output()
		if err == nil {
			cmdStr = strings.TrimSpace(string(out))
		}
	}
	if cmdStr == "" {
		// Fallback to comm
		cmdStr = GetCommandForPID(pid)
	}

	expectedBase := filepath.Base(expectedCommand)
	actualBase := filepath.Base(cmdStr)
	if strings.Contains(cmdStr, expectedCommand) || strings.Contains(expectedCommand, cmdStr) || expectedBase == actualBase {
		return true
	}
	return false
}

// FetchHealth checks a health URL and returns (healthy, detailMessage).
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
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, fmt.Sprintf("health returned HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	// Try to parse JSON for rich status, embed_error, startup_message, data directory error
	var m map[string]interface{}
	if json.Unmarshal(body, &m) == nil {
		if st, ok := m["status"].(string); ok && st != "ok" && st != "healthy" && st != "" {
			return false, fmt.Sprintf("status: %s", st)
		}
		if errStr, ok := m["embed_error"].(string); ok && errStr != "" {
			return false, errStr
		}
		if startup, ok := m["startup_message"].(string); ok && startup != "" && !strings.Contains(strings.ToLower(startup), "ok") {
			return false, startup
		}
		if dataErr, ok := m["data_directory_error"].(string); ok && dataErr != "" {
			return false, dataErr
		}
		if dataErr, ok := m["error"].(string); ok && dataErr != "" {
			return false, dataErr
		}
	}
	return true, "healthy"
}
