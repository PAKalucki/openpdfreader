package pdf

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/form"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func TestParseBoolValue(t *testing.T) {
	trues := []string{"true", "yes", "on", "1", "2"}
	for _, v := range trues {
		got, err := parseBoolValue(v)
		if err != nil {
			t.Fatalf("parseBoolValue(%q) returned error: %v", v, err)
		}
		if !got {
			t.Fatalf("parseBoolValue(%q) = false, want true", v)
		}
	}

	falses := []string{"false", "no", "off", "0"}
	for _, v := range falses {
		got, err := parseBoolValue(v)
		if err != nil {
			t.Fatalf("parseBoolValue(%q) returned error: %v", v, err)
		}
		if got {
			t.Fatalf("parseBoolValue(%q) = true, want false", v)
		}
	}
}

func TestParseBoolValueInvalid(t *testing.T) {
	if _, err := parseBoolValue("maybe"); err == nil {
		t.Fatal("parseBoolValue(\"maybe\") expected error")
	}
}

func TestParseListValues(t *testing.T) {
	got := parseListValues("A, B, , C")
	if len(got) != 3 || got[0] != "A" || got[1] != "B" || got[2] != "C" {
		t.Fatalf("parseListValues() = %#v, want [A B C]", got)
	}
}

func TestApplyFormValues(t *testing.T) {
	data := &form.Form{
		TextFields: []*form.TextField{
			{ID: "id_text", Name: "name_text"},
		},
		CheckBoxes: []*form.CheckBox{
			{ID: "id_check", Name: "name_check"},
		},
		ListBoxes: []*form.ListBox{
			{ID: "id_list", Name: "name_list"},
		},
	}

	updated, err := applyFormValues(data, map[string]string{
		"name_text":  "hello",
		"id_check":   "true",
		"name_list":  "x,y",
		"unmatched1": "ignored",
	})
	if err != nil {
		t.Fatalf("applyFormValues() returned error: %v", err)
	}
	if updated != 3 {
		t.Fatalf("updated = %d, want 3", updated)
	}
	if data.TextFields[0].Value != "hello" {
		t.Fatalf("text value = %q, want hello", data.TextFields[0].Value)
	}
	if !data.CheckBoxes[0].Value {
		t.Fatal("checkbox value = false, want true")
	}
	if len(data.ListBoxes[0].Values) != 2 || data.ListBoxes[0].Values[0] != "x" || data.ListBoxes[0].Values[1] != "y" {
		t.Fatalf("list values = %#v, want [x y]", data.ListBoxes[0].Values)
	}
}

func TestApplyFormValuesInvalidCheckbox(t *testing.T) {
	data := &form.Form{
		CheckBoxes: []*form.CheckBox{
			{ID: "id_check", Name: "name_check"},
		},
	}

	_, err := applyFormValues(data, map[string]string{
		"id_check": "not-bool",
	})
	if err == nil {
		t.Fatal("applyFormValues() expected checkbox parse error")
	}
}

func TestFormManagerListFields(t *testing.T) {
	original := formFieldsAPI
	defer func() {
		formFieldsAPI = original
	}()

	formFieldsAPI = func(rs io.ReadSeeker, conf *model.Configuration) ([]form.Field, error) {
		return []form.Field{
			{Pages: []int{2}, Typ: form.FTText, ID: "b", Name: "B", V: "v2"},
			{Pages: []int{1}, Typ: form.FTCheckBox, ID: "a", Name: "A", V: "v1"},
		}, nil
	}

	input := filepath.Join(t.TempDir(), "in.pdf")
	if err := os.WriteFile(input, []byte("dummy"), 0644); err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}

	fields, err := NewFormManager().ListFields(input)
	if err != nil {
		t.Fatalf("ListFields() returned error: %v", err)
	}
	if len(fields) != 2 {
		t.Fatalf("field count = %d, want 2", len(fields))
	}
	if fields[0].Name != "A" {
		t.Fatalf("first field name = %q, want A", fields[0].Name)
	}
}

func TestFormManagerFillFields(t *testing.T) {
	origExport := exportFormAPI
	origFill := fillFormAPI
	defer func() {
		exportFormAPI = origExport
		fillFormAPI = origFill
	}()

	exportFormAPI = func(rs io.ReadSeeker, source string, conf *model.Configuration) (*form.FormGroup, error) {
		return &form.FormGroup{
			Forms: []form.Form{
				{
					TextFields: []*form.TextField{
						{ID: "id_text", Name: "name_text"},
					},
				},
			},
		}, nil
	}

	fillFormAPI = func(rs io.ReadSeeker, rd io.Reader, w io.Writer, conf *model.Configuration) error {
		b, err := io.ReadAll(rd)
		if err != nil {
			return err
		}
		if !strings.Contains(string(b), "\"value\":\"updated\"") {
			return errors.New("missing updated value")
		}
		_, err = w.Write([]byte("ok"))
		return err
	}

	tmpDir := t.TempDir()
	inFile := filepath.Join(tmpDir, "in.pdf")
	outFile := filepath.Join(tmpDir, "out.pdf")
	if err := os.WriteFile(inFile, []byte("pdf"), 0644); err != nil {
		t.Fatalf("WriteFile(in) failed: %v", err)
	}

	if err := NewFormManager().FillFields(inFile, outFile, map[string]string{"name_text": "updated"}); err != nil {
		t.Fatalf("FillFields() returned error: %v", err)
	}

	b, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("ReadFile(out) failed: %v", err)
	}
	if !bytes.Equal(b, []byte("ok")) {
		t.Fatalf("output bytes = %q, want ok", string(b))
	}
}

func TestFormManagerFillFieldsNoMatch(t *testing.T) {
	origExport := exportFormAPI
	origFill := fillFormAPI
	defer func() {
		exportFormAPI = origExport
		fillFormAPI = origFill
	}()

	exportFormAPI = func(rs io.ReadSeeker, source string, conf *model.Configuration) (*form.FormGroup, error) {
		return &form.FormGroup{
			Forms: []form.Form{
				{
					TextFields: []*form.TextField{
						{ID: "id_text", Name: "name_text"},
					},
				},
			},
		}, nil
	}

	fillCalled := false
	fillFormAPI = func(rs io.ReadSeeker, rd io.Reader, w io.Writer, conf *model.Configuration) error {
		fillCalled = true
		return nil
	}

	inFile := filepath.Join(t.TempDir(), "in.pdf")
	if err := os.WriteFile(inFile, []byte("pdf"), 0644); err != nil {
		t.Fatalf("WriteFile(in) failed: %v", err)
	}

	err := NewFormManager().FillFields(inFile, "", map[string]string{"missing": "value"})
	if err == nil {
		t.Fatal("FillFields() expected no-match error")
	}
	if fillCalled {
		t.Fatal("fillFormAPI should not be called when nothing matches")
	}
}
