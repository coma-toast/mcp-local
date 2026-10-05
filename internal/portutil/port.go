package portutil

import (
	"fmt"
	"net"
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
