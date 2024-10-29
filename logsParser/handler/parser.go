package parser

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type LogEntry struct {
	Ts      int    `json:"-"`
	Level   string `json:"level"`
	Message string `json:"-"`
}

func Parse(ctx context.Context, n int, filePath string) (map[string]int64, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("No file")
		return nil, err
	}

	defer jsonFile.Close()

	results := make(chan map[string]int64, n)
	lines := make(chan []string, n)

	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localResult := make(map[string]int64)
			for batch := range lines {
				for _, line := range batch {
					var entry LogEntry
					if err := json.Unmarshal([]byte(line), &entry); err == nil {
						localResult[entry.Level]++
					}
				}
			}
			results <- localResult
		}()
	}

	scanner := bufio.NewScanner(jsonFile)
	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	batchSize := 10000
	var batch []string
	for scanner.Scan() {
		batch = append(batch, scanner.Text())
		if len(batch) >= batchSize {
			select {
			case <-ctx.Done():
				close(lines)
				return nil, ctx.Err()
			case lines <- batch:
				batch = nil
			}
		}
	}

	if len(batch) > 0 {
		lines <- batch
	}
	close(lines)

	go func() {
		wg.Wait()
		close(results)
	}()

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return nil, err
	}

	finalResult := make(map[string]int64)
	for res := range results {
		for level, count := range res {
			finalResult[level] += count
		}
	}

	return finalResult, nil
}
