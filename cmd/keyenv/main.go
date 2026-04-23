package main

import (
	"fmt"
	"os"

	"github.com/asimmittal/key-env/internal/cli"
)

const name = "key-env"

var version = "dev"

const usage = `Usage: key-env <command> [flags] -- <child command>

Commands:
  run    Load env vars, resolve secrets, and run a child command

Flags:
  --env <file>        Path to .env file
  --secrets <kdbx>    Path to KeePassXC .kdbx database
  --password <value>  Password to unlock the database
  --version           Print version and exit
  --help              Show this help message

Example:
  key-env run --env .env.dev --secrets ./secrets.kdbx --password 'pw' -- npm test`

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version":
			fmt.Println(name + " " + version)
			return
		case "--help", "-h":
			fmt.Println(usage)
			return
		}
	}
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		os.Exit(1)
	}
}
