package pdf

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// PrintFile sends a PDF file to the default system printer.
func PrintFile(path string) error {
	if path == "" {
		return errors.New("no file path set")
	}
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("cannot access file for printing: %w", err)
	}

	command, args, err := resolvePrintCommand(exec.LookPath, runtime.GOOS, path)
	if err != nil {
		return err
	}

	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(output))
		if msg != "" {
			return fmt.Errorf("print command failed: %w: %s", err, msg)
		}
		return fmt.Errorf("print command failed: %w", err)
	}

	return nil
}

func resolvePrintCommand(lookPath func(string) (string, error), goos, path string) (string, []string, error) {
	switch goos {
	case "windows":
		if _, err := lookPath("rundll32"); err == nil {
			return "rundll32", []string{"shell32.dll,ShellExec_RunDLL", path}, nil
		}
		return "", nil, errors.New("no print command found: rundll32 is not available")
	default:
		if _, err := lookPath("lp"); err == nil {
			return "lp", []string{path}, nil
		}
		if _, err := lookPath("lpr"); err == nil {
			return "lpr", []string{path}, nil
		}
		return "", nil, errors.New("no print command found: install lp or lpr")
	}
}
