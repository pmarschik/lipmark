package theme

import (
	"fmt"
	"io"
	"strings"

	"charm.land/lipgloss/v2"
)

// PreviewItem is a single styled sample in a detailed preview.
type PreviewItem struct {
	Style string
	Text  string
}

// PreviewSection is a titled group of styled samples in a detailed preview.
type PreviewSection struct {
	Title string
	Items []PreviewItem
}

// PreviewOption customizes detailed preview rendering.
type PreviewOption func(*previewConfig)

type previewConfig struct {
	sections []PreviewSection
}

// PreviewLineMode controls how PreviewLine renders each style segment.
type PreviewLineMode int

const (
	PreviewLineText     PreviewLineMode = iota // styled text labels (default)
	PreviewLineSwatches                        // solid █ blocks
)

// PreviewLineOption customizes PreviewLine rendering.
type PreviewLineOption func(*previewLineConfig)

type previewLineConfig struct {
	extraStyles []string
	mode        PreviewLineMode
}

// WithPreviewLineMode sets the rendering mode for PreviewLine.
func WithPreviewLineMode(mode PreviewLineMode) PreviewLineOption {
	return func(cfg *previewLineConfig) {
		cfg.mode = mode
	}
}

// WithSwatchStyles appends extra style names to the swatch band.
// Extra styles are appended after the default color-bearing built-ins.
// Has no effect in PreviewLineText mode.
func WithSwatchStyles(names ...string) PreviewLineOption {
	return func(cfg *previewLineConfig) {
		cfg.extraStyles = append(cfg.extraStyles, names...)
	}
}

var defaultSwatchStyles = []string{
	Success, Error, Warning, Info, Muted,
	Command, Flag, Heading, Key, Value, Path,
}

var previewLineSegments = []struct {
	style string
	text  string
}{
	{Success, "red"},
	{Error, "orange"},
	{Warning, "yellow"},
	{Info, "green"},
	{Command, "blue"},
	{Flag, "indigo"},
	{Heading, "violet"},
}

var defaultPreviewSections = []PreviewSection{
	{
		Title: "",
		Items: []PreviewItem{
			{Success, "Success message"},
			{Error, "Error message"},
			{Warning, "Warning message"},
			{Info, "Info message"},
			{Muted, "Muted note"},
		},
	},
	{
		Title: "",
		Items: []PreviewItem{
			{Bold, "Bold text"},
			{Dim, "Dim text"},
			{Italic, "Italic text"},
			{Underline, "Underlined"},
		},
	},
	{
		Title: "Commands:",
		Items: []PreviewItem{
			{Command, "serve"},
			{Command, "build"},
		},
	},
	{
		Title: "Flags:",
		Items: []PreviewItem{
			{Flag, "--port"},
			{Flag, "--config"},
		},
	},
}

// WithPreviewSection appends a titled section to the detailed preview.
func WithPreviewSection(title string, items ...PreviewItem) PreviewOption {
	return func(cfg *previewConfig) {
		cfg.sections = append(cfg.sections, PreviewSection{
			Title: title,
			Items: items,
		})
	}
}

// PreviewLine renders a short single-line style preview.
// If width is > 0, the output is bounded to the given display width.
// Optional PreviewLineOption values customize rendering (e.g. swatch mode).
func PreviewLine(set Set, width int, opts ...PreviewLineOption) string {
	cfg := &previewLineConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.mode == PreviewLineSwatches {
		return renderSwatchLine(set, width, cfg.extraStyles)
	}

	parts := make([]string, 0, len(previewLineSegments))
	currentWidth := 0

	for _, segment := range previewLineSegments {
		rendered := set.Get(segment.style).Render(segment.text)
		partWidth := lipgloss.Width(rendered)
		if len(parts) > 0 {
			partWidth++
		}
		if width > 0 && len(parts) > 0 && currentWidth+partWidth > width {
			break
		}

		parts = append(parts, rendered)
		currentWidth += partWidth
	}

	return strings.Join(parts, " ")
}

// Preview renders a longer sample showing a theme's styles.
func Preview(w io.Writer, set Set, opts ...PreviewOption) {
	fmt.Fprintln(w, set.Get(Heading).Render("Theme Preview"))
	fmt.Fprintln(w)
	fmt.Fprintln(w, "  "+PreviewLine(set, 0, WithPreviewLineMode(PreviewLineSwatches)))
	fmt.Fprintln(w)

	cfg := previewConfig{
		sections: append([]PreviewSection(nil), defaultPreviewSections...),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	renderPreviewSections(w, set, cfg.sections)
}

// PreviewCompare renders two themes side by side for comparison.
func PreviewCompare(w io.Writer, name1 string, set1 Set, name2 string, set2 Set) {
	fmt.Fprintf(w, "%-30s %s\n", name1, name2)
	fmt.Fprintf(w, "%-30s %s\n",
		strings.Repeat("─", len(name1)),
		strings.Repeat("─", len(name2)),
	)

	pairs := []struct {
		label string
		style string
	}{
		{"Success", Success},
		{"Error", Error},
		{"Warning", Warning},
		{"Info", Info},
		{"Command", Command},
		{"Heading", Heading},
		{"Muted", Muted},
	}

	for _, pair := range pairs {
		left := set1.Get(pair.style).Render(pair.label)
		right := set2.Get(pair.style).Render(pair.label)
		fmt.Fprintf(w, "  %-28s   %s\n", left, right)
	}
}

func renderSwatchLine(set Set, width int, extraStyles []string) string {
	styles := make([]string, 0, len(defaultSwatchStyles)+len(extraStyles))
	styles = append(styles, defaultSwatchStyles...)
	styles = append(styles, extraStyles...)
	var sb strings.Builder
	count := 0
	for _, name := range styles {
		style := set.Get(name)
		if _, ok := style.GetForeground().(lipgloss.NoColor); ok {
			continue
		}
		if width > 0 && count >= width {
			break
		}
		sb.WriteString(style.Render("█"))
		count++
	}
	return sb.String()
}

func renderPreviewSections(w io.Writer, set Set, sections []PreviewSection) {
	for i, section := range sections {
		switch {
		case section.Title != "":
			fmt.Fprintln(w, "  "+set.Get(Heading).Render(section.Title))
			for _, item := range section.Items {
				fmt.Fprintf(w, "    %s\n", set.Get(item.Style).Render(item.Text))
			}
		case len(section.Items) > 0 && isTextStyleRow(section.Items):
			fmt.Fprintln(w, "  "+renderInlineItems(set, section.Items))
		default:
			for _, item := range section.Items {
				fmt.Fprintf(w, "  %s\n", set.Get(item.Style).Render(item.Text))
			}
		}

		if i < len(sections)-1 {
			fmt.Fprintln(w)
		}
	}
}

func isTextStyleRow(items []PreviewItem) bool {
	if len(items) == 0 {
		return false
	}
	for _, item := range items {
		switch item.Style {
		case Bold, Dim, Italic, Underline:
		default:
			return false
		}
	}
	return true
}

func renderInlineItems(set Set, items []PreviewItem) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		parts = append(parts, set.Get(item.Style).Render(item.Text))
	}
	return strings.Join(parts, "  ")
}
