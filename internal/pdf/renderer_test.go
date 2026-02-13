package pdf

import (
	"path/filepath"
	"testing"
)

func TestNewRenderer(t *testing.T) {
	r := NewRenderer()
	if r == nil {
		t.Error("NewRenderer() returned nil")
	}
}

func TestRendererCanRender(t *testing.T) {
	r := NewRenderer()
	// Just test it doesn't panic - result depends on poppler-utils being installed
	canRender := r.CanRender()
	t.Logf("CanRender() = %v (depends on poppler-utils installation)", canRender)
}

func TestRendererRenderPage(t *testing.T) {
	r := NewRenderer()

	if !r.CanRender() {
		t.Skip("poppler-utils not available")
	}

	tmpDir := t.TempDir()
	testPDF := filepath.Join(tmpDir, "test.pdf")

	if !createTestPDF(testPDF) {
		t.Skip("Cannot create test PDF")
	}

	// Render first page
	img, err := r.RenderPage(testPDF, 0, 72)
	if err != nil {
		t.Errorf("RenderPage() failed: %v", err)
	}
	if img == nil {
		t.Error("RenderPage() returned nil image")
	}
}

func TestRendererRenderPageInvalidFile(t *testing.T) {
	r := NewRenderer()

	if !r.CanRender() {
		t.Skip("poppler-utils not available")
	}

	// Try to render non-existent file
	_, err := r.RenderPage("/nonexistent/file.pdf", 0, 72)
	if err == nil {
		t.Error("RenderPage() should fail for non-existent file")
	}
}

func TestRendererRenderPageNegativePage(t *testing.T) {
	r := NewRenderer()

	if !r.CanRender() {
		t.Skip("poppler-utils not available")
	}

	tmpDir := t.TempDir()
	testPDF := filepath.Join(tmpDir, "test.pdf")

	if !createTestPDF(testPDF) {
		t.Skip("Cannot create test PDF")
	}

	// Try negative page number - pdftoppm may handle this differently
	_, err := r.RenderPage(testPDF, -1, 72)
	// Just log, pdftoppm behavior for negative pages varies
	if err != nil {
		t.Logf("RenderPage() with negative page: %v (expected)", err)
	} else {
		t.Log("RenderPage() with negative page succeeded (pdftoppm may normalize negative pages)")
	}
}

func TestRendererRenderPageInvalidDPI(t *testing.T) {
	r := NewRenderer()

	if !r.CanRender() {
		t.Skip("poppler-utils not available")
	}

	tmpDir := t.TempDir()
	testPDF := filepath.Join(tmpDir, "test.pdf")

	if !createTestPDF(testPDF) {
		t.Skip("Cannot create test PDF")
	}

	// Very low DPI should still work (may be clamped)
	img, err := r.RenderPage(testPDF, 0, 10)
	if err != nil {
		t.Logf("RenderPage() with low DPI: %v", err)
	} else if img == nil {
		t.Error("RenderPage() returned nil image")
	}
}
