package glyphs

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

var previewLineNames = []string{
	Success, Error, Warning, Info, Arrow, Star, Heart, Folder, File, Lock, Search, Gear,
}

var previewGroups = []struct {
	title string
	names []string
}{
	{
		title: "Status",
		names: []string{Success, Error, Warning, Info, Check, Cross},
	},
	{
		title: "Navigation",
		names: []string{Arrow, ArrowLeft, ArrowUp, ArrowDown, Bullet, Ellipsis, Separator},
	},
	{
		title: "Objects",
		names: []string{Star, Heart, Folder, File, Lock, Unlock, Search, Gear, User, Link},
	},
	{
		title: "Actions",
		names: []string{Edit, Trash, Refresh, Download, Upload, Clock, Play, Pause, Stop, Plus, Minus},
	},
}

// PreviewLine renders a short single-line glyph preview.
// If width is > 0, the output is bounded to the given display width.
func PreviewLine(set Set, width int) string {
	parts := make([]string, 0, len(previewLineNames))
	currentWidth := 0

	for _, name := range previewLineNames {
		glyph := set.Get(name)
		if glyph == "" {
			continue
		}

		partWidth := utf8.RuneCountInString(glyph)
		if len(parts) > 0 {
			partWidth++
		}
		if width > 0 && len(parts) > 0 && currentWidth+partWidth > width {
			break
		}

		parts = append(parts, glyph)
		currentWidth += partWidth
	}

	if len(parts) == 0 {
		return "(none)"
	}
	return strings.Join(parts, " ")
}

// Preview renders a longer sample showing a glyph set.
func Preview(w io.Writer, set Set) {
	fmt.Fprintf(w, "Glyph Preview (%s)\n\n", set.Mode())
	fmt.Fprintf(w, "  %s\n\n", PreviewLine(set, 0))

	for _, group := range previewGroups {
		lines := previewGroupLines(set, group.names)
		if len(lines) == 0 {
			continue
		}
		fmt.Fprintf(w, "  %s:\n", group.title)
		for _, line := range lines {
			fmt.Fprintf(w, "    %s\n", line)
		}
		fmt.Fprintln(w)
	}
}

func previewGroupLines(set Set, names []string) []string {
	lines := make([]string, 0, len(names))
	for _, name := range names {
		glyph := set.Get(name)
		if glyph == "" {
			continue
		}
		lines = append(lines, fmt.Sprintf("%-10s %s", name, glyph))
	}
	return lines
}
