package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"

	"github.com/marcelohmariano/pg-bench/internal/bench"
)

const (
	dbURLFlagDesc  = "Postgres connection URL in format 'user:pass@host:port/dbname'"
	dbURLFlagName  = "db-url"
	dbURLShorthand = "u"

	sqlFileFlagDesc  = "Text file containing the SQL statements to be benchmarked"
	sqlFileFlagName  = "sql-file"
	sqlFileShorthand = "f"

	workersFlagDesc  = "Number of concurrent bench workers"
	workersFlagName  = "workers"
	workersShorthand = "w"
)

var workersDefaultValue = runtime.NumCPU()

func main() {
	command := newPgBenchCommand()
	if err := command.Execute(); err != nil {
		fail(err)
	}
}

func fail(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func newPgBenchCommand() *cobra.Command {
	var (
		dbURL   string
		sqlFile string
		workers int
	)

	cmd := &cobra.Command{
		Use:  filepath.Base(os.Args[0]),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return run(normalize(dbURL), sqlFile, workers)
		},
	}

	cmd.Flags().StringVarP(&dbURL, dbURLFlagName, dbURLShorthand, "", dbURLFlagDesc)
	cmd.Flags().StringVarP(&sqlFile, sqlFileFlagName, sqlFileShorthand, "", sqlFileFlagDesc)
	cmd.Flags().IntVarP(&workers, workersFlagName, workersShorthand, workersDefaultValue, workersFlagDesc)

	_ = cmd.MarkFlagRequired(dbURLFlagName)
	_ = cmd.MarkFlagRequired(sqlFileFlagName)

	return cmd
}

func normalize(dbURL string) string {
	if strings.Index(dbURL, "postgres://") == 0 {
		return dbURL
	}
	return fmt.Sprintf("postgres://%s", dbURL)
}

func run(dbURL string, sqlFile string, workers int) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	input, err := os.Open(sqlFile)
	if err != nil {
		return err
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return err
	}

	stats, err := bench.Run(ctx, input, bench.NewTaskRunner(pool), bench.NumWorkers(workers))
	if err != nil {
		return err
	}
	fmt.Println(stats)

	return nil
}
