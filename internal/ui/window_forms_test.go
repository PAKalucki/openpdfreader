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
