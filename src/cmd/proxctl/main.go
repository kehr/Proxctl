package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kehr/proxctl/src/internal/cli"
)

func main() {
	rt := cli.NewRuntime(os.Stdout, os.Stderr, os.Stdin)
	if err := cli.Execute(context.Background(), os.Args[1:], rt); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
