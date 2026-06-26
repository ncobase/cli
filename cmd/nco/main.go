package main

import (
	"fmt"
	"os"

	"github.com/ncobase/cli/internal/app"
)

func main() {
	if err := app.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
