package pdf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	openTextExportDocument = Open
	extractTextForPage     = func(doc *Document, pageNum int) (string, error) {
		return doc.ExtractText(pageNum)
	}
)

// TextExporter exports PDF pages into a plain text document.
type TextExporter struct{}

// NewTextExporter creates a text exporter.
func NewTextExporter() *TextExporter {
	return &TextExporter{}
}

// ExportToText exports all pages from inputPath into outputPath.
func (e *TextExporter) ExportToText(inputPath, outputPath string) error {
	if inputPath == "" {
		return errors.New("input path is required")
	}
	if outputPath == "" {
		return errors.New("output path is required")
	}

	doc, err := openTextExportDocument(inputPath)
	if err != nil {
		return err
	}
	defer doc.Close()

	var b strings.Builder
	for i := 0; i < doc.PageCount(); i++ {
		text, err := extractTextForPage(doc, i)
		if err != nil {
			return fmt.Errorf("extract page %d: %w", i+1, err)
		}
		if i > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(fmt.Sprintf("--- Page %d ---\n", i+1))
		b.WriteString(text)
		b.WriteString("\n")
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(outputPath, []byte(b.String()), 0644)
}
