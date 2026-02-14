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

// fixedSizeLayout is a layout that gives its content a fixed size.
type fixedSizeLayout struct {
	size fyne.Size
}

func (f *fixedSizeLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return f.size
}

func (f *fixedSizeLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Resize(f.size)
		o.Move(fyne.NewPos(0, 0))
	}
}

// Viewer displays PDF pages.
type Viewer struct {
	container   *fyne.Container
	scroll      *container.Scroll
	pageImage   *canvas.Image
	imageHolder *fyne.Container
	sizeLayout  *fixedSizeLayout
	document    *pdf.Document
	currentPage int
	zoom        float64
	pageLabel   *widget.Label
	zoomLabel   *widget.Label
	prevBtn     *widget.Button
	nextBtn     *widget.Button
	cachedImage image.Image
	cachedPage  int
	baseWidth   int
	baseHeight  int
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

	// Create image
	v.pageImage = canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 100, 100)))
	v.pageImage.FillMode = canvas.ImageFillContain
	v.pageImage.ScaleMode = canvas.ImageScaleSmooth

	// Use fixed size layout to control image size
	v.sizeLayout = &fixedSizeLayout{size: fyne.NewSize(100, 100)}
	v.imageHolder = container.New(v.sizeLayout, v.pageImage)

	v.scroll = container.NewScroll(v.imageHolder)

	v.container = container.NewBorder(
		nil,
		v.createPageControls(),
		nil,
		nil,
		v.scroll,
	)

	v.updateButtonStates()

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
	v.cachedPage = -1
	v.cachedImage = nil
	v.renderCurrentPage()
}

// GoToPage navigates to the specified page (0-indexed).
func (v *Viewer) GoToPage(page int) {
	if v.document == nil {
		return
	}

	pageCount := v.document.PageCount()
	if page < 0 || page >= pageCount || page == v.currentPage {
		return
	}

	v.currentPage = page
	v.renderCurrentPage()
}

// CurrentPage returns the currently displayed page (0-indexed).
func (v *Viewer) CurrentPage() int {
	return v.currentPage
}

// ZoomIn increases the zoom level.
func (v *Viewer) ZoomIn() {
	newZoom := v.zoom * 1.25
	if newZoom > 5.0 {
		newZoom = 5.0
	}
	if newZoom != v.zoom {
		v.zoom = newZoom
		v.applyZoom()
	}
}

// ZoomOut decreases the zoom level.
func (v *Viewer) ZoomOut() {
	newZoom := v.zoom / 1.25
	if newZoom < 0.25 {
		newZoom = 0.25
	}
	if newZoom != v.zoom {
		v.zoom = newZoom
		v.applyZoom()
	}
}

// FitToPage fits the page to the viewport.
func (v *Viewer) FitToPage() {
	if v.cachedImage == nil {
		return
	}

	viewport := v.scroll.Size()
	if viewport.Width <= 0 || viewport.Height <= 0 {
		return
	}

	baseWidth := float32(v.baseWidth) / 2.0
	baseHeight := float32(v.baseHeight) / 2.0
	if baseWidth <= 0 || baseHeight <= 0 {
		return
	}

	widthZoom := float64(viewport.Width / baseWidth)
	heightZoom := float64(viewport.Height / baseHeight)
	v.zoom = clampZoom(minFloat(widthZoom, heightZoom))
	v.applyZoom()
}

// FitToWidth fits the page width to the viewport.
func (v *Viewer) FitToWidth() {
	if v.cachedImage == nil {
		return
	}

	viewport := v.scroll.Size()
	if viewport.Width <= 0 {
		return
	}

	baseWidth := float32(v.baseWidth) / 2.0
	if baseWidth <= 0 {
		return
	}

	v.zoom = clampZoom(float64(viewport.Width / baseWidth))
	v.applyZoom()
}

func (v *Viewer) applyZoom() {
	if v.cachedImage == nil {
		return
	}

	// Calculate scaled size based on cached image and zoom
	// Base render is at 2.0x, so divide by 2 to get 100% size
	scaledWidth := float32(v.baseWidth) * float32(v.zoom) / 2.0
	scaledHeight := float32(v.baseHeight) * float32(v.zoom) / 2.0

	// Update the layout size and refresh
	v.sizeLayout.size = fyne.NewSize(scaledWidth, scaledHeight)
	v.imageHolder.Refresh()

	v.zoomLabel.SetText(fmt.Sprintf("%.0f%%", v.zoom*100))
}

func clampZoom(zoom float64) float64 {
	if zoom < 0.25 {
		return 0.25
	}
	if zoom > 5.0 {
		return 5.0
	}
	return zoom
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func (v *Viewer) renderCurrentPage() {
	if v.document == nil {
		v.pageLabel.SetText("No document loaded")
		v.zoomLabel.SetText("--")
		v.updateButtonStates()
		return
	}

	// Only re-render from PDF if page changed
	if v.cachedPage != v.currentPage || v.cachedImage == nil {
		img, err := v.document.RenderPage(v.currentPage, 2.0)
		if err != nil {
			v.pageLabel.SetText("Error: " + err.Error())
			return
		}
		v.cachedImage = img
		v.cachedPage = v.currentPage
		bounds := img.Bounds()
		v.baseWidth = bounds.Dx()
		v.baseHeight = bounds.Dy()
		v.pageImage.Image = img
		v.pageImage.Refresh()
	}

	v.applyZoom()
	v.scroll.ScrollToTop()

	v.pageLabel.SetText(fmt.Sprintf("Page %d of %d", v.currentPage+1, v.document.PageCount()))
	v.updateButtonStates()
}

func (v *Viewer) updateButtonStates() {
	if v.prevBtn == nil || v.nextBtn == nil {
		return
	}

	if v.document == nil {
		v.prevBtn.Disable()
		v.nextBtn.Disable()
		return
	}

	if v.currentPage <= 0 {
		v.prevBtn.Disable()
	} else {
		v.prevBtn.Enable()
	}

	if v.currentPage >= v.document.PageCount()-1 {
		v.nextBtn.Disable()
	} else {
		v.nextBtn.Enable()
	}
}

func (v *Viewer) createPageControls() *fyne.Container {
	v.prevBtn = widget.NewButton("<", func() {
		v.GoToPage(v.currentPage - 1)
	})

	v.nextBtn = widget.NewButton(">", func() {
		v.GoToPage(v.currentPage + 1)
	})

	zoomOutBtn := widget.NewButton("-", func() {
		v.ZoomOut()
	})

	zoomInBtn := widget.NewButton("+", func() {
		v.ZoomIn()
	})

	return container.NewHBox(
		v.prevBtn,
		v.pageLabel,
		v.nextBtn,
		widget.NewSeparator(),
		zoomOutBtn,
		v.zoomLabel,
		zoomInBtn,
	)
}

func intToStr(n int) string {
	return fmt.Sprintf("%d", n)
}
