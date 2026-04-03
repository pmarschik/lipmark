package theme

import (
	"bytes"
	"strings"
	"testing"
)

func testPreviewSet() Set {
	return NewWithOptions(WithDefaults()).Resolve("nocolor")
}

func TestPreviewLine(t *testing.T) {
	got := PreviewLine(testPreviewSet(), 0)
	if got != "red orange yellow green blue indigo violet" {
		t.Fatalf("unexpected preview line: %q", got)
	}
}

func TestPreviewLine_BoundsWidth(t *testing.T) {
	got := PreviewLine(testPreviewSet(), 18)
	if got != "red orange yellow" {
		t.Fatalf("expected bounded preview line, got %q", got)
	}
}

func TestPreview(t *testing.T) {
	var buf bytes.Buffer
	Preview(&buf, testPreviewSet())

	got := buf.String()
	if !strings.Contains(got, "Theme Preview") {
		t.Fatalf("expected theme heading, got %q", got)
	}
	if !strings.Contains(got, "red orange yellow green blue indigo violet") {
		t.Fatalf("expected single-line rainbow preview, got %q", got)
	}
	if !strings.Contains(got, "Success message") {
		t.Fatalf("expected detailed theme preview, got %q", got)
	}
}

func TestPreview_WithPreviewSection(t *testing.T) {
	var buf bytes.Buffer
	Preview(
		&buf,
		testPreviewSet(),
		WithPreviewSection(
			"Custom:",
			PreviewItem{Style: Success, Text: "ship it"},
			PreviewItem{Style: Heading, Text: "PROJECT-123"},
		),
	)

	got := buf.String()
	if !strings.Contains(got, "Custom:") {
		t.Fatalf("expected custom section heading, got %q", got)
	}
	if !strings.Contains(got, "ship it") || !strings.Contains(got, "PROJECT-123") {
		t.Fatalf("expected custom preview items, got %q", got)
	}
}

func TestPreviewCompare(t *testing.T) {
	var buf bytes.Buffer
	set := testPreviewSet()
	PreviewCompare(&buf, "one", set, "two", set)

	got := buf.String()
	if !strings.Contains(got, "one") || !strings.Contains(got, "two") {
		t.Fatalf("expected compare headings, got %q", got)
	}
	if !strings.Contains(got, "Success") {
		t.Fatalf("expected compare rows, got %q", got)
	}
}
