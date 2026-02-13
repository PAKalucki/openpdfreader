package pdf

import (
	"errors"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// Merger handles PDF merging operations.
type Merger struct{}

// NewMerger creates a new Merger.
func NewMerger() *Merger {
	return &Merger{}
}

// Merge combines multiple PDF files into one.
func (m *Merger) Merge(inputPaths []string, outputPath string) error {
	if len(inputPaths) < 2 {
		return errors.New("merge requires at least 2 input files")
	}
	return api.MergeCreateFile(inputPaths, outputPath, false, nil)
}

// Split splits a PDF into individual pages.
func (m *Merger) Split(inputPath, outputDir string) error {
	return api.SplitFile(inputPath, outputDir, 1, nil)
}

// ExtractPages extracts specific pages from a PDF.
// pages is a slice of 1-indexed page numbers.
func (m *Merger) ExtractPages(inputPath string, pages []int, outputPath string) error {
	pageSelection := make([]string, len(pages))
	for i, p := range pages {
		pageSelection[i] = intToStr(p)
	}

	return api.ExtractPagesFile(inputPath, outputPath, pageSelection, nil)
}

// DeletePages removes specific pages from a PDF.
// pages is a slice of 1-indexed page numbers.
func (m *Merger) DeletePages(inputPath string, pages []int, outputPath string) error {
	pageSelection := make([]string, len(pages))
	for i, p := range pages {
		pageSelection[i] = intToStr(p)
	}

	return api.RemovePagesFile(inputPath, outputPath, pageSelection, nil)
}

// RotatePages rotates specific pages.
// rotation should be 90, 180, or 270.
func (m *Merger) RotatePages(inputPath string, rotation int, pages []int, outputPath string) error {
	pageSelection := make([]string, len(pages))
	for i, p := range pages {
		pageSelection[i] = intToStr(p)
	}

	return api.RotateFile(inputPath, outputPath, rotation, pageSelection, nil)
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}

	var digits []byte
	negative := n < 0
	if negative {
		n = -n
	}

	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}

	if negative {
		digits = append([]byte{'-'}, digits...)
	}

	return string(digits)
}
