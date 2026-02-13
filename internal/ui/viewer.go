package ui

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/openpdfreader/openpdfreader/internal/pdf"
)

// Viewer displays PDF pages.
type Viewer struct {
	container   *fyne.Container
	scroll      *container.Scroll
	pageImage   *canvas.Image
	document    *pdf.Document
	currentPage int
	zoom        float64
	pageLabel   *widget.Label
	zoomLabel   *widget.Label
}

// NewViewer creates a new PDF viewer widget.
func NewViewer() *Viewer {
	v := &Viewer{
		currentPage: 0,
		zoom:        1.0,
		pageLabel:   widget.NewLabel("No document loaded"),
		zoomLabel:   widget.NewLabel("100%"),
	}

	// Placeholder image
	v.pageImage = canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 1, 1)))
	v.pageImage.FillMode = canvas.ImageFillOriginal
	v.pageImage.ScaleMode = canvas.ImageScaleFastest

	v.scroll = container.NewScroll(v.pageImage)

	v.container = container.NewBorder(
		nil,
		v.createPageControls(),
		nil,
		nil,
		v.scroll,
	)

	return v
}

// Container returns the viewer's container.
func (v *Viewer) Container() *fyne.Container {
	return v.container
}

// SetDocument sets the PDF document to display.
func (v *Viewer) SetDocument(doc *pdf.Document) {
	v.document = doc
	v.currentPage = 0
	v.renderCurrentPage()
}

// GoToPage navigates to the specified page (0-indexed).
func (v *Viewer) GoToPage(page int) {
	if v.document == nil {
		return
	}

	pageCount := v.document.PageCount()
	if page < 0 {
		page = 0
	}
	if page >= pageCount {
		page = pageCount - 1
	}

	v.currentPage = page
	v.renderCurrentPage()
}

// ZoomIn increases the zoom level.
func (v *Viewer) ZoomIn() {
	v.zoom *= 1.25
	if v.zoom > 5.0 {
		v.zoom = 5.0
	}
	v.renderCurrentPage()
}

// ZoomOut decreases the zoom level.
func (v *Viewer) ZoomOut() {
	v.zoom /= 1.25
	if v.zoom < 0.1 {
		v.zoom = 0.1
	}
	v.renderCurrentPage()
}

// FitToPage fits the page to the viewport.
func (v *Viewer) FitToPage() {
	// TODO: Calculate zoom to fit page in viewport
	v.zoom = 1.0
	v.renderCurrentPage()
}

// FitToWidth fits the page width to the viewport.
func (v *Viewer) FitToWidth() {
	// TODO: Calculate zoom to fit width in viewport
	v.zoom = 1.0
	v.renderCurrentPage()
}

func (v *Viewer) renderCurrentPage() {
	if v.document == nil {
		v.pageLabel.SetText("No document loaded")
		v.zoomLabel.SetText("--")
		return
	}

	img, err := v.document.RenderPage(v.currentPage, v.zoom)
	if err != nil {
		v.pageLabel.SetText("Error rendering page: " + err.Error())
		return
	}

	// Set the image and update its minimum size to match the rendered dimensions
	bounds := img.Bounds()
	v.pageImage.Image = img
	v.pageImage.SetMinSize(fyne.NewSize(float32(bounds.Dx()), float32(bounds.Dy())))
	v.pageImage.Refresh()

	// Scroll to top-left when page changes
	v.scroll.ScrollToTop()

	v.pageLabel.SetText(
		"Page " + intToStr(v.currentPage+1) + " of " + intToStr(v.document.PageCount()),
	)
	v.zoomLabel.SetText(fmt.Sprintf("%.0f%%", v.zoom*100))
}

func (v *Viewer) createPageControls() *fyne.Container {
	prevBtn := widget.NewButton("<", func() {
		v.GoToPage(v.currentPage - 1)
	})

	nextBtn := widget.NewButton(">", func() {
		v.GoToPage(v.currentPage + 1)
	})

	zoomOutBtn := widget.NewButton("-", func() {
		v.ZoomOut()
	})

	zoomInBtn := widget.NewButton("+", func() {
		v.ZoomIn()
	})

	return container.NewHBox(
		prevBtn,
		v.pageLabel,
		nextBtn,
		widget.NewSeparator(),
		zoomOutBtn,
		v.zoomLabel,
		zoomInBtn,
	)
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}

	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
