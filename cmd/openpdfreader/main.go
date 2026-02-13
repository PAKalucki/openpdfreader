// Package main is the entry point for OpenPDF Reader application.
package main

import (
	"github.com/openpdfreader/openpdfreader/internal/app"
)

func main() {
	application := app.New()
	application.Run()
}
