package pdf

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestResolvePrintCommandUnixPrefersLP(t *testing.T) {
	lookPath := func(name string) (string, error) {
		switch name {
		case "lp":
			return "/usr/bin/lp", nil
		case "lpr":
			return "/usr/bin/lpr", nil
		default:
			return "", errors.New("not found")
		}
	}

	command, args, err := resolvePrintCommand(lookPath, "linux", "/tmp/test.pdf")
	if err != nil {
		t.Fatalf("resolvePrintCommand() returned error: %v", err)
	}
	if command != "lp" {
		t.Fatalf("command = %q, want %q", command, "lp")
	}
	if len(args) != 1 || args[0] != "/tmp/test.pdf" {
		t.Fatalf("args = %#v, want file path only", args)
	}
}

func TestResolvePrintCommandUnixFallsBackToLPR(t *testing.T) {
	lookPath := func(name string) (string, error) {
		if name == "lpr" {
			return "/usr/bin/lpr", nil
		}
		return "", errors.New("not found")
	}

	command, _, err := resolvePrintCommand(lookPath, "linux", "/tmp/test.pdf")
	if err != nil {
		t.Fatalf("resolvePrintCommand() returned error: %v", err)
	}
	if command != "lpr" {
		t.Fatalf("command = %q, want %q", command, "lpr")
	}
}

func TestResolvePrintCommandUnixNoCommand(t *testing.T) {
	lookPath := func(name string) (string, error) {
		return "", errors.New("not found")
	}

	_, _, err := resolvePrintCommand(lookPath, "linux", "/tmp/test.pdf")
	if err == nil {
		t.Fatal("resolvePrintCommand() expected error when no print command exists")
	}
}

func TestResolvePrintCommandWindows(t *testing.T) {
	lookPath := func(name string) (string, error) {
		if name == "rundll32" {
			return "C:\\Windows\\System32\\rundll32.exe", nil
		}
		return "", errors.New("not found")
	}

	command, args, err := resolvePrintCommand(lookPath, "windows", "C:\\tmp\\test.pdf")
	if err != nil {
		t.Fatalf("resolvePrintCommand() returned error: %v", err)
	}
	if command != "rundll32" {
		t.Fatalf("command = %q, want %q", command, "rundll32")
	}
	if len(args) != 2 {
		t.Fatalf("args length = %d, want 2", len(args))
	}
}

func TestPrintFileErrors(t *testing.T) {
	if err := PrintFile(""); err == nil {
		t.Fatal("PrintFile(\"\") expected error")
	}

	missing := filepath.Join(t.TempDir(), "missing.pdf")
	if err := PrintFile(missing); err == nil {
		t.Fatal("PrintFile(missing) expected error")
	}
}
