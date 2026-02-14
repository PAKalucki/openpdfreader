package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func extractPageText(pdfPath string, pageNum int, userPassword, ownerPassword string) (string, error) {
	command, args, err := resolveTextExtractCommand(exec.LookPath, pdfPath, pageNum, userPassword, ownerPassword)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		details := strings.TrimSpace(stderr.String())
		if details != "" {
			return "", fmt.Errorf("pdftotext failed: %v: %s", err, details)
		}
		return "", fmt.Errorf("pdftotext failed: %w", err)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func resolveTextExtractCommand(
	lookPath func(string) (string, error),
	pdfPath string,
	pageNum int,
	userPassword string,
	ownerPassword string,
) (string, []string, error) {
	if _, err := lookPath("pdftotext"); err != nil {
		return "", nil, errors.New("no text extraction backend available: install pdftotext from poppler-utils")
	}

	page := strconv.Itoa(pageNum + 1)
	args := []string{
		"-f", page,
		"-l", page,
		"-layout",
	}

	if userPassword != "" {
		args = append(args, "-upw", userPassword)
	}
	if ownerPassword != "" {
		args = append(args, "-opw", ownerPassword)
	}

	args = append(args, pdfPath, "-")

	return "pdftotext", args, nil
}
