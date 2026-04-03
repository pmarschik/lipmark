package lipmark

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pmarschik/lipmark/glyphs"
	"github.com/pmarschik/lipmark/theme"
)

func testUI() (ui *UI, stdout, stderr *bytes.Buffer) {
	gs := glyphs.NewWithOptions(glyphs.WithDefaults()).Resolve("ascii")
	ts := theme.NewWithOptions(theme.WithDefaults()).Resolve("nocolor")
	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	ui = NewUI(
		WithUIGlyphs(gs),
		WithUITheme(ts),
		WithUIOut(stdout),
		WithUIErr(stderr),
	)
	return ui, stdout, stderr
}

func TestUI_Sprintf(t *testing.T) {
	ui, _, _ := testUI()
	got := ui.Sprintf("{success} deployed %s", "my-repo")
	if got != "[ok] deployed my-repo" {
		t.Errorf("expected '[ok] deployed my-repo', got %q", got)
	}
}

func TestUI_Sprintf_MultipleGlyphs(t *testing.T) {
	ui, _, _ := testUI()
	got := ui.Sprintf("{success} pushed %d commits {arrow} %s", 3, "origin")
	if got != "[ok] pushed 3 commits -> origin" {
		t.Errorf("got %q", got)
	}
}

func TestUI_Sprintf_StyledText(t *testing.T) {
	ui, _, _ := testUI()
	got := ui.Sprintf("repo {bold %s} updated", "my-repo")
	if got != "repo my-repo updated" { // nocolor
		t.Errorf("got %q", got)
	}
}

func TestUI_Stdout(t *testing.T) {
	ui, stdout, stderr := testUI()
	ui.Stdout("{success} done")
	if !strings.Contains(stdout.String(), "[ok] done") {
		t.Errorf("stdout: %q", stdout.String())
	}
	if stderr.Len() > 0 {
		t.Error("should not write to stderr")
	}
}

func TestUI_Stderr(t *testing.T) {
	ui, stdout, stderr := testUI()
	ui.Stderr("{warning} careful")
	if !strings.Contains(stderr.String(), "[warn] careful") {
		t.Errorf("stderr: %q", stderr.String())
	}
	if stdout.Len() > 0 {
		t.Error("should not write to stdout")
	}
}

func TestUI_Success(t *testing.T) {
	ui, _, stderr := testUI()
	ui.Success("deployed %s", "app")
	if !strings.Contains(stderr.String(), "[ok]") || !strings.Contains(stderr.String(), "deployed app") {
		t.Errorf("stderr: %q", stderr.String())
	}
}

func TestUI_Error(t *testing.T) {
	ui, _, stderr := testUI()
	ui.Error("failed: %s", "timeout")
	if !strings.Contains(stderr.String(), "[err]") || !strings.Contains(stderr.String(), "failed: timeout") {
		t.Errorf("stderr: %q", stderr.String())
	}
}

func TestUI_Warning(t *testing.T) {
	ui, _, stderr := testUI()
	ui.Warning("disk at %d%%", 90)
	if !strings.Contains(stderr.String(), "[warn]") {
		t.Errorf("stderr: %q", stderr.String())
	}
}

func TestUI_Info(t *testing.T) {
	ui, _, stderr := testUI()
	ui.Info("building")
	if !strings.Contains(stderr.String(), "[info]") {
		t.Errorf("stderr: %q", stderr.String())
	}
}

func TestUI_Note(t *testing.T) {
	ui, _, stderr := testUI()
	ui.Note("see docs")
	if !strings.Contains(stderr.String(), "->") || !strings.Contains(stderr.String(), "see docs") {
		t.Errorf("stderr: %q", stderr.String())
	}
}

func TestUI_Glyph(t *testing.T) {
	ui, _, _ := testUI()
	if ui.Glyph("success") != "[ok]" {
		t.Errorf("got %q", ui.Glyph("success"))
	}
}

func TestUI_StyledGlyph(t *testing.T) {
	ui, _, _ := testUI()
	got := ui.StyledGlyph("success")
	if got != "[ok]" { // nocolor
		t.Errorf("got %q", got)
	}
}

func TestUI_StyledGlyph_Empty(t *testing.T) {
	ui, _, _ := testUI()
	got := ui.StyledGlyph("nonexistent")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestUI_Print_NoNewline(t *testing.T) {
	ui, stdout, _ := testUI()
	ui.Print("hello")
	if strings.HasSuffix(stdout.String(), "\n") {
		t.Error("Print should not add newline")
	}
}

func TestNewUI_Defaults(t *testing.T) {
	ui := NewUI()
	if ui == nil {
		t.Fatal("NewUI returned nil")
	}
	if ui.Out == nil || ui.Err == nil {
		t.Error("writers should default to stdout/stderr")
	}
	// Should be able to render without panic.
	got := ui.Sprintf("{success}")
	if got == "" {
		t.Error("expected non-empty render")
	}
}
