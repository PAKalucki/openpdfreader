package pdf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImageExporterValidation(t *testing.T) {
	e := NewImageExporter()

	if _, err := e.ExportToImages("", "/tmp", "png", 1.0); err == nil {
		t.Fatal("expected error for empty input path")
	}
	if _, err := e.ExportToImages("in.pdf", "", "png", 1.0); err == nil {
		t.Fatal("expected error for empty output directory")
	}
	if _, err := e.ExportToImages("in.pdf", "/tmp", "png", 0); err == nil {
		t.Fatal("expected error for invalid scale")
	}
	if _, err := e.ExportToImages("in.pdf", "/tmp", "gif", 1.0); err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestImageExporterExportToImagesPNG(t *testing.T) {
	tmpDir := t.TempDir()
	inFile := filepath.Join(tmpDir, "in.pdf")
	outDir := filepath.Join(tmpDir, "images")

	if !createTestPDF(inFile) {
		t.Skip("Cannot create test PDF")
	}

	e := NewImageExporter()
	files, err := e.ExportToImages(inFile, outDir, "png", 1.0)
	if err != nil {
		t.Fatalf("ExportToImages() returned error: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 exported file, got %d", len(files))
	}
	if filepath.Ext(files[0]) != ".png" {
		t.Fatalf("expected .png output, got %s", filepath.Ext(files[0]))
	}
	if _, err := os.Stat(files[0]); err != nil {
		t.Fatalf("expected output file to exist: %v", err)
	}
}

func TestImageExporterExportToImagesJPEGAlias(t *testing.T) {
	tmpDir := t.TempDir()
	inFile := filepath.Join(tmpDir, "in.pdf")
	outDir := filepath.Join(tmpDir, "images")

	if !createTestPDF(inFile) {
		t.Skip("Cannot create test PDF")
	}

	e := NewImageExporter()
	files, err := e.ExportToImages(inFile, outDir, "jpeg", 1.0)
	if err != nil {
		t.Fatalf("ExportToImages() returned error: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 exported file, got %d", len(files))
	}
	if filepath.Ext(files[0]) != ".jpg" {
		t.Fatalf("expected .jpg output, got %s", filepath.Ext(files[0]))
	}
}
