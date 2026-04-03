package lipmark

import (
	"fmt"
	"io"

	"github.com/pmarschik/lipmark/glyphs"
	"github.com/pmarschik/lipmark/theme"
)

// Preview renders the combined theme and glyph previews.
// Prefer the package-level helpers in [theme] and [glyphs] when you want only one section.
func Preview(w io.Writer, gs glyphs.Set, ts theme.Set) {
	theme.Preview(w, ts)
	fmt.Fprintln(w)
	glyphs.Preview(w, gs)
}

// PreviewCompare renders two themes side by side for comparison.
func PreviewCompare(w io.Writer, _ glyphs.Set, name1 string, ts1 theme.Set, name2 string, ts2 theme.Set) {
	theme.PreviewCompare(w, name1, ts1, name2, ts2)
}
