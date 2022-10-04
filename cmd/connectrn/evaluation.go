package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"

	"connectrn-evaluation/cmd/connectrn/commands"
	"connectrn-evaluation/internal/logging"
)

func main() {
	parser := flags.NewParser(nil, flags.Default)
	parser.ShortDescription = "connectRN evaluation"
	parser.LongDescription = ""

	var err error

	logging.Configure("info", "root", []string{"stdout"}, []string{"stdout"})
	logger := logging.GetLoggerUnsafe()

	_, err = parser.AddCommand(
		"server",
		"Starts the API service",
		` `,
		&commands.ServerCmd{},
	)
	if err != nil {
		logger.Fatal(err)
	}
	if _, err = parser.Parse(); err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}
