package logs

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

func StreamLogs(services []string, stopChan chan struct{}) {
	var wg sync.WaitGroup

	for _, s := range services {
		service := s
		wg.Add(1)
		go func() {
			defer wg.Done()
			logFile := fmt.Sprintf("/tmp/%s.log", service)
			file, err := os.Open(logFile)
			if err != nil {
				fmt.Printf("[%s] Error opening log: %v\n", service, err)
				return
			}
			defer file.Close()

			// Seek to end
			file.Seek(0, os.SEEK_END)
			scanner := bufio.NewScanner(file)
			
			for {
				select {
				case <-stopChan:
					return
				default:
					if scanner.Scan() {
						fmt.Printf("[%s] %s\n", service, scanner.Text())
					} else {
						// Small sleep to avoid spinning
						// In a production app, we'd use fsnotify
						// but this is a simple TUI logger.
					}
				}
			}
		}()
	}
	wg.Wait()
}
