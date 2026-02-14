package app

import "testing"

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
		{" LIGHT ", "light"},
	}

	for _, tt := range tests {
		got := normalizeThemeName(tt.in)
		if got != tt.want {
			t.Fatalf("normalizeThemeName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestAppIconResource(t *testing.T) {
	res := appIconResource()
	if res == nil {
		t.Fatal("appIconResource() returned nil")
	}
	if len(res.Content()) == 0 {
		t.Fatal("app icon resource has empty content")
	}
}
