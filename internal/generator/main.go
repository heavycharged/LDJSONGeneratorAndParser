package main

import (
	generator "GeneratorAndParser/internal/generator/handlers"
	"GeneratorAndParser/internal/handlers"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	n        int
	fromDate string
	toDate   string
)

func init() {
	flag.IntVar(&n, "n", 100000, "Number of rows")
	flag.StringVar(&fromDate, "from", "", "Start date (YYYY-MM-DD)")
	flag.StringVar(&toDate, "to", "", "End date (YYYY-MM-DD)")
}

func main() {
	flag.Parse()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	startTime := time.Now()

	if fromDate == "" || toDate == "" {
		fmt.Println("Need to write --from and --to dates.")
		os.Exit(1)
	}

	fromTime, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		fmt.Printf("Error to parse from date: %v\n", err)
		os.Exit(1)
	}

	toTime, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		fmt.Printf("Error to parse to date: %v\n", err)
		os.Exit(1)
	}

	if !fromTime.Before(toTime) {
		fmt.Println("from date must be before to date.")
		os.Exit(1)
	}

	ctx := handlers.SetupContext()

	err = generator.Generate(ctx, n, fromTime, toTime)

	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	elapsedTime := time.Since(startTime)

	fmt.Printf("File generated in %d milliseconds.\n", elapsedTime.Milliseconds())
}
