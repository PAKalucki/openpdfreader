package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Toolbar provides quick access to common actions.
type Toolbar struct {
	container *fyne.Container
	window    *MainWindow
}

// NewToolbar creates a new toolbar.
func NewToolbar(window *MainWindow) *Toolbar {
	t := &Toolbar{
		window: window,
	}

	openBtn := widget.NewButton("Open", window.onOpenFile)
	saveBtn := widget.NewButton("Save", window.onSave)
	printBtn := widget.NewButton("Print", window.onPrint)

	sep1 := widget.NewSeparator()

	zoomInBtn := widget.NewButton("+", window.onZoomIn)
	zoomOutBtn := widget.NewButton("-", window.onZoomOut)
	fitPageBtn := widget.NewButton("Fit", window.onFitToPage)

	sep2 := widget.NewSeparator()

	prevPageBtn := widget.NewButton("◀", func() {
		if window.viewer != nil {
			window.viewer.GoToPage(window.viewer.currentPage - 1)
		}
	})

	nextPageBtn := widget.NewButton("▶", func() {
		if window.viewer != nil {
			window.viewer.GoToPage(window.viewer.currentPage + 1)
		}
	})

	t.container = container.NewHBox(
		openBtn,
		saveBtn,
		printBtn,
		sep1,
		zoomOutBtn,
		zoomInBtn,
		fitPageBtn,
		sep2,
		prevPageBtn,
		nextPageBtn,
	)

	return t
}

// Container returns the toolbar's container.
func (t *Toolbar) Container() *fyne.Container {
	return t.container
}
