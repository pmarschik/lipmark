package theme

import (
	"bytes"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func testPreviewSet() Set {
	return NewWithOptions(WithDefaults()).Resolve("nocolor")
}

func testColoredPreviewSet() Set {
	return NewWithOptions(WithDefaults()).Resolve("default")
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

func TestPreviewLine_TextMode(t *testing.T) {
	set := testPreviewSet()
	got := PreviewLine(set, 0, WithPreviewLineMode(PreviewLineText))
	if got != "red orange yellow green blue indigo violet" {
		t.Fatalf("unexpected preview line: %q", got)
	}
}

func TestPreviewLine_TextMode_BoundsWidth(t *testing.T) {
	set := testPreviewSet()
	got := PreviewLine(set, 18, WithPreviewLineMode(PreviewLineText))
	if got != "red orange yellow" {
		t.Fatalf("expected bounded preview line, got %q", got)
	}
}

func TestPreviewLine_DefaultsToText(t *testing.T) {
	set := testPreviewSet()
	got := PreviewLine(set, 0)
	if got != "red orange yellow green blue indigo violet" {
		t.Fatalf("unexpected preview line: %q", got)
	}
}

func TestPreviewLine_SwatchMode(t *testing.T) {
	set := testColoredPreviewSet()
	got := PreviewLine(set, 0, WithPreviewLineMode(PreviewLineSwatches))
	// default preset: Success, Error, Warning, Info, Command, Flag, Heading, Key have foreground;
	// Muted, Value, Path are faint/empty with no foreground → skipped
	if lipgloss.Width(got) != 8 {
		t.Fatalf("expected 8 swatch chars, got width %d: %q", lipgloss.Width(got), got)
	}
	if !strings.Contains(got, "█") {
		t.Fatalf("expected block chars in swatch output, got %q", got)
	}
}

func TestPreviewLine_SwatchMode_BoundsWidth(t *testing.T) {
	set := testColoredPreviewSet()
	got := PreviewLine(set, 5, WithPreviewLineMode(PreviewLineSwatches))
	if lipgloss.Width(got) != 5 {
		t.Fatalf("expected 5 swatch chars, got width %d: %q", lipgloss.Width(got), got)
	}
}

func TestPreviewLine_SwatchMode_SkipsNoColor(t *testing.T) {
	set := testPreviewSet() // nocolor preset has no foreground colors
	got := PreviewLine(set, 0, WithPreviewLineMode(PreviewLineSwatches))
	if got != "" {
		t.Fatalf("expected empty swatch line for nocolor set, got %q", got)
	}
}

func TestPreviewLine_WithSwatchStyles(t *testing.T) {
	reg := NewWithOptions(WithDefaults())
	reg.Register("custom", map[string]StyleDef{
		"my-brand": {Foreground: "#FF00FF"},
	})
	set := reg.Resolve("custom")

	got := PreviewLine(set, 0,
		WithPreviewLineMode(PreviewLineSwatches),
		WithSwatchStyles("my-brand"),
	)
	// "my-brand" has a foreground; default styles have no foreground in "custom" preset
	if lipgloss.Width(got) != 1 {
		t.Fatalf("expected 1 swatch for custom style, got width %d: %q", lipgloss.Width(got), got)
	}
	if !strings.Contains(got, "█") {
		t.Fatalf("expected block char, got %q", got)
	}
}

func TestPreview(t *testing.T) {
	var buf bytes.Buffer
	Preview(&buf, testPreviewSet())

	got := buf.String()
	if !strings.Contains(got, "Theme Preview") {
		t.Fatalf("expected theme heading, got %q", got)
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
