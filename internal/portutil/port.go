package portutil

import (
	"fmt"
	"net"
	"os/exec"
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
	out, err := exec.Command("lsof", "-t", "-iTCP:"+portStr, "-sTCP:LISTEN").Output()
	if err != nil {
		return 0
	}
	line := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	pid, _ := strconv.Atoi(line)
	return pid
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
	out, err := exec.Command("ps", "-o", "etime=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return "?"
	}
	return strings.TrimSpace(string(out))
}
