// Package app provides the main application logic for OpenPDF Reader.
package app

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

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
	cfg := config.Load()
	fyneApp := app.NewWithID("com.openpdfreader.app")
	applyConfiguredTheme(fyneApp, cfg)

	return &App{
		fyneApp: fyneApp,
		config:  cfg,
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

func applyConfiguredTheme(fyneApp fyne.App, cfg *config.Config) {
	switch normalizeThemeName(cfg.Theme) {
	case "light":
		cfg.Theme = "light"
		fyneApp.Settings().SetTheme(theme.LightTheme())
	case "dark":
		cfg.Theme = "dark"
		fyneApp.Settings().SetTheme(theme.DarkTheme())
	default:
		cfg.Theme = "system"
		fyneApp.Settings().SetTheme(theme.DefaultTheme())
	}
}

func normalizeThemeName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "light":
		return "light"
	case "dark":
		return "dark"
	default:
		return "system"
	}
}
