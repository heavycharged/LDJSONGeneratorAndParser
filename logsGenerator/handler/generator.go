package generator

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var outputFile = "logs.ldjson"

var levels = []string{"trace", "debug", "info", "warn", "error"}

func Generate(ctx context.Context, n int, fromTime, toTime time.Time) error {
	delta := toTime.Sub(fromTime).Milliseconds()

	src := rand.NewSource(time.Now().UnixMilli())
	r := rand.New(src)

	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating file logs.ldjson: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, 64*1024*1024)
	defer writer.Flush()

	for i := 0; i < n; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			ts := fromTime.Add(time.Duration(delta)).UnixMilli() + int64(i)
			level := levels[r.Intn(len(levels))]
			message := fmt.Sprintf("Message #%d", i+1)

			logEntry := fmt.Sprintf(`{"ts": %d, "level": "%s", "message": "%s"}`, ts, level, message)
			fmt.Fprintln(writer, logEntry)
		}
	}
	return nil
}
