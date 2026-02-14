package ui

import "testing"

func TestParseFieldAssignments(t *testing.T) {
	values, err := parseFieldAssignments("name=Alice\nactive=true\n# comment\ncity=NYC")
	if err != nil {
		t.Fatalf("parseFieldAssignments() returned error: %v", err)
	}

	if len(values) != 3 {
		t.Fatalf("assignment count = %d, want 3", len(values))
	}
	if values["name"] != "Alice" {
		t.Fatalf("name = %q, want Alice", values["name"])
	}
	if values["active"] != "true" {
		t.Fatalf("active = %q, want true", values["active"])
	}
}

func TestParseFieldAssignmentsErrors(t *testing.T) {
	if _, err := parseFieldAssignments(""); err == nil {
		t.Fatal("parseFieldAssignments(\"\") expected error")
	}
	if _, err := parseFieldAssignments("novalue"); err == nil {
		t.Fatal("parseFieldAssignments(\"novalue\") expected error")
	}
	if _, err := parseFieldAssignments("=x"); err == nil {
		t.Fatal("parseFieldAssignments(\"=x\") expected error")
	}
}

func TestParsePageSelection(t *testing.T) {
	got, err := parsePageSelection("1, 3-5, 3,10", 10)
	if err != nil {
		t.Fatalf("parsePageSelection() returned error: %v", err)
	}

	want := []int{1, 3, 4, 5, 10}
	if len(got) != len(want) {
		t.Fatalf("page count = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("page[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestParsePageSelectionAll(t *testing.T) {
	got, err := parsePageSelection("all", 3)
	if err != nil {
		t.Fatalf("parsePageSelection(all) returned error: %v", err)
	}

	want := []int{1, 2, 3}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("page[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestParsePageSelectionErrors(t *testing.T) {
	tests := []string{
		"",
		"1,,2",
		"9",
		"3-1",
		"2-a",
		"a",
	}

	for _, in := range tests {
		if _, err := parsePageSelection(in, 5); err == nil {
			t.Fatalf("parsePageSelection(%q) expected error", in)
		}
	}
}

func TestPageNumbersToString(t *testing.T) {
	if got := pageNumbersToString(nil); got != "-" {
		t.Fatalf("pageNumbersToString(nil) = %q, want -", got)
	}
	if got := pageNumbersToString([]int{1, 2, 5}); got != "1,2,5" {
		t.Fatalf("pageNumbersToString([1,2,5]) = %q, want 1,2,5", got)
	}
}

func TestTabTitleForPath(t *testing.T) {
	if got := tabTitleForPath(""); got != "Untitled" {
		t.Fatalf("tabTitleForPath(\"\") = %q, want Untitled", got)
	}
	if got := tabTitleForPath("/tmp/docs/test.pdf"); got != "test.pdf" {
		t.Fatalf("tabTitleForPath(path) = %q, want test.pdf", got)
	}
}

func TestNormalizeThemeName(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"light", "light"},
		{"dark", "dark"},
		{"system", "system"},
		{"", "system"},
		{"invalid", "system"},
	}

	for _, tt := range tests {
		got := normalizeThemeName(tt.in)
		if got != tt.want {
			t.Fatalf("normalizeThemeName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
