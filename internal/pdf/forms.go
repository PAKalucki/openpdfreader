package pdf

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdfform "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/form"
)

var (
	formFieldsAPI = api.FormFields
	exportFormAPI = api.ExportForm
	fillFormAPI   = api.FillForm
)

// FormField contains form field metadata for UI display.
type FormField struct {
	Pages  []int
	Locked bool
	Type   string
	ID     string
	Name   string
	Value  string
}

// FormManager provides PDF form detection and fill operations.
type FormManager struct{}

// NewFormManager creates a form manager.
func NewFormManager() *FormManager {
	return &FormManager{}
}

// ListFields returns all form fields in a PDF.
func (m *FormManager) ListFields(inputPath string) ([]FormField, error) {
	if inputPath == "" {
		return nil, errors.New("input path is required")
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fields, err := formFieldsAPI(f, nil)
	if err != nil {
		return nil, err
	}

	result := make([]FormField, 0, len(fields))
	for _, field := range fields {
		result = append(result, FormField{
			Pages:  append([]int(nil), field.Pages...),
			Locked: field.Locked,
			Type:   field.Typ.String(),
			ID:     field.ID,
			Name:   field.Name,
			Value:  field.V,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if len(result[i].Pages) > 0 && len(result[j].Pages) > 0 && result[i].Pages[0] != result[j].Pages[0] {
			return result[i].Pages[0] < result[j].Pages[0]
		}
		if result[i].Name != result[j].Name {
			return result[i].Name < result[j].Name
		}
		return result[i].ID < result[j].ID
	})

	return result, nil
}

// FillFields updates form fields by ID or field name.
func (m *FormManager) FillFields(inputPath, outputPath string, values map[string]string) error {
	if inputPath == "" {
		return errors.New("input path is required")
	}
	if len(values) == 0 {
		return errors.New("no form values provided")
	}

	src, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	formGroup, err := exportFormAPI(src, inputPath, nil)
	src.Close()
	if err != nil {
		return err
	}

	if len(formGroup.Forms) == 0 {
		return errors.New("no form fields found")
	}

	updated, err := applyFormValues(&formGroup.Forms[0], values)
	if err != nil {
		return err
	}
	if updated == 0 {
		return errors.New("no matching form fields found")
	}

	payload, err := json.Marshal(formGroup)
	if err != nil {
		return err
	}

	return writeFilledForm(inputPath, outputPath, bytes.NewReader(payload))
}

func writeFilledForm(inputPath, outputPath string, formJSON io.Reader) (err error) {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	tmpFile := inputPath + ".tmp"
	if outputPath != "" && outputPath != inputPath {
		tmpFile = outputPath
	}

	outFile, err := os.Create(tmpFile)
	if err != nil {
		inFile.Close()
		return err
	}

	defer func() {
		if err != nil {
			outFile.Close()
			inFile.Close()
			os.Remove(tmpFile)
			return
		}
		if err = outFile.Close(); err != nil {
			return
		}
		if err = inFile.Close(); err != nil {
			return
		}
		if outputPath == "" || outputPath == inputPath {
			err = os.Rename(tmpFile, inputPath)
		}
	}()

	return fillFormAPI(inFile, formJSON, outFile, nil)
}

func applyFormValues(formData *pdfform.Form, values map[string]string) (int, error) {
	updated := 0

	for key, value := range values {
		matched := false

		for _, f := range formData.TextFields {
			if matchesFormField(f.ID, f.Name, key) && !f.Locked {
				f.Value = value
				updated++
				matched = true
			}
		}
		for _, f := range formData.DateFields {
			if matchesFormField(f.ID, f.Name, key) && !f.Locked {
				f.Value = value
				updated++
				matched = true
			}
		}
		for _, f := range formData.ComboBoxes {
			if matchesFormField(f.ID, f.Name, key) && !f.Locked {
				f.Value = value
				updated++
				matched = true
			}
		}
		for _, f := range formData.RadioButtonGroups {
			if matchesFormField(f.ID, f.Name, key) && !f.Locked {
				f.Value = value
				updated++
				matched = true
			}
		}
		for _, f := range formData.ListBoxes {
			if matchesFormField(f.ID, f.Name, key) && !f.Locked {
				f.Values = parseListValues(value)
				updated++
				matched = true
			}
		}
		for _, f := range formData.CheckBoxes {
			if matchesFormField(f.ID, f.Name, key) && !f.Locked {
				b, parseErr := parseBoolValue(value)
				if parseErr != nil {
					return 0, fmt.Errorf("field %q: %w", key, parseErr)
				}
				f.Value = b
				updated++
				matched = true
			}
		}

		_ = matched
	}

	return updated, nil
}

func matchesFormField(id, name, key string) bool {
	return key == id || key == name
}

func parseListValues(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		v := strings.TrimSpace(part)
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

func parseBoolValue(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true, nil
	case "0", "false", "no", "off":
		return false, nil
	}

	if i, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
		return i != 0, nil
	}

	return false, errors.New("expected boolean value (true/false/yes/no/1/0)")
}
