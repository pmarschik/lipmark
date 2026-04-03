package lipmark

import (
	"testing"

	"github.com/pmarschik/lipmark/glyphs"
	"github.com/pmarschik/lipmark/theme"
)

func noColorSets() (gs glyphs.Set, ts theme.Set) {
	gs = glyphs.NewWithOptions(glyphs.WithDefaults()).Resolve("ascii")
	ts = theme.NewWithOptions(theme.WithDefaults()).Resolve("nocolor")
	return gs, ts
}

func TestRender_StyledGlyph(t *testing.T) {
	gs, ts := noColorSets()
	got := Render("{success}", gs, ts)
	if got != "[ok]" {
		t.Errorf("expected [ok], got %q", got)
	}
}

func TestRender_StyledText(t *testing.T) {
	gs, ts := noColorSets()
	got := Render("{bold hello}", gs, ts)
	if got != "hello" { // nocolor renders plain
		t.Errorf("expected hello, got %q", got)
	}
}

func TestRender_Mixed(t *testing.T) {
	gs, ts := noColorSets()
	got := Render("{success} done {arrow} next", gs, ts)
	if got != "[ok] done -> next" {
		t.Errorf("expected '[ok] done -> next', got %q", got)
	}
}

func TestRender_EscapedBraces(t *testing.T) {
	gs, ts := noColorSets()
	got := Render("use {{name}} syntax", gs, ts)
	if got != "use {name} syntax" {
		t.Errorf("expected 'use {name} syntax', got %q", got)
	}
}

func TestRender_UnknownPlaceholder(t *testing.T) {
	gs, ts := noColorSets()
	got := Render("{unknown}", gs, ts)
	if got != "{unknown}" {
		t.Errorf("unknown placeholder should be left as-is, got %q", got)
	}
}

func TestRender_EmptyPlaceholder(t *testing.T) {
	gs, ts := noColorSets()
	got := Render("{}", gs, ts)
	if got != "{}" {
		t.Errorf("empty placeholder should be left as-is, got %q", got)
	}
}

func TestRender_NoPlaceholders(t *testing.T) {
	gs, ts := noColorSets()
	got := Render("plain text", gs, ts)
	if got != "plain text" {
		t.Errorf("expected plain text, got %q", got)
	}
}

func TestRender_MultipleGlyphs(t *testing.T) {
	gs, ts := noColorSets()
	got := Render("{success} {arrow} {error}", gs, ts)
	if got != "[ok] -> [err]" {
		t.Errorf("expected '[ok] -> [err]', got %q", got)
	}
}

func TestRender_UnclosedBrace(t *testing.T) {
	gs, ts := noColorSets()
	got := Render("text {unclosed", gs, ts)
	if got != "text {unclosed" {
		t.Errorf("unclosed brace should be literal, got %q", got)
	}
}

func TestRender_StyledTextWithSpaces(t *testing.T) {
	gs, ts := noColorSets()
	got := Render("{bold some longer text}", gs, ts)
	if got != "some longer text" { // nocolor
		t.Errorf("expected 'some longer text', got %q", got)
	}
}
