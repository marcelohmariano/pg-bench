package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

const (
	dbURLFlagName = "db-url"
	dbURLFlagDesc = `DB connection URL in format "postgres://<user>:<pass>@<host>:<port>/<dbname>"`

	queryFlagName = "query"
	queryFlagDesc = "Text file containing the SQL commands to benchmarked"

	queryArgsFlagName = "query-args"
	queryArgsFlagDesc = "CSV file containing the args for the generated queries"

	csvDelimFlagName = "csv-delim"
	csvDelimFlagDesc = "Character used as delimiter in the CSV file"

	csvHeaderFlagName = "csv-header"
	csvHeaderFlagDesc = "Whether to treat first line as the file header"

	workersFlagName = "workers"
	workersFlagDesc = "Number of concurrent benchmark workers"
)

var (
	dbURLFlagValue string

	queryFlagValue     string
	queryArgsFlagValue string

	csvDelimFlagValue  string
	csvHeaderFlagValue bool

	workersFlagValue uint
)

func parseFlags() error {
	pflag.ErrHelp = fmt.Errorf("Showing help for %s.", os.Args[0])
	pflag.StringVar(&dbURLFlagValue, dbURLFlagName, "", dbURLFlagDesc)

	pflag.StringVar(&queryFlagValue, queryFlagName, "", queryFlagDesc)
	pflag.StringVar(&queryArgsFlagValue, queryArgsFlagName, "", queryArgsFlagDesc)

	pflag.StringVar(&csvDelimFlagValue, csvDelimFlagName, ",", csvDelimFlagDesc)
	pflag.BoolVar(&csvHeaderFlagValue, csvHeaderFlagName, true, csvHeaderFlagDesc)

	pflag.UintVar(&workersFlagValue, workersFlagName, 4, workersFlagDesc)

	pflag.Parse()
	return validateFlags()
}

func validateFlags() error {
	if err := requireFlag(dbURLFlagValue != "", dbURLFlagName); err != nil {
		return err
	}

	if err := requireFlag(queryFlagValue != "", queryFlagName); err != nil {
		return err
	}

	if err := requireFlag(queryArgsFlagValue != "", queryArgsFlagName); err != nil {
		return err
	}

	return nil
}

func requireFlag(condition bool, name string) error {
	const flagNeedsArgErrMsg = "flag needs an argument: --%s"

	if !condition {
		return fmt.Errorf(flagNeedsArgErrMsg, name)
	}

	return nil
}
