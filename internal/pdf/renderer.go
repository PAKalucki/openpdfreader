package pdf

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os/exec"
	"strconv"
)

// Renderer handles PDF page rendering.
type Renderer struct {
	usePoppler bool
}

// NewRenderer creates a new PDF renderer.
// It checks for available rendering backends.
func NewRenderer() *Renderer {
	r := &Renderer{}

	// Check if pdftoppm (poppler-utils) is available
	if _, err := exec.LookPath("pdftoppm"); err == nil {
		r.usePoppler = true
	}

	return r
}

// CanRender returns true if a rendering backend is available.
func (r *Renderer) CanRender() bool {
	return r.usePoppler
}

// RenderPage renders a PDF page to an image.
// pageNum is 0-indexed, scale is DPI (72 = 100%, 144 = 200%).
func (r *Renderer) RenderPage(pdfPath string, pageNum int, dpi int) (image.Image, error) {
	if r.usePoppler {
		return r.renderWithPoppler(pdfPath, pageNum, dpi)
	}
	return nil, fmt.Errorf("no rendering backend available")
}

func (r *Renderer) renderWithPoppler(pdfPath string, pageNum int, dpi int) (image.Image, error) {
	// pdftoppm uses 1-indexed pages
	pageStr := strconv.Itoa(pageNum + 1)
	dpiStr := strconv.Itoa(dpi)

	// Run pdftoppm to convert page to PNG
	cmd := exec.Command("pdftoppm",
		"-png",
		"-f", pageStr, // first page
		"-l", pageStr, // last page (same = single page)
		"-r", dpiStr, // resolution
		"-singlefile", // don't add page suffix
		pdfPath,
		"-", // output to stdout
	)

	// pdftoppm with "-" outputs to stdout but still needs a prefix
	// Use different approach: output to stdout with -png
	cmd = exec.Command("pdftoppm",
		"-png",
		"-f", pageStr,
		"-l", pageStr,
		"-r", dpiStr,
		"-singlefile",
		pdfPath,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("pdftoppm failed: %v: %s", err, stderr.String())
	}

	// Decode PNG from stdout
	img, err := png.Decode(&stdout)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %v", err)
	}

	return img, nil
}

// GetPageCount returns the number of pages using pdfinfo.
func (r *Renderer) GetPageCount(pdfPath string) (int, error) {
	cmd := exec.Command("pdfinfo", pdfPath)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// Parse "Pages: N" from output
	lines := bytes.Split(output, []byte("\n"))
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte("Pages:")) {
			parts := bytes.Fields(line)
			if len(parts) >= 2 {
				return strconv.Atoi(string(parts[1]))
			}
		}
	}

	return 0, fmt.Errorf("could not determine page count")
}
