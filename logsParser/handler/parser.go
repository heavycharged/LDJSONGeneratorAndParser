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

func createWorker(wg *sync.WaitGroup, lines <-chan []string, results chan<- map[string]int64) {
	defer wg.Done()
	localResult := make(map[string]int64)
	for bag := range lines {
		for _, line := range bag {
			var entry LogEntry
			if err := json.Unmarshal([]byte(line), &entry); err == nil {
				localResult[entry.Level]++
			}
		}
	}
	results <- localResult
}

func scanFile(ctx context.Context, file *os.File, lines chan<- []string) error {
	scanner := bufio.NewScanner(file)
	maxCapacity := 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	batchSize := 100000
	batch := make([]string, batchSize)
	for scanner.Scan() {
		batch = append(batch, scanner.Text())
		if len(batch) >= batchSize {
			select {
			case <-ctx.Done():
				close(lines)
				return ctx.Err()
			case lines <- batch:
				batch = nil
			}
		}
	}

	if len(batch) > 0 {
		lines <- batch
	}
	close(lines)

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return err
	}

	return nil
}

func waitAndClose(wg *sync.WaitGroup, result chan map[string]int64) {
	wg.Wait()
	close(result)
}

func getResult(results chan map[string]int64) map[string]int64 {
	finalResult := make(map[string]int64)
	for res := range results {
		for level, count := range res {
			finalResult[level] += count
		}
	}
	return finalResult
}

func Parse(ctx context.Context, n int, filePath string) (map[string]int64, error) {
	file, err := os.Open(filePath)

	if err != nil {
		fmt.Println("No file")
		return nil, err
	}

	defer file.Close()

	results := make(chan map[string]int64, n)
	lines := make(chan []string, n)

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go createWorker(&wg, lines, results)
	}

	if err := scanFile(ctx, file, lines); err != nil {
		return nil, err
	}

	go waitAndClose(&wg, results)

	var finalResult = getResult(results)

	return finalResult, nil
}
