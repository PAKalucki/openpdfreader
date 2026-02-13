package pdf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsPasswordError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"password error", &testError{"password required"}, true},
		{"encrypted error", &testError{"file is encrypted"}, true},
		{"decrypt error", &testError{"failed to decrypt"}, true},
		{"other error", &testError{"file not found"}, false},
		{"case insensitive", &testError{"PASSWORD REQUIRED"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPasswordError(tt.err)
			if result != tt.expected {
				t.Errorf("IsPasswordError(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestDocumentOpen(t *testing.T) {
	// Test opening non-existent file
	_, err := Open("/nonexistent/path/to/file.pdf")
	if err == nil {
		t.Error("Open should fail for non-existent file")
	}
}

func TestDocumentMethods(t *testing.T) {
	// Create a minimal test PDF using pdfcpu
	tmpDir := t.TempDir()
	testPDF := filepath.Join(tmpDir, "test.pdf")

	// Skip if we can't create a test PDF
	if !createTestPDF(testPDF) {
		t.Skip("Cannot create test PDF - pdfcpu may not be available")
	}

	doc, err := Open(testPDF)
	if err != nil {
		t.Fatalf("Failed to open test PDF: %v", err)
	}
	defer doc.Close()

	// Test Path
	if doc.Path() != testPDF {
		t.Errorf("Path() = %v, want %v", doc.Path(), testPDF)
	}

	// Test PageCount
	if doc.PageCount() < 1 {
		t.Errorf("PageCount() = %d, want >= 1", doc.PageCount())
	}

	// Test IsModified
	if doc.IsModified() {
		t.Error("IsModified() should be false for newly opened document")
	}

	// Test SaveAs
	outputPath := filepath.Join(tmpDir, "output.pdf")
	err = doc.SaveAs(outputPath)
	if err != nil {
		t.Errorf("SaveAs() failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("SaveAs() did not create output file")
	}
}

func TestDocumentRenderPage(t *testing.T) {
	tmpDir := t.TempDir()
	testPDF := filepath.Join(tmpDir, "test.pdf")

	if !createTestPDF(testPDF) {
		t.Skip("Cannot create test PDF")
	}

	doc, err := Open(testPDF)
	if err != nil {
		t.Fatalf("Failed to open test PDF: %v", err)
	}
	defer doc.Close()

	// Test rendering valid page
	img, err := doc.RenderPage(0, 1.0)
	if err != nil {
		t.Logf("RenderPage failed (may need poppler-utils): %v", err)
	} else if img == nil {
		t.Error("RenderPage returned nil image without error")
	}

	// Test rendering invalid page
	_, err = doc.RenderPage(-1, 1.0)
	if err == nil {
		t.Error("RenderPage should fail for negative page number")
	}

	_, err = doc.RenderPage(9999, 1.0)
	if err == nil {
		t.Error("RenderPage should fail for page number out of range")
	}
}

// createTestPDF creates a minimal PDF for testing
func createTestPDF(path string) bool {
	// Minimal PDF 1.4 structure
	pdfContent := `%PDF-1.4
1 0 obj
<< /Type /Catalog /Pages 2 0 R >>
endobj
2 0 obj
<< /Type /Pages /Kids [3 0 R] /Count 1 >>
endobj
3 0 obj
<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Resources << >> >>
endobj
4 0 obj
<< /Length 44 >>
stream
BT
/F1 12 Tf
100 700 Td
(Test) Tj
ET
endstream
endobj
xref
0 5
0000000000 65535 f 
0000000009 00000 n 
0000000058 00000 n 
0000000115 00000 n 
0000000214 00000 n 
trailer
<< /Size 5 /Root 1 0 R >>
startxref
308
%%EOF
`
	err := os.WriteFile(path, []byte(pdfContent), 0644)
	return err == nil
}
