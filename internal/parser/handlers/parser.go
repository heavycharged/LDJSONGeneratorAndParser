//go:generate easyjson -all

package handler

import (
	"bufio"
	"context"
	"io"
	"sync"

	"github.com/mailru/easyjson"
)

type LogTs int
type LogLevel string
type LogMessage string

// easyjson:json
type LogEntry struct {
	Ts      LogTs      `json:"-"`
	Level   LogLevel   `json:"level"`
	Message LogMessage `json:"-"`
}

func createWorker(lines <-chan [][]byte, results chan<- map[LogLevel]int64) {
	localResult := make(map[LogLevel]int64)

	for bag := range lines {
		for _, line := range bag {
			entry := &LogEntry{}
			if err := easyjson.Unmarshal(line, entry); err == nil {
				localResult[entry.Level]++
			}
		}
	}

	results <- localResult
}

func scanFile(ctx context.Context, input io.Reader, lines chan<- [][]byte) error {
	const (
		BATCH_SIZE       = 1024
		MAX_BUF_CAPACITY = 8 * 1024 * 1024
	)

	defer close(lines)
	buf := make([]byte, MAX_BUF_CAPACITY)

	scanner := bufio.NewScanner(input)
	scanner.Buffer(buf, MAX_BUF_CAPACITY)

	batch := make([][]byte, 0, BATCH_SIZE)

	for scanner.Scan() {
		line := make([]byte, len(scanner.Bytes()))
		copy(line, scanner.Bytes())
		batch = append(batch, line)

		if len(batch) >= BATCH_SIZE {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case lines <- batch:
				batch = make([][]byte, 0, BATCH_SIZE)
			}
		}
	}

	if len(batch) > 0 {
		lines <- batch
	}

	return scanner.Err()
}

func getResult(results chan map[LogLevel]int64) map[LogLevel]int64 {
	finalResult := make(map[LogLevel]int64)
	for res := range results {
		for level, count := range res {
			finalResult[level] += count
		}
	}
	return finalResult
}

func Parse(ctx context.Context, n int, input io.Reader) (map[LogLevel]int64, error) {
	var (
		results = make(chan map[LogLevel]int64)
		lines   = make(chan [][]byte, n*2)
	)

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			createWorker(lines, results)
		}()
	}

	if err := scanFile(ctx, input, lines); err != nil {
		return nil, err
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var finalResult = getResult(results)

	return finalResult, nil
}
