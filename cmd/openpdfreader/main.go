// Package main is the entry point for OpenPDF Reader application.
package main

import (
	"fmt"
	"os"

	"github.com/openpdfreader/openpdfreader/internal/app"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--cli" || os.Args[1] == "cli") {
		if err := app.RunCLI(os.Args[2:], os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		return
	}

	a := app.New()

	// Open file from command line argument
	if len(os.Args) > 1 {
		a.SetInitialFile(os.Args[1])
	}

	a.Run()
}
