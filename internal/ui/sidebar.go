package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/openpdfreader/openpdfreader/internal/pdf"
)

// Sidebar provides page thumbnails and bookmarks.
type Sidebar struct {
	container *fyne.Container
	list      *widget.List
	window    *MainWindow
	document  *pdf.Document
	visible   bool
}

// NewSidebar creates a new sidebar.
func NewSidebar(window *MainWindow) *Sidebar {
	s := &Sidebar{
		window:  window,
		visible: true,
	}

	s.list = widget.NewList(
		func() int {
			if s.document == nil {
				return 0
			}
			return s.document.PageCount()
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Page 000")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			label.SetText("Page " + intToStr(id+1))
		},
	)

	s.list.OnSelected = func(id widget.ListItemID) {
		if s.window.viewer != nil {
			s.window.viewer.GoToPage(int(id))
		}
	}

	header := widget.NewLabel("Pages")
	header.TextStyle = fyne.TextStyle{Bold: true}

	s.container = container.NewBorder(
		header,
		nil,
		nil,
		nil,
		s.list,
	)

	return s
}

// Container returns the sidebar's container.
func (s *Sidebar) Container() *fyne.Container {
	return s.container
}

// SetDocument updates the sidebar with a new document.
func (s *Sidebar) SetDocument(doc *pdf.Document) {
	s.document = doc
	s.list.Refresh()
}

// Toggle shows or hides the sidebar.
func (s *Sidebar) Toggle() {
	s.visible = !s.visible
	if s.visible {
		s.container.Show()
	} else {
		s.container.Hide()
	}
}
