package pdf

import (
	"errors"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func TestPageSelection(t *testing.T) {
	got := pageSelection(4)
	if len(got) != 1 || got[0] != "5" {
		t.Fatalf("pageSelection(4) = %#v, want []string{\"5\"}", got)
	}
}

func TestValidateAnnotationInput(t *testing.T) {
	if err := validateAnnotationInput("", 0); err == nil {
		t.Fatal("validateAnnotationInput() expected error for empty input path")
	}
	if err := validateAnnotationInput("in.pdf", -1); err == nil {
		t.Fatal("validateAnnotationInput() expected error for negative page number")
	}
}

func TestAnnotatorMethodsCallAPI(t *testing.T) {
	original := addAnnotationsFile
	defer func() {
		addAnnotationsFile = original
	}()

	called := 0
	addAnnotationsFile = func(inFile, outFile string, selectedPages []string, ar model.AnnotationRenderer, conf *model.Configuration, incr bool) error {
		called++
		if inFile != "in.pdf" {
			t.Fatalf("inFile = %q, want in.pdf", inFile)
		}
		if outFile != "out.pdf" {
			t.Fatalf("outFile = %q, want out.pdf", outFile)
		}
		if len(selectedPages) != 1 || selectedPages[0] != "1" {
			t.Fatalf("selectedPages = %#v, want []string{\"1\"}", selectedPages)
		}
		if ar == nil {
			t.Fatal("annotation renderer must not be nil")
		}
		if conf != nil {
			t.Fatal("expected nil conf")
		}
		if incr {
			t.Fatal("expected incr=false")
		}
		return nil
	}

	a := NewAnnotator()
	if err := a.AddHighlight("in.pdf", "out.pdf", 0, "hl"); err != nil {
		t.Fatalf("AddHighlight() returned error: %v", err)
	}
	if err := a.AddText("in.pdf", "out.pdf", 0, "note"); err != nil {
		t.Fatalf("AddText() returned error: %v", err)
	}
	if err := a.AddShape("in.pdf", "out.pdf", 0, "shape"); err != nil {
		t.Fatalf("AddShape() returned error: %v", err)
	}

	if called != 3 {
		t.Fatalf("API calls = %d, want 3", called)
	}
}

func TestAnnotatorPropagatesAPIError(t *testing.T) {
	original := addAnnotationsFile
	defer func() {
		addAnnotationsFile = original
	}()

	addAnnotationsFile = func(inFile, outFile string, selectedPages []string, ar model.AnnotationRenderer, conf *model.Configuration, incr bool) error {
		return errors.New("api failed")
	}

	a := NewAnnotator()
	err := a.AddText("in.pdf", "out.pdf", 0, "note")
	if err == nil {
		t.Fatal("AddText() expected error")
	}
	if !strings.Contains(err.Error(), "api failed") {
		t.Fatalf("error = %q, want api failed", err.Error())
	}
}
