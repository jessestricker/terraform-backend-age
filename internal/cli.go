package internal

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type cliOptions struct {
	port      uint
	keyFile   string
	stateFile string
}

var defaultCliOptions = cliOptions{
	port:      4321,
	keyFile:   "state-key.txt",
	stateFile: "terraform.tfstate.age",
}

const usageHeader = `Usage:
  terraform-backend-age [OPTIONS] [STATE_FILE]

Arguments:
  STATE_FILE
    The encrypted state file to serve. (default %q)

Options:
`

type cliEnvironment struct {
	key string
}

const envNamePrefix = "TF_BACKEND_AGE_"

const usageFooter = `
Environment:
  %sKEY
    The contents of the key file, takes precedence over -key-file.
`

func parseCliOptions(args []string) cliOptions {
	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, usageHeader, defaultCliOptions.stateFile)
		flagSet.PrintDefaults()
		fmt.Fprintf(os.Stderr, usageFooter, envNamePrefix)
	}

	var options cliOptions
	flagSet.UintVar(&options.port, "port", defaultCliOptions.port, "The port to listen on.")
	flagSet.StringVar(&options.keyFile, "key-file", defaultCliOptions.keyFile, "The key file to use for encryption and decryption.")

	flagSet.Parse(args[1:])

	if flagSet.NArg() > 0 {
		options.stateFile = flagSet.Arg(0)
	} else {
		options.stateFile = defaultCliOptions.stateFile
	}

	return options
}

func parseEnvironment() cliEnvironment {
	return cliEnvironment{
		key: os.Getenv(envNamePrefix + "KEY"),
	}
}

func Main() int {
	var err error

	options := parseCliOptions(os.Args)
	environment := parseEnvironment()

	// load key file
	var keyFile *keyFile
	if environment.key != "" {
		keyFile, err = parseKeyFile(strings.NewReader(environment.key))
		if err != nil {
			slog.Error("failed to parse key from environment variable", "error", err)
			return 1
		}
	} else {
		keyFile, err = loadKeyFile(options.keyFile)
		if err != nil {
			slog.Error("failed to load key file", "error", err)
			return 1
		}
	}

	srv := server{
		port:          options.port,
		keyFile:       keyFile,
		stateFilePath: options.stateFile,
	}
	err = srv.listenAndServe()
	if err != nil {
		slog.Error("server error", "error", err)
		return 1
	}

	return 0
}
