package pdf

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdfcolor "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

var addAnnotationsFile = api.AddAnnotationsFile

// Annotator provides basic PDF annotation operations.
type Annotator struct{}

// NewAnnotator creates a new annotator.
func NewAnnotator() *Annotator {
	return &Annotator{}
}

// AddHighlight adds a highlight annotation to the selected page.
func (a *Annotator) AddHighlight(inputPath, outputPath string, pageNum int, contents string) error {
	if err := validateAnnotationInput(inputPath, pageNum); err != nil {
		return err
	}

	rect := types.NewRectangle(100, 600, 500, 640)
	quad := types.NewQuadLiteralForRect(rect)
	ann := model.NewHighlightAnnotation(
		*rect,
		contents,
		nextAnnotationID("hl"),
		"",
		0,
		&pdfcolor.Yellow,
		0,
		0,
		1,
		"OpenPDF Reader",
		nil,
		nil,
		"",
		"",
		types.QuadPoints{*quad},
	)

	return addAnnotationsFile(inputPath, outputPath, pageSelection(pageNum), ann, nil, false)
}

// AddText adds a text annotation to the selected page.
func (a *Annotator) AddText(inputPath, outputPath string, pageNum int, contents string) error {
	if err := validateAnnotationInput(inputPath, pageNum); err != nil {
		return err
	}

	ann := model.NewTextAnnotation(
		*types.NewRectangle(80, 620, 320, 760),
		contents,
		nextAnnotationID("txt"),
		"",
		0,
		&pdfcolor.Blue,
		"OpenPDF Reader",
		nil,
		nil,
		"",
		"",
		0,
		0,
		1,
		true,
		"Comment",
	)

	return addAnnotationsFile(inputPath, outputPath, pageSelection(pageNum), ann, nil, false)
}

// AddShape adds a square annotation to the selected page.
func (a *Annotator) AddShape(inputPath, outputPath string, pageNum int, contents string) error {
	if err := validateAnnotationInput(inputPath, pageNum); err != nil {
		return err
	}

	ann := model.NewSquareAnnotation(
		*types.NewRectangle(180, 440, 430, 650),
		contents,
		nextAnnotationID("shape"),
		"",
		0,
		&pdfcolor.Red,
		"OpenPDF Reader",
		nil,
		nil,
		"",
		"",
		&pdfcolor.LightGray,
		0,
		0,
		0,
		0,
		1,
		model.BSSolid,
		false,
		0,
	)

	return addAnnotationsFile(inputPath, outputPath, pageSelection(pageNum), ann, nil, false)
}

func validateAnnotationInput(inputPath string, pageNum int) error {
	if inputPath == "" {
		return errors.New("input path is required")
	}
	if pageNum < 0 {
		return errors.New("page number out of range")
	}
	return nil
}

func pageSelection(pageNum int) []string {
	return []string{strconv.Itoa(pageNum + 1)}
}

func nextAnnotationID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
