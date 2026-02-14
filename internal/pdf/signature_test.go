package pdf

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func TestSignatureManagerValidation(t *testing.T) {
	m := NewSignatureManager()

	if err := m.AddSignatureToPage("", "", 0, []byte("png")); err == nil {
		t.Fatal("expected error for empty input path")
	}
	if err := m.AddSignatureToPage("in.pdf", "", -1, []byte("png")); err == nil {
		t.Fatal("expected error for negative page number")
	}
	if err := m.AddSignatureToPage("in.pdf", "", 0, nil); err == nil {
		t.Fatal("expected error for empty signature image")
	}
}

func TestSignatureManagerCallsAPI(t *testing.T) {
	original := addImageWatermarksForReaderFile
	defer func() {
		addImageWatermarksForReaderFile = original
	}()

	addImageWatermarksForReaderFile = func(
		inFile, outFile string,
		selectedPages []string,
		onTop bool,
		r io.Reader,
		desc string,
		conf *model.Configuration,
	) error {
		if inFile != "in.pdf" {
			t.Fatalf("inFile = %q, want in.pdf", inFile)
		}
		if outFile != "out.pdf" {
			t.Fatalf("outFile = %q, want out.pdf", outFile)
		}
		if len(selectedPages) != 1 || selectedPages[0] != "3" {
			t.Fatalf("selectedPages = %#v, want []string{\"3\"}", selectedPages)
		}
		if !onTop {
			t.Fatal("expected onTop = true")
		}
		if conf != nil {
			t.Fatal("expected nil conf")
		}
		if !strings.Contains(desc, "pos:br") {
			t.Fatalf("desc = %q, want bottom-right placement", desc)
		}
		if r == nil {
			t.Fatal("reader should not be nil")
		}
		return nil
	}

	m := NewSignatureManager()
	if err := m.AddSignatureToPage("in.pdf", "out.pdf", 2, []byte("pngdata")); err != nil {
		t.Fatalf("AddSignatureToPage() returned error: %v", err)
	}
}

func TestSignatureManagerPropagatesAPIError(t *testing.T) {
	original := addImageWatermarksForReaderFile
	defer func() {
		addImageWatermarksForReaderFile = original
	}()

	addImageWatermarksForReaderFile = func(
		inFile, outFile string,
		selectedPages []string,
		onTop bool,
		r io.Reader,
		desc string,
		conf *model.Configuration,
	) error {
		return errors.New("stamp failed")
	}

	m := NewSignatureManager()
	err := m.AddSignatureToPage("in.pdf", "", 0, []byte("png"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "stamp failed") {
		t.Fatalf("error = %q, want propagated error", err.Error())
	}
}
