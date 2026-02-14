package pdf

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func TestTextExporterValidation(t *testing.T) {
	e := NewTextExporter()

	if err := e.ExportToText("", "out.txt"); err == nil {
		t.Fatal("expected error for empty input path")
	}
	if err := e.ExportToText("in.pdf", ""); err == nil {
		t.Fatal("expected error for empty output path")
	}
}

func TestTextExporterExport(t *testing.T) {
	origOpen := openTextExportDocument
	origExtract := extractTextForPage
	defer func() {
		openTextExportDocument = origOpen
		extractTextForPage = origExtract
	}()

	openTextExportDocument = func(path string) (*Document, error) {
		return &Document{
			path:      path,
			ctx:       &model.Context{},
			pageCount: 2,
		}, nil
	}

	extractTextForPage = func(doc *Document, pageNum int) (string, error) {
		if pageNum == 0 {
			return "first page text", nil
		}
		return "second page text", nil
	}

	outPath := filepath.Join(t.TempDir(), "out.txt")
	e := NewTextExporter()
	if err := e.ExportToText("in.pdf", outPath); err != nil {
		t.Fatalf("ExportToText() returned error: %v", err)
	}

	b, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("ReadFile() returned error: %v", err)
	}
	out := string(b)

	if !strings.Contains(out, "--- Page 1 ---") || !strings.Contains(out, "first page text") {
		t.Fatalf("missing page 1 content: %q", out)
	}
	if !strings.Contains(out, "--- Page 2 ---") || !strings.Contains(out, "second page text") {
		t.Fatalf("missing page 2 content: %q", out)
	}
}

func TestTextExporterPropagatesExtractionError(t *testing.T) {
	origOpen := openTextExportDocument
	origExtract := extractTextForPage
	defer func() {
		openTextExportDocument = origOpen
		extractTextForPage = origExtract
	}()

	openTextExportDocument = func(path string) (*Document, error) {
		return &Document{
			path:      path,
			ctx:       &model.Context{},
			pageCount: 1,
		}, nil
	}

	extractTextForPage = func(doc *Document, pageNum int) (string, error) {
		return "", errors.New("extract failed")
	}

	e := NewTextExporter()
	err := e.ExportToText("in.pdf", filepath.Join(t.TempDir(), "out.txt"))
	if err == nil {
		t.Fatal("expected extraction error")
	}
	if !strings.Contains(err.Error(), "extract failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}
