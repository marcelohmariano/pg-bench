package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/marcelohmariano/timescaledb-benchmark/internal/benchmark"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if err := parseFlags(); err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()

	pool, err := pgxpool.Connect(ctx, dbURLFlagValue)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(queryFlagValue)
	if err != nil {
		return err
	}
	query := string(content)

	args := benchmark.NewQueryArgsScanner(
		benchmark.WithFileReader(queryArgsFlagValue),
		benchmark.WithCSVHeader(csvHeaderFlagValue),
		benchmark.WithCSVDelim(rune(csvDelimFlagValue[0])),
		benchmark.WithQueryArgsParser(benchmark.QueryArgsParserFunc(parseQueryArgs)),
	)
	defer func() { _ = args.Close() }()

	runner := benchmark.NewRunner(pool, args, benchmark.WithWorkers(workersFlagValue))
	summary := runner.Run(ctx, query)
	printSummary(summary)

	return nil
}

func parseQueryArgs(args []string) (benchmark.QueryArgs, error) {
	const timeLayout = "2006-01-02 15:04:05"

	if len(args) != 3 {
		return nil, benchmark.ErrInvalidQueryArgs
	}

	hostname := args[0]

	startTime, err := time.Parse(timeLayout, args[1])
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(timeLayout, args[2])
	if err != nil {
		return nil, err
	}

	return benchmark.QueryArgs{hostname, startTime, endTime}, nil
}

func printSummary(summary benchmark.Summary) {
	fmt.Printf("Queries processed: %v\n", summary.NumberOfQueries())
	fmt.Printf("Queries with success: %v\n", summary.NumberOfSuccesses())
	fmt.Printf("Queries with error: %v\n\n", summary.NumberOfErrors())

	fmt.Printf("Minimum query time: %v\n", summary.MinQueryTime())
	fmt.Printf("Maximum query time: %v\n\n", summary.MaxQueryTime())

	fmt.Printf("Median query time: %v\n", summary.MedianQueryTime())
	fmt.Printf("Average query time: %v\n\n", summary.AvgQueryTime())

	fmt.Printf("Overall processing time: %v\n", summary.OverallQueryTime())
}
