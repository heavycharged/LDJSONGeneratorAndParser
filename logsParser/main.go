package main

import (
	"context"
	"flag"
	"fmt"
	parser "logsParser/handler"
	"runtime"
	"time"
)

var levels = []string{"trace", "debug", "info", "warn", "error"}

var (
	n        int
	filePath string
)

func init() {
	flag.IntVar(&n, "n", runtime.NumCPU(), "Number of streams")
	flag.StringVar(&filePath, "filePath", "logs.ldjson", "Name of file")
}

func main() {
	flag.Parse()

	startTime := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result, err := parser.Parse(ctx, n, filePath)

	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Println("Log's level statistics:")
	for _, v := range levels {
		fmt.Printf("%s: %d\n", v, result[v])
	}

	elapsedTime := time.Since(startTime)

	fmt.Printf("File was parsed in %d milliseconds.\n", elapsedTime.Milliseconds())
}
