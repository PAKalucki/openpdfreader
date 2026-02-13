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
	container    *fyne.Container
	scroll       *container.Scroll
	pageImage    *canvas.Image
	document     *pdf.Document
	currentPage  int
	zoom         float64
	pageLabel    *widget.Label
	zoomLabel    *widget.Label
	cachedImage  image.Image // cached render at base DPI
	cachedPage   int         // which page is cached
	baseWidth    int         // original image width
	baseHeight   int         // original image height
}

// NewViewer creates a new PDF viewer widget.
func NewViewer() *Viewer {
	v := &Viewer{
		currentPage: 0,
		zoom:        1.0,
		cachedPage:  -1,
		pageLabel:   widget.NewLabel("No document loaded"),
		zoomLabel:   widget.NewLabel("100%"),
	}

	// Placeholder image
	v.pageImage = canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 1, 1)))
	v.pageImage.FillMode = canvas.ImageFillOriginal
	v.pageImage.ScaleMode = canvas.ImageScaleSmooth

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
	v.cachedPage = -1 // invalidate cache
	v.cachedImage = nil
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

	if page != v.currentPage {
		v.currentPage = page
		v.renderCurrentPage()
		v.scroll.ScrollToTop()
	}
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

	// Only re-render from PDF if page changed
	if v.cachedPage != v.currentPage || v.cachedImage == nil {
		// Render at high DPI (2.0 = 144 DPI) for quality when zooming
		img, err := v.document.RenderPage(v.currentPage, 2.0)
		if err != nil {
			v.pageLabel.SetText("Error rendering page: " + err.Error())
			return
		}
		v.cachedImage = img
		v.cachedPage = v.currentPage
		bounds := img.Bounds()
		v.baseWidth = bounds.Dx()
		v.baseHeight = bounds.Dy()
		v.pageImage.Image = img
	}

	// Apply zoom by scaling the display size (not re-rendering)
	// Base render is at 2.0x, so divide by 2 to get 100% size, then multiply by zoom
	scaledWidth := float32(v.baseWidth) * float32(v.zoom) / 2.0
	scaledHeight := float32(v.baseHeight) * float32(v.zoom) / 2.0
	v.pageImage.SetMinSize(fyne.NewSize(scaledWidth, scaledHeight))
	v.pageImage.Refresh()

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
