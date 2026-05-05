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
  --env <file>        Path to .env file (required)
  --secrets <kdbx>    Path to KeePassXC .kdbx database (required if env contains kp:// refs)
  --password <value>  Password to unlock the database (required if env contains kp:// refs)
  --verbose           Print detailed logging and stats
  --version           Print version and exit
  --help              Show this help message

1Password (op://) references use the op CLI's existing authentication
(desktop app integration, op signin session, or OP_SERVICE_ACCOUNT_TOKEN).

Examples:
  # KeePassXC only
  key-env run --env .env.dev --secrets ./secrets.kdbx --password 'pw' -- npm test

  # 1Password only (no --secrets/--password needed)
  key-env run --env .env.dev -- npm test`

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
