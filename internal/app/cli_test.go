package app

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRunCLIHelp(t *testing.T) {
	var out bytes.Buffer
	if err := RunCLI([]string{"--help"}, &out); err != nil {
		t.Fatalf("RunCLI(help) returned error: %v", err)
	}
	if !strings.Contains(out.String(), "Usage:") {
		t.Fatalf("help output missing usage: %q", out.String())
	}
}

func TestRunCLIUnknownCommand(t *testing.T) {
	var out bytes.Buffer
	err := RunCLI([]string{"unknown"}, &out)
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
}

func TestRunCLIMergeDispatch(t *testing.T) {
	orig := cliMerge
	defer func() { cliMerge = orig }()

	called := false
	cliMerge = func(inputs []string, output string) error {
		called = true
		if len(inputs) != 2 || inputs[0] != "a.pdf" || inputs[1] != "b.pdf" {
			t.Fatalf("inputs = %#v", inputs)
		}
		if output != "out.pdf" {
			t.Fatalf("output = %q, want out.pdf", output)
		}
		return nil
	}

	var out bytes.Buffer
	err := RunCLI([]string{"merge", "--inputs", "a.pdf,b.pdf", "--output", "out.pdf"}, &out)
	if err != nil {
		t.Fatalf("RunCLI(merge) returned error: %v", err)
	}
	if !called {
		t.Fatal("expected merge function to be called")
	}
}

func TestRunCLIExportImagesValidation(t *testing.T) {
	var out bytes.Buffer
	err := RunCLI([]string{"export-images", "--input", "a.pdf", "--output-dir", "/tmp", "--scale", "0"}, &out)
	if err == nil {
		t.Fatal("expected validation error for scale=0")
	}
}

func TestRunCLIExportTextDispatch(t *testing.T) {
	orig := cliExportText
	defer func() { cliExportText = orig }()

	cliExportText = func(input, output string) error {
		if input != "a.pdf" || output != "out.txt" {
			t.Fatalf("got input=%q output=%q", input, output)
		}
		return nil
	}

	var out bytes.Buffer
	err := RunCLI([]string{"export-text", "--input", "a.pdf", "--output", "out.txt"}, &out)
	if err != nil {
		t.Fatalf("RunCLI(export-text) returned error: %v", err)
	}
}

func TestRunCLIPropagatesOperationError(t *testing.T) {
	orig := cliSplit
	defer func() { cliSplit = orig }()

	cliSplit = func(input, outputDir string) error {
		return errors.New("split failed")
	}

	var out bytes.Buffer
	err := RunCLI([]string{"split", "--input", "a.pdf", "--output-dir", "out"}, &out)
	if err == nil || !strings.Contains(err.Error(), "split failed") {
		t.Fatalf("expected propagated split error, got: %v", err)
	}
}
