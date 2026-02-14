package pdf

import (
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// ImageExporter exports PDF pages to image files.
type ImageExporter struct{}

// NewImageExporter creates a page image exporter.
func NewImageExporter() *ImageExporter {
	return &ImageExporter{}
}

// ExportToImages exports all pages from inputPath into outputDir.
// Supported formats: png, jpg/jpeg.
func (e *ImageExporter) ExportToImages(inputPath, outputDir, format string, scale float64) ([]string, error) {
	if inputPath == "" {
		return nil, errors.New("input path is required")
	}
	if outputDir == "" {
		return nil, errors.New("output directory is required")
	}
	if scale <= 0 {
		return nil, errors.New("scale must be greater than zero")
	}

	format = normalizeImageFormat(format)
	if format != "png" && format != "jpg" {
		return nil, errors.New("unsupported image format: use png or jpg")
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, err
	}

	doc, err := Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer doc.Close()

	outFiles := make([]string, 0, doc.PageCount())
	for i := 0; i < doc.PageCount(); i++ {
		img, err := doc.RenderPage(i, scale)
		if err != nil {
			return nil, fmt.Errorf("render page %d: %w", i+1, err)
		}

		outPath := filepath.Join(outputDir, fmt.Sprintf("page-%03d.%s", i+1, format))
		f, err := os.Create(outPath)
		if err != nil {
			return nil, err
		}

		if format == "png" {
			err = png.Encode(f, img)
		} else {
			err = jpeg.Encode(f, img, &jpeg.Options{Quality: 92})
		}
		closeErr := f.Close()
		if err != nil {
			return nil, err
		}
		if closeErr != nil {
			return nil, closeErr
		}

		outFiles = append(outFiles, outPath)
	}

	return outFiles, nil
}

func normalizeImageFormat(format string) string {
	format = strings.ToLower(strings.TrimSpace(format))
	if format == "jpeg" {
		return "jpg"
	}
	return format
}
