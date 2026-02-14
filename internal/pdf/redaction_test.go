package pdf

import (
	"errors"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func TestRedactorValidation(t *testing.T) {
	r := NewRedactor()

	if err := r.ApplyVisualRedaction("", "", 0, "r"); err == nil {
		t.Fatal("expected error for empty input path")
	}
	if err := r.ApplyVisualRedaction("in.pdf", "", -1, "r"); err == nil {
		t.Fatal("expected error for negative page number")
	}
}

func TestRedactorCallsAPI(t *testing.T) {
	original := addAnnotationsFile
	defer func() {
		addAnnotationsFile = original
	}()

	addAnnotationsFile = func(inFile, outFile string, selectedPages []string, ar model.AnnotationRenderer, conf *model.Configuration, incr bool) error {
		if inFile != "in.pdf" {
			t.Fatalf("inFile = %q, want in.pdf", inFile)
		}
		if outFile != "out.pdf" {
			t.Fatalf("outFile = %q, want out.pdf", outFile)
		}
		if len(selectedPages) != 1 || selectedPages[0] != "4" {
			t.Fatalf("selectedPages = %#v, want []string{\"4\"}", selectedPages)
		}
		if ar == nil {
			t.Fatal("annotation renderer should not be nil")
		}
		if conf != nil {
			t.Fatal("expected nil conf")
		}
		if incr {
			t.Fatal("expected incr=false")
		}
		return nil
	}

	r := NewRedactor()
	if err := r.ApplyVisualRedaction("in.pdf", "out.pdf", 3, "confidential"); err != nil {
		t.Fatalf("ApplyVisualRedaction() returned error: %v", err)
	}
}

func TestRedactorPropagatesError(t *testing.T) {
	original := addAnnotationsFile
	defer func() {
		addAnnotationsFile = original
	}()

	addAnnotationsFile = func(inFile, outFile string, selectedPages []string, ar model.AnnotationRenderer, conf *model.Configuration, incr bool) error {
		return errors.New("annotate failed")
	}

	r := NewRedactor()
	err := r.ApplyVisualRedaction("in.pdf", "", 0, "x")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "annotate failed") {
		t.Fatalf("error = %q, want propagated error", err.Error())
	}
}
