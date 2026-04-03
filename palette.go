package lipmark

import (
	"strings"

	"github.com/pmarschik/lipmark/theme"
)

// PaletteSwatches returns a compact color swatch string showing the key palette
// colors of a theme. Each swatch is a colored block character. Useful for
// rendering next to theme names in listings.
func PaletteSwatches(ts theme.Set) string {
	// Show the semantic colors as colored blocks.
	styles := []string{
		theme.Error, theme.Warning, theme.Success, theme.Info,
		theme.Command, theme.Heading, theme.Muted,
	}

	var parts []string
	for _, name := range styles {
		s := ts.Get(name)
		if isZeroStyle(s) {
			continue
		}
		parts = append(parts, s.Render("█"))
	}

	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, "")
}

// PaletteSwatchesFull returns a wider swatch showing all base16 palette colors
// (base08-base0F) as colored blocks.
func PaletteSwatchesFull(ts theme.Set) string {
	keys := []string{
		"base08", "base09", "base0A", "base0B",
		"base0C", "base0D", "base0E", "base0F",
	}

	var parts []string
	for _, key := range keys {
		s := ts.Get(key)
		if isZeroStyle(s) {
			continue
		}
		parts = append(parts, s.Render("██"))
	}

	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, "")
}

// isZeroStyle is defined in markup.go
