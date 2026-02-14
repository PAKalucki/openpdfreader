package ui

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/openpdfreader/openpdfreader/internal/pdf"
)

// Sidebar provides page thumbnails and bookmarks.
type Sidebar struct {
	container  *fyne.Container
	list       *widget.List
	viewer     *Viewer
	document   *pdf.Document
	visible    bool
	cacheMu    sync.RWMutex
	cache      map[int]image.Image
	emptyThumb image.Image
}

const thumbnailScale = 0.35

// NewSidebar creates a new sidebar.
func NewSidebar(viewer *Viewer) *Sidebar {
	s := &Sidebar{
		viewer:     viewer,
		visible:    true,
		cache:      make(map[int]image.Image),
		emptyThumb: buildPlaceholderThumbnail(),
	}

	s.list = widget.NewList(
		func() int {
			if s.document == nil {
				return 0
			}
			return s.document.PageCount()
		},
		func() fyne.CanvasObject {
			thumb := canvas.NewImageFromImage(s.emptyThumb)
			thumb.FillMode = canvas.ImageFillContain
			thumb.ScaleMode = canvas.ImageScaleSmooth
			thumb.SetMinSize(fyne.NewSize(90, 120))

			label := widget.NewLabel("Page 000")
			label.Alignment = fyne.TextAlignCenter

			return container.NewVBox(thumb, label)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			item := obj.(*fyne.Container)
			thumb := item.Objects[0].(*canvas.Image)
			label := item.Objects[1].(*widget.Label)
			label.SetText(fmt.Sprintf("Page %d", id+1))

			pageNum := int(id)
			if cached := s.getCachedThumbnail(pageNum); cached != nil {
				thumb.Image = cached
				thumb.Refresh()
				return
			}

			if s.document == nil {
				thumb.Image = s.emptyThumb
				thumb.Refresh()
				return
			}

			rendered, err := s.document.RenderPage(pageNum, thumbnailScale)
			if err != nil || rendered == nil {
				rendered = s.emptyThumb
			}
			s.setCachedThumbnail(pageNum, rendered)

			thumb.Image = rendered
			thumb.Refresh()
		},
	)

	s.list.OnSelected = func(id widget.ListItemID) {
		if s.viewer != nil {
			s.viewer.GoToPage(int(id))
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
	s.cacheMu.Lock()
	s.cache = make(map[int]image.Image)
	s.cacheMu.Unlock()
	s.list.Refresh()
	if doc != nil && doc.PageCount() > 0 {
		s.list.Select(0)
	}
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

func (s *Sidebar) getCachedThumbnail(pageNum int) image.Image {
	s.cacheMu.RLock()
	defer s.cacheMu.RUnlock()
	return s.cache[pageNum]
}

func (s *Sidebar) setCachedThumbnail(pageNum int, img image.Image) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	s.cache[pageNum] = img
}

func buildPlaceholderThumbnail() image.Image {
	const width = 120
	const height = 160

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	bg := color.RGBA{245, 245, 245, 255}
	border := color.RGBA{200, 200, 200, 255}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bg)
		}
	}

	for x := 0; x < width; x++ {
		img.Set(x, 0, border)
		img.Set(x, height-1, border)
	}
	for y := 0; y < height; y++ {
		img.Set(0, y, border)
		img.Set(width-1, y, border)
	}

	return img
}
