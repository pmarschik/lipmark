package lipmark

import (
	"fmt"
	"io"
	"os"

	"github.com/pmarschik/lipmark/glyphs"
	"github.com/pmarschik/lipmark/theme"

	"charm.land/lipgloss/v2"
)

// UI provides themed, markup-aware output for CLI applications.
// It wraps resolved glyph and theme sets and provides formatted output
// with {placeholder} syntax to stdout and stderr.
//
// Create a UI with [NewUI], then use [UI.Stdout], [UI.Stderr], or the
// semantic helpers ([UI.Success], [UI.Error], etc.) for output.
//
// Markup syntax:
//   - {name} — styled glyph (looks up glyph + style by name)
//   - {name text} — styled text (applies named style to text)
//   - {{ / }} — literal braces
//   - Standard %s, %d, etc. format verbs work as in [fmt.Sprintf]
type UI struct {
	Out    io.Writer
	Err    io.Writer
	Glyphs glyphs.Set
	Theme  theme.Set
}

// UIOption configures a [UI].
type UIOption func(*UI)

// WithUIGlyphs sets the glyph set.
func WithUIGlyphs(g glyphs.Set) UIOption {
	return func(u *UI) { u.Glyphs = g }
}

// WithUITheme sets the theme set.
func WithUITheme(t theme.Set) UIOption {
	return func(u *UI) { u.Theme = t }
}

// WithUIOut sets the stdout writer.
func WithUIOut(w io.Writer) UIOption {
	return func(u *UI) { u.Out = w }
}

// WithUIErr sets the stderr writer.
func WithUIErr(w io.Writer) UIOption {
	return func(u *UI) { u.Err = w }
}

// NewUI creates a UI with auto-detected glyphs and theme.
// Override with options.
func NewUI(opts ...UIOption) *UI {
	u := &UI{
		Glyphs: glyphs.NewWithOptions(glyphs.WithDefaults()).Resolve("auto"),
		Theme:  theme.NewWithOptions(theme.WithDefaults()).Resolve("auto"),
		Out:    os.Stdout,
		Err:    os.Stderr,
	}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

// --- Markup rendering ---

// Sprintf formats a markup string and returns the result.
// Format verbs (%s, %d, etc.) are resolved first, then {placeholders}.
func (u *UI) Sprintf(format string, args ...any) string {
	s := fmt.Sprintf(format, args...)
	return Render(s, u.Glyphs, u.Theme)
}

// --- Output to stdout ---

// Stdout writes a markup-formatted line to stdout.
func (u *UI) Stdout(format string, args ...any) {
	_, _ = fmt.Fprintln(u.Out, u.Sprintf(format, args...))
}

// Print writes a markup-formatted string to stdout (no trailing newline).
func (u *UI) Print(format string, args ...any) {
	_, _ = fmt.Fprint(u.Out, u.Sprintf(format, args...))
}

// --- Output to stderr ---

// Stderr writes a markup-formatted line to stderr.
func (u *UI) Stderr(format string, args ...any) {
	_, _ = fmt.Fprintln(u.Err, u.Sprintf(format, args...))
}

// --- Semantic helpers (write to stderr) ---

// Success writes a success message: "{success} message".
func (u *UI) Success(format string, args ...any) {
	u.Stderr("{success} "+format, args...)
}

// Error writes an error message: "{error} message" with styled text.
func (u *UI) Error(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	styled := Render("{error}", u.Glyphs, u.Theme) + " " +
		u.Theme.Get(theme.Error).Render(msg)
	_, _ = fmt.Fprintln(u.Err, styled)
}

// Warning writes a warning message: "{warning} message".
func (u *UI) Warning(format string, args ...any) {
	u.Stderr("{warning} "+format, args...)
}

// Info writes an info message: "{info} message".
func (u *UI) Info(format string, args ...any) {
	u.Stderr("{info} "+format, args...)
}

// Note writes a muted note: "{arrow} message" in muted style.
func (u *UI) Note(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	styled := Render("{arrow}", u.Glyphs, u.Theme) + " " +
		u.Theme.Get(theme.Muted).Render(msg)
	_, _ = fmt.Fprintln(u.Err, styled)
}

// Status writes a plain status line to stderr.
func (u *UI) Status(format string, args ...any) {
	_, _ = fmt.Fprintln(u.Err, fmt.Sprintf(format, args...))
}

// --- Direct access ---

// Glyph returns the raw glyph character for the given name.
func (u *UI) Glyph(name string) string {
	return u.Glyphs.Get(name)
}

// StyledGlyph returns the styled glyph for the given name.
func (u *UI) StyledGlyph(name string) string {
	g := u.Glyphs.Get(name)
	if g == "" {
		return ""
	}
	return u.Theme.Get(name).Render(g)
}

// Style returns the lipgloss.Style for the given name.
func (u *UI) Style(name string) lipgloss.Style {
	return u.Theme.Get(name)
}
