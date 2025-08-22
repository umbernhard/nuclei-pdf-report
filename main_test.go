package main

import (
	"strings"
	"testing"
)

func TestSanitize(t *testing.T) {
	input := "foo_{bar}&baz"
	want := "foo\\_\\{bar\\}\\&baz"
	got := sanitize(input)
	if got != want {
		t.Errorf("sanitize(%q) = %q, want %q", input, got, want)
	}
}

func TestProcessHttp(t *testing.T) {
	input := "GET /path?foo=bar&baz=qux\r\nCookie: sessionid=abc123"
	got := processHttp(input)
	if !strings.Contains(got, "\\seqsplit{") {
		t.Errorf("processHttp did not seqsplit long cookie or query string")
	}
}

func TestParseDuration(t *testing.T) {
	input := "01:02:03"
	want := "1h2m3s"
	got, err := ParseDuration(input)
	if err != nil {
		t.Fatalf("ParseDuration returned error: %v", err)
	}
	if got != want {
		t.Errorf("ParseDuration(%q) = %q, want %q", input, got, want)
	}
}

func TestFindIndex(t *testing.T) {
	slice := []string{"critical", "high", "medium"}
	if idx := findIndex(slice, "high"); idx != 1 {
		t.Errorf("findIndex returned %d, want 1", idx)
	}
	if idx := findIndex(slice, "low"); idx != 3 {
		t.Errorf("findIndex returned %d, want 3 (not found)", idx)
	}
}
