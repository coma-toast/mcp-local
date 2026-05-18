package logs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry is a service name and optional absolute log path (empty → /tmp/<name>.log).
type Entry struct {
	Name string
	Path string
}

func StreamLogs(entries []Entry, stopChan chan struct{}) {
	var wg sync.WaitGroup
	for _, e := range entries {
		entry := e
		wg.Add(1)
		go func() {
			defer wg.Done()
			path := entry.Path
			if path == "" {
				path = filepath.Join(os.TempDir(), entry.Name+".log")
			}
			file, err := os.Open(path)
			if err != nil {
				fmt.Printf("[%s] Error opening log %s: %v\n", entry.Name, path, err)
				return
			}
			defer file.Close()
			_, _ = file.Seek(0, os.SEEK_END)
			scanner := bufio.NewScanner(file)
			for {
				select {
				case <-stopChan:
					return
				default:
					if scanner.Scan() {
						fmt.Printf("[%s] %s\n", entry.Name, scanner.Text())
					} else {
						time.Sleep(100 * time.Millisecond)
					}
				}
			}
		}()
	}
	wg.Wait()
}
