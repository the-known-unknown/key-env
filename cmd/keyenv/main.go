package main

import (
	"fmt"
	"os"

	"github.com/asimmittal/key-env/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "key-env: %v\n", err)
		os.Exit(1)
	}
}
