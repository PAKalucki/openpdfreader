// Package app provides the main application logic for OpenPDF Reader.
package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/openpdfreader/openpdfreader/internal/config"
	"github.com/openpdfreader/openpdfreader/internal/ui"
)

// App represents the main application.
type App struct {
	fyneApp     fyne.App
	window      *ui.MainWindow
	config      *config.Config
	initialFile string
}

// New creates a new application instance.
func New() *App {
	fyneApp := app.NewWithID("com.openpdfreader.app")

	return &App{
		fyneApp: fyneApp,
		config:  config.Load(),
	}
}

// SetInitialFile sets a file to open on startup.
func (a *App) SetInitialFile(path string) {
	a.initialFile = path
}

// Run starts the application.
func (a *App) Run() {
	a.window = ui.NewMainWindow(a.fyneApp, a.config)

	if a.initialFile != "" {
		a.window.OpenFile(a.initialFile)
	}

	a.window.ShowAndRun()
}
