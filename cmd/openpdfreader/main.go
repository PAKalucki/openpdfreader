// Package main is the entry point for OpenPDF Reader application.
package main

import (
	"os"

	"github.com/openpdfreader/openpdfreader/internal/app"
)

func main() {
	a := app.New()

	// Open file from command line argument
	if len(os.Args) > 1 {
		a.SetInitialFile(os.Args[1])
	}

	a.Run()
}
