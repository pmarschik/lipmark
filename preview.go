package lipmark

import (
	"fmt"
	"io"
	"strings"

	"github.com/pmarschik/lipmark/glyphs"
	"github.com/pmarschik/lipmark/theme"
)

// Preview renders a sample output demonstrating a theme's styles.
// Useful for "help themes <name>" preview commands.
func Preview(w io.Writer, gs glyphs.Set, ts theme.Set) {
	ui := &UI{Glyphs: gs, Theme: ts, Out: w, Err: w}

	fmt.Fprintln(w, ui.Sprintf("{heading Theme Preview}"))
	fmt.Fprintln(w)

	// Semantic styles.
	fmt.Fprintln(w, ui.Sprintf("  {success} {success Success message}"))
	fmt.Fprintln(w, ui.Sprintf("  {error} {error Error message}"))
	fmt.Fprintln(w, ui.Sprintf("  {warning} {warning Warning message}"))
	fmt.Fprintln(w, ui.Sprintf("  {info} {info Info message}"))
	fmt.Fprintln(w, ui.Sprintf("  {arrow} {muted Muted note}"))
	fmt.Fprintln(w)

	// Text styles.
	fmt.Fprintln(w, ui.Sprintf("  {bold Bold text}  {dim Dim text}  {italic Italic text}  {underline Underlined}"))
	fmt.Fprintln(w)

	// CLI chrome.
	fmt.Fprintln(w, ui.Sprintf("  {heading Commands:}"))
	fmt.Fprintln(w, ui.Sprintf("    {command serve}    Start the server"))
	fmt.Fprintln(w, ui.Sprintf("    {command build}    Build the project"))
	fmt.Fprintln(w)
	fmt.Fprintln(w, ui.Sprintf("  {heading Flags:}"))
	fmt.Fprintln(w, ui.Sprintf("    {flag --port}={dim 8080}    Listen port"))
	fmt.Fprintln(w, ui.Sprintf("    {flag --config}={path ~/.config/app.yaml}"))
	fmt.Fprintln(w)

	// Glyphs.
	printGlyphRow(w, gs, ts)
}

func printGlyphRow(w io.Writer, gs glyphs.Set, ts theme.Set) {
	names := []string{
		glyphs.Success, glyphs.Error, glyphs.Warning, glyphs.Info,
		glyphs.Arrow, glyphs.Star, glyphs.Heart, glyphs.Folder,
		glyphs.File, glyphs.Lock, glyphs.Search, glyphs.Gear,
	}
	var parts []string
	for _, n := range names {
		g := gs.Get(n)
		if g == "" {
			continue
		}
		parts = append(parts, ts.Get(n).Render(g))
	}
	if len(parts) > 0 {
		fmt.Fprintf(w, "  Glyphs: %s\n", strings.Join(parts, " "))
	}
}

// PreviewCompare renders two themes side by side for comparison.
func PreviewCompare(w io.Writer, _ glyphs.Set, name1 string, ts1 theme.Set, name2 string, ts2 theme.Set) {
	fmt.Fprintf(w, "%-30s %s\n", name1, name2)
	fmt.Fprintf(w, "%-30s %s\n",
		strings.Repeat("─", len(name1)),
		strings.Repeat("─", len(name2)),
	)

	pairs := []struct {
		label string
		style string
	}{
		{"Success", theme.Success},
		{"Error", theme.Error},
		{"Warning", theme.Warning},
		{"Info", theme.Info},
		{"Command", theme.Command},
		{"Heading", theme.Heading},
		{"Muted", theme.Muted},
	}

	for _, p := range pairs {
		left := ts1.Get(p.style).Render(p.label)
		right := ts2.Get(p.style).Render(p.label)
		fmt.Fprintf(w, "  %-28s   %s\n", left, right)
	}
}
