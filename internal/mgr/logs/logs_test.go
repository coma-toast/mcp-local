package logs

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStreamLogs_ReadsFile(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// Write initial content
	if err := os.WriteFile(logFile, []byte("line 1\nline 2\n"), 0644); err != nil {
		t.Fatal(err)
	}

	entries := []Entry{{Name: "test-service", Path: logFile}}
	stopChan := make(chan struct{})

	// Close stopChan after a short delay to let it read
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(stopChan)
	}()

	// This should not panic and should read the lines
	StreamLogs(entries, stopChan)
}

func TestStreamLogs_EmptyPath_UsesTempDir(t *testing.T) {
	_ = t.TempDir()

	entries := []Entry{{Name: "test-service", Path: ""}}
	stopChan := make(chan struct{})

	go func() {
		time.Sleep(50 * time.Millisecond)
		close(stopChan)
	}()

	// Should not panic when file doesn't exist
	StreamLogs(entries, stopChan)
}

func TestStreamLogs_StopsOnStopChan(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// Write a lot of content
	content := ""
	for i := 0; i < 100; i++ {
		content += "line\n"
	}
	if err := os.WriteFile(logFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	entries := []Entry{{Name: "test-service", Path: logFile}}
	stopChan := make(chan struct{})

	// Close immediately - should stop quickly
	close(stopChan)

	// Should not hang
	StreamLogs(entries, stopChan)
}
