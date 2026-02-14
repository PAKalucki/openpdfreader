package pdf

import (
	"errors"

	pdfcolor "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// Redactor provides basic redaction operations.
// Note: this implementation applies a visual redaction overlay annotation.
type Redactor struct{}

// NewRedactor creates a redactor.
func NewRedactor() *Redactor {
	return &Redactor{}
}

// ApplyVisualRedaction adds a filled black rectangle annotation to the selected page.
func (r *Redactor) ApplyVisualRedaction(inputPath, outputPath string, pageNum int, reason string) error {
	if inputPath == "" {
		return errors.New("input path is required")
	}
	if pageNum < 0 {
		return errors.New("page number out of range")
	}

	ann := model.NewSquareAnnotation(
		*types.NewRectangle(150, 430, 460, 660),
		reason,
		nextAnnotationID("redact"),
		"",
		0,
		&pdfcolor.Black,
		"OpenPDF Reader",
		nil,
		nil,
		"",
		"Visual redaction",
		&pdfcolor.Black,
		0,
		0,
		0,
		0,
		0,
		model.BSSolid,
		false,
		0,
	)

	return addAnnotationsFile(inputPath, outputPath, pageSelection(pageNum), ann, nil, false)
}
