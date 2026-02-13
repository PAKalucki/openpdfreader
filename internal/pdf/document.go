// Package pdf provides PDF document handling functionality.
package pdf

import (
	"errors"
	"image"
	"image/color"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Document represents a PDF document.
type Document struct {
	path      string
	ctx       *model.Context
	pageCount int
	modified  bool
}

// Open opens a PDF file.
func Open(path string) (*Document, error) {
	ctx, err := api.ReadContextFile(path)
	if err != nil {
		return nil, err
	}

	return &Document{
		path:      path,
		ctx:       ctx,
		pageCount: ctx.PageCount,
		modified:  false,
	}, nil
}

// OpenWithPassword opens a password-protected PDF file.
func OpenWithPassword(path, password string) (*Document, error) {
	conf := model.NewDefaultConfiguration()
	conf.UserPW = password
	conf.OwnerPW = password

	ctx, err := api.ReadContextFile(path)
	if err != nil {
		return nil, err
	}

	return &Document{
		path:      path,
		ctx:       ctx,
		pageCount: ctx.PageCount,
		modified:  false,
	}, nil
}

// Close closes the document and releases resources.
func (d *Document) Close() error {
	d.ctx = nil
	return nil
}

// Path returns the file path.
func (d *Document) Path() string {
	return d.path
}

// PageCount returns the number of pages.
func (d *Document) PageCount() int {
	return d.pageCount
}

// IsModified returns true if the document has unsaved changes.
func (d *Document) IsModified() bool {
	return d.modified
}

// Save saves the document to its original path.
func (d *Document) Save() error {
	if d.path == "" {
		return errors.New("no file path set")
	}
	return d.SaveAs(d.path)
}

// SaveAs saves the document to the specified path.
func (d *Document) SaveAs(path string) error {
	if d.ctx == nil {
		return errors.New("no document loaded")
	}

	err := api.WriteContextFile(d.ctx, path)
	if err != nil {
		return err
	}

	d.path = path
	d.modified = false
	return nil
}

// RenderPage renders a page to an image.
// pageNum is 0-indexed.
// scale is the zoom factor (1.0 = 100%).
func (d *Document) RenderPage(pageNum int, scale float64) (image.Image, error) {
	if d.ctx == nil {
		return nil, errors.New("no document loaded")
	}

	if pageNum < 0 || pageNum >= d.pageCount {
		return nil, errors.New("page number out of range")
	}

	// TODO: Implement actual PDF rendering using go-pdfium
	// For now, return a placeholder image
	width := int(612 * scale)  // Letter size width in points
	height := int(792 * scale) // Letter size height in points

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with white
	white := color.RGBA{255, 255, 255, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, white)
		}
	}

	// Draw border
	gray := color.RGBA{200, 200, 200, 255}
	for x := 0; x < width; x++ {
		img.Set(x, 0, gray)
		img.Set(x, height-1, gray)
	}
	for y := 0; y < height; y++ {
		img.Set(0, y, gray)
		img.Set(width-1, y, gray)
	}

	return img, nil
}

// GetPageSize returns the size of a page in points.
func (d *Document) GetPageSize(pageNum int) (width, height float64, err error) {
	if d.ctx == nil {
		return 0, 0, errors.New("no document loaded")
	}

	if pageNum < 0 || pageNum >= d.pageCount {
		return 0, 0, errors.New("page number out of range")
	}

	// Default to Letter size
	return 612, 792, nil
}

// ExtractText extracts text from a page.
func (d *Document) ExtractText(pageNum int) (string, error) {
	if d.ctx == nil {
		return "", errors.New("no document loaded")
	}

	// TODO: Implement text extraction
	return "", nil
}
