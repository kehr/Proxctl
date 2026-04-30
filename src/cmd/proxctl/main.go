package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kehr/proxctl/src/internal/app"
)

func main() {
	a := app.New(os.Stdout, os.Stderr, os.Stdin)
	if err := a.Run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
