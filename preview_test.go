package lipmark

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pmarschik/lipmark/glyphs"
	"github.com/pmarschik/lipmark/theme"
)

func TestPreview(t *testing.T) {
	gs := glyphs.NewWithOptions(glyphs.WithDefaults()).Resolve("ascii")
	ts := theme.NewWithOptions(theme.WithDefaults()).Resolve("nocolor")

	var buf bytes.Buffer
	Preview(&buf, gs, ts)

	got := buf.String()
	if !strings.Contains(got, "Theme Preview") {
		t.Fatalf("expected theme preview output, got %q", got)
	}
	if !strings.Contains(got, "Glyph Preview (ascii)") {
		t.Fatalf("expected glyph preview output, got %q", got)
	}
}
