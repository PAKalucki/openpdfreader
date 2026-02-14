package pdf

import (
	"errors"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func TestResolveTextExtractCommand(t *testing.T) {
	lookPath := func(name string) (string, error) {
		if name == "pdftotext" {
			return "/usr/bin/pdftotext", nil
		}
		return "", errors.New("not found")
	}

	command, args, err := resolveTextExtractCommand(lookPath, "/tmp/test.pdf", 1, "userpw", "ownerpw")
	if err != nil {
		t.Fatalf("resolveTextExtractCommand() returned error: %v", err)
	}
	if command != "pdftotext" {
		t.Fatalf("command = %q, want pdftotext", command)
	}

	expected := []string{
		"-f", "2",
		"-l", "2",
		"-layout",
		"-upw", "userpw",
		"-opw", "ownerpw",
		"/tmp/test.pdf",
		"-",
	}
	if len(args) != len(expected) {
		t.Fatalf("args length = %d, want %d", len(args), len(expected))
	}
	for i := range expected {
		if args[i] != expected[i] {
			t.Fatalf("args[%d] = %q, want %q", i, args[i], expected[i])
		}
	}
}

func TestResolveTextExtractCommandMissingBackend(t *testing.T) {
	lookPath := func(name string) (string, error) {
		return "", errors.New("not found")
	}

	_, _, err := resolveTextExtractCommand(lookPath, "/tmp/test.pdf", 0, "", "")
	if err == nil {
		t.Fatal("resolveTextExtractCommand() expected error when pdftotext is unavailable")
	}
}

func TestDocumentExtractTextValidation(t *testing.T) {
	doc := &Document{}
	if _, err := doc.ExtractText(0); err == nil {
		t.Fatal("ExtractText() expected error when no document loaded")
	}

	doc.ctx = &model.Context{}
	doc.pageCount = 2
	doc.path = "/tmp/test.pdf"

	if _, err := doc.ExtractText(-1); err == nil {
		t.Fatal("ExtractText() expected error for negative page number")
	}
	if _, err := doc.ExtractText(3); err == nil {
		t.Fatal("ExtractText() expected error for out-of-range page number")
	}
}
