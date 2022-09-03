package benchmark

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
)

var (
	ErrInvalidQueryArgs = errors.New("invalid query args")
)

type QueryArgs []any

type QueryArgsParser interface {
	Parse(args []string) (QueryArgs, error)
}

type QueryArgsParserFunc func(args []string) (QueryArgs, error)

var _ QueryArgsParser = (QueryArgsParserFunc)(nil)

func (f QueryArgsParserFunc) Parse(args []string) (QueryArgs, error) {
	return f(args)
}

type QueryArgsScanner struct {
	reader io.Reader

	csv           *csv.Reader
	csvHeader     bool
	csvHeaderRead bool
	csvDelim      rune

	args   QueryArgs
	parser QueryArgsParser

	err error
}

type QueryArgsScannerOption func(r *QueryArgsScanner)

func NewQueryArgsScanner(opts ...QueryArgsScannerOption) *QueryArgsScanner {
	s := &QueryArgsScanner{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithReader(reader io.Reader) QueryArgsScannerOption {
	return func(r *QueryArgsScanner) {
		r.reader = reader
	}
}

func WithFileReader(fileName string) QueryArgsScannerOption {
	return func(r *QueryArgsScanner) {
		var f *os.File

		f, r.err = os.Open(fileName)
		if r.err != nil {
			return
		}

		r.reader = f
	}
}

func WithCSVDelim(delim rune) QueryArgsScannerOption {
	return func(r *QueryArgsScanner) {
		r.csvDelim = delim
	}
}

func WithCSVHeader(header bool) QueryArgsScannerOption {
	return func(r *QueryArgsScanner) {
		r.csvHeader = header
	}
}

func WithQueryArgsParser(parser QueryArgsParser) QueryArgsScannerOption {
	return func(r *QueryArgsScanner) {
		r.parser = parser
	}
}

func (s *QueryArgsScanner) Scan() bool {
	if s.err != nil {
		return false
	}

	var record []string

	record, s.err = s.read()
	if s.err != nil {
		return false
	}

	s.args, s.err = s.parse(record)
	return s.err == nil
}

func (s *QueryArgsScanner) Data() QueryArgs {
	if s.err != nil {
		return nil
	}
	return s.args
}

func (s *QueryArgsScanner) Err() error {
	return s.err
}

func (s *QueryArgsScanner) Close() error {
	closer, ok := s.reader.(io.Closer)
	if ok {
		return closer.Close()
	}
	return nil
}

func (s *QueryArgsScanner) read() ([]string, error) {
	if s.csv == nil {
		s.csv = csv.NewReader(s.reader)
		s.csv.Comma = s.csvDelim
	}

	if s.csvHeader && !s.csvHeaderRead {
		_, err := s.csv.Read()
		if err != nil {
			return nil, err
		}
		s.csvHeaderRead = true
	}

	return s.csv.Read()
}

func (s *QueryArgsScanner) parse(record []string) (QueryArgs, error) {
	args := make([]any, len(record))

	if s.parser == nil {
		for i, r := range record {
			args[i] = r
		}
		return args, nil
	}

	return s.parser.Parse(record)
}
