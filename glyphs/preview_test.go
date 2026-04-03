package glyphs

import (
	"bytes"
	"strings"
	"testing"
)

func testPreviewSet() Set {
	return NewWithOptions(WithDefaults()).Resolve("ascii")
}

func TestPreviewLine(t *testing.T) {
	got := PreviewLine(testPreviewSet(), 0)
	if !strings.HasPrefix(got, "[ok] [err] [warn] [info]") {
		t.Fatalf("unexpected preview line: %q", got)
	}
}

func TestPreviewLine_BoundsWidth(t *testing.T) {
	got := PreviewLine(testPreviewSet(), 16)
	if got != "[ok] [err]" {
		t.Fatalf("expected bounded preview line, got %q", got)
	}
}

func TestPreview(t *testing.T) {
	var buf bytes.Buffer
	Preview(&buf, testPreviewSet())

	got := buf.String()
	if !strings.Contains(got, "Glyph Preview (ascii)") {
		t.Fatalf("expected glyph heading, got %q", got)
	}
	if !strings.Contains(got, "Status:") || !strings.Contains(got, "Actions:") {
		t.Fatalf("expected grouped glyph preview, got %q", got)
	}
	if !strings.Contains(got, "success    [ok]") {
		t.Fatalf("expected named glyph rows, got %q", got)
	}
}
