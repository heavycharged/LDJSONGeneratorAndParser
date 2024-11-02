package handler

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var levels = []string{"trace", "debug", "info", "warn", "error"}
var outputFile = "logs.ldjson"

func isOutputRedirected() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func Generate(ctx context.Context, n int, fromTime, toTime time.Time) error {
	output := os.Stdout
	if !isOutputRedirected() {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("Error creating file logs.ldjson: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		output = file
	}

	delta := toTime.Sub(fromTime).Milliseconds()

	src := rand.NewSource(time.Now().UnixMilli())
	r := rand.New(src)

	writer := bufio.NewWriterSize(output, 64*1024*1024)
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
