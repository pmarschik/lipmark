package lipmark

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/pmarschik/lipmark/glyphs"
	"github.com/pmarschik/lipmark/theme"
)

// Render processes a markup string, replacing placeholders with styled text.
//
// Placeholder syntax:
//   - {name} — styled glyph: looks up glyph by name, renders with matching style
//   - {name text} — styled text: applies the named style to the text
//   - {{ — literal { (escaped)
//   - }} — literal } (escaped)
//
// If a name has no matching glyph or style, the placeholder is left as-is.
func Render(s string, glyphSet glyphs.Set, themeSet theme.Set) string {
	return renderMarkup(s, glyphSet, themeSet)
}

func renderMarkup(s string, gs glyphs.Set, ts theme.Set) string {
	var b strings.Builder
	b.Grow(len(s))

	i := 0
	for i < len(s) {
		// Escaped braces.
		if i+1 < len(s) && s[i] == '{' && s[i+1] == '{' {
			b.WriteByte('{')
			i += 2
			continue
		}
		if i+1 < len(s) && s[i] == '}' && s[i+1] == '}' {
			b.WriteByte('}')
			i += 2
			continue
		}

		// Start of placeholder.
		if s[i] == '{' {
			end := strings.IndexByte(s[i:], '}')
			if end < 0 {
				// No closing brace — write rest literally.
				b.WriteString(s[i:])
				break
			}
			inner := s[i+1 : i+end]
			b.WriteString(resolvePlaceholder(inner, gs, ts))
			i += end + 1
			continue
		}

		b.WriteByte(s[i])
		i++
	}

	return b.String()
}

func resolvePlaceholder(inner string, gs glyphs.Set, ts theme.Set) string {
	inner = strings.TrimSpace(inner)
	if inner == "" {
		return "{}"
	}

	// Split into name and optional text.
	name, text, hasText := strings.Cut(inner, " ")

	if hasText {
		// {name text} — apply style to text.
		style := ts.Get(name)
		return style.Render(text)
	}

	// {name} — styled glyph.
	glyph := gs.Get(name)
	style := ts.Get(name)

	if glyph != "" {
		return style.Render(glyph)
	}

	// No glyph found — if there's a style, render the name; otherwise literal.
	if !isZeroStyle(style) {
		return style.Render(name)
	}

	return "{" + inner + "}"
}

func isZeroStyle(s lipgloss.Style) bool {
	return s.Render("x") == "x"
}
