package pdf

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

var addImageWatermarksForReaderFile = api.AddImageWatermarksForReaderFile

// SignatureManager applies drawn signatures to PDF pages.
type SignatureManager struct{}

// NewSignatureManager creates a signature manager.
func NewSignatureManager() *SignatureManager {
	return &SignatureManager{}
}

// AddSignatureToPage adds a signature image as an on-top stamp on one page.
func (m *SignatureManager) AddSignatureToPage(inputPath, outputPath string, pageNum int, signaturePNG []byte) error {
	if inputPath == "" {
		return errors.New("input path is required")
	}
	if pageNum < 0 {
		return errors.New("page number out of range")
	}
	if len(signaturePNG) == 0 {
		return errors.New("signature image is empty")
	}

	selectedPages := []string{strconv.Itoa(pageNum + 1)}
	// Stamp near bottom-right. Scale is relative to page size.
	desc := "pos:br, off:18 18, scalefactor:0.18 rel, rot:0"

	return addImageWatermarksForReaderFile(
		inputPath,
		outputPath,
		selectedPages,
		true,
		bytes.NewReader(signaturePNG),
		desc,
		nil,
	)
}
