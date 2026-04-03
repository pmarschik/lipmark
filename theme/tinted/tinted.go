// Package tinted converts base16 and base24 color schemes from the tinted-theming
// project (https://github.com/tinted-theming/schemes) into lipmark theme StyleDefs.
//
// Base16 defines 16 palette colors (base00-base0F) with semantic meaning:
//
//	base00-base07: background to foreground gradient
//	base08: red        base0C: cyan
//	base09: orange     base0D: blue
//	base0A: yellow     base0E: purple/magenta
//	base0B: green      base0F: accent
//
// Base24 extends base16 with 8 additional colors (base10-base17):
//
//	base10-base11: darker background shades
//	base12-base17: brighter/lighter variants of base08-base0D
//
// Both formats are supported transparently. The [DefaultMapping] only uses
// base00-base0F keys which exist in both formats. Apps can create custom
// mappings that use base10-base17 for base24 schemes.
//
// # Loading from file
//
// Download scheme files from https://github.com/tinted-theming/schemes
// and load them with the [load] sub-package:
//
//	import "github.com/pmarschik/lipmark/theme/tinted/load"
//	defs, err := load.File("dracula.yaml")
//	themeReg.Register("dracula", defs)
//
// # Embedded themes
//
// Use [embedded.WithEmbedded] from the embedded sub-package to register the
// bundled schemes (solarized-dark, solarized-light, dracula, catppuccin-mocha, nord):
//
//	import "github.com/pmarschik/lipmark/theme/tinted/embedded"
//	reg := theme.NewWithOptions(theme.WithDefaults(), embedded.WithEmbedded())
//
// # Custom mapping
//
// Override the default palette-to-style mapping:
//
//	defs, err := load.FileWithMapping("dracula.yaml", tinted.Mapping{
//	    Success: "base0B", Error: "base08", Warning: "base0A",
//	    // ...
//	})
package tinted

import (
	"charm.land/lipgloss/v2"
	"github.com/pmarschik/lipmark/theme"
)

// Palette key constants for use with [theme.Set.Get] on tinted themes.
// These are the raw palette color names from tinted-theming schemes,
// available alongside the semantic style names (bold, error, success, etc.).
//
// Base16 keys (available in all schemes):
const (
	Base00 = "base00" // Darkest background
	Base01 = "base01" // Lighter background (status bars, line numbers)
	Base02 = "base02" // Selection background
	Base03 = "base03" // Comments, invisibles
	Base04 = "base04" // Dark foreground (status bars)
	Base05 = "base05" // Default foreground
	Base06 = "base06" // Light foreground
	Base07 = "base07" // Lightest foreground
	Base08 = "base08" // Red / errors
	Base09 = "base09" // Orange / integers
	Base0A = "base0A" // Yellow / warnings / classes
	Base0B = "base0B" // Green / success / strings
	Base0C = "base0C" // Cyan / support / regex
	Base0D = "base0D" // Blue / info / functions
	Base0E = "base0E" // Purple / keywords
	Base0F = "base0F" // Accent / deprecated
)

// Base24 keys (only available in base24 schemes):
const (
	Base10 = "base10" // Darker background shade
	Base11 = "base11" // Darkest background shade
	Base12 = "base12" // Bright red (base08 lighter)
	Base13 = "base13" // Bright yellow (base0A lighter)
	Base14 = "base14" // Bright green (base0B lighter)
	Base15 = "base15" // Bright cyan (base0C lighter)
	Base16 = "base16" // Bright blue (base0D lighter)
	Base17 = "base17" // Bright purple (base0E lighter)
)

// Palette provides typed access to all resolved palette colors.
// Embed this in your app's theme struct alongside [theme.BaseStyles]:
//
//	type MyTheme struct {
//	    theme.BaseStyles
//	    tinted.Palette
//	    IssueKey lipgloss.Style `style:"issue_key"`
//	}
//	var t MyTheme
//	reg.ResolveInto("dracula", &t)
//	t.Base0C  // cyan from palette
//	t.Success // semantic style from BaseStyles
type Palette struct {
	// Base16 colors (available in all schemes).
	Base00 lipgloss.Style `style:"base00"` // Darkest background
	Base01 lipgloss.Style `style:"base01"` // Lighter background
	Base02 lipgloss.Style `style:"base02"` // Selection background
	Base03 lipgloss.Style `style:"base03"` // Comments
	Base04 lipgloss.Style `style:"base04"` // Dark foreground
	Base05 lipgloss.Style `style:"base05"` // Default foreground
	Base06 lipgloss.Style `style:"base06"` // Light foreground
	Base07 lipgloss.Style `style:"base07"` // Lightest foreground
	Base08 lipgloss.Style `style:"base08"` // Red
	Base09 lipgloss.Style `style:"base09"` // Orange
	Base0A lipgloss.Style `style:"base0A"` // Yellow
	Base0B lipgloss.Style `style:"base0B"` // Green
	Base0C lipgloss.Style `style:"base0C"` // Cyan
	Base0D lipgloss.Style `style:"base0D"` // Blue
	Base0E lipgloss.Style `style:"base0E"` // Purple
	Base0F lipgloss.Style `style:"base0F"` // Accent

	// Base24 colors (only available in base24 schemes; zero-value in base16).
	Base10 lipgloss.Style `style:"base10"` // Darker background
	Base11 lipgloss.Style `style:"base11"` // Darkest background
	Base12 lipgloss.Style `style:"base12"` // Bright red
	Base13 lipgloss.Style `style:"base13"` // Bright yellow
	Base14 lipgloss.Style `style:"base14"` // Bright green
	Base15 lipgloss.Style `style:"base15"` // Bright cyan
	Base16 lipgloss.Style `style:"base16"` // Bright blue
	Base17 lipgloss.Style `style:"base17"` // Bright purple
}

// PaletteDefs provides typed access to palette color definitions for registration.
type PaletteDefs struct {
	Base00 theme.StyleDef `style:"base00"`
	Base01 theme.StyleDef `style:"base01"`
	Base02 theme.StyleDef `style:"base02"`
	Base03 theme.StyleDef `style:"base03"`
	Base04 theme.StyleDef `style:"base04"`
	Base05 theme.StyleDef `style:"base05"`
	Base06 theme.StyleDef `style:"base06"`
	Base07 theme.StyleDef `style:"base07"`
	Base08 theme.StyleDef `style:"base08"`
	Base09 theme.StyleDef `style:"base09"`
	Base0A theme.StyleDef `style:"base0A"`
	Base0B theme.StyleDef `style:"base0B"`
	Base0C theme.StyleDef `style:"base0C"`
	Base0D theme.StyleDef `style:"base0D"`
	Base0E theme.StyleDef `style:"base0E"`
	Base0F theme.StyleDef `style:"base0F"`
	Base10 theme.StyleDef `style:"base10"`
	Base11 theme.StyleDef `style:"base11"`
	Base12 theme.StyleDef `style:"base12"`
	Base13 theme.StyleDef `style:"base13"`
	Base14 theme.StyleDef `style:"base14"`
	Base15 theme.StyleDef `style:"base15"`
	Base16 theme.StyleDef `style:"base16"`
	Base17 theme.StyleDef `style:"base17"`
}

// Scheme represents a parsed tinted-theming color scheme (base16 or base24).
type Scheme struct {
	Palette map[string]string
	System  string
	Name    string
	Author  string
	Variant string
}

// Mapping defines which palette entry maps to each style name.
// Use empty string to skip a style. Keys are palette names like "base05", "base0D", etc.
// base24-specific keys (base10-base17) are only available in base24 schemes.
type Mapping struct {
	Bold      string
	Dim       string
	Italic    string
	Underline string
	Success   string
	Error     string
	Warning   string
	Info      string
	Muted     string
	Command   string
	Flag      string
	Heading   string
	Key       string
	Value     string
	Path      string
}

// DefaultMapping returns the standard palette-to-style mapping.
// It uses only base00-base0F keys, which are available in both base16 and base24 schemes.
func DefaultMapping() Mapping {
	return Mapping{
		Bold:    "base05",
		Dim:     "base03",
		Italic:  "base05",
		Success: "base0B",
		Error:   "base08",
		Warning: "base0A",
		Info:    "base0D",
		Muted:   "base04",
		Command: "base0C",
		Flag:    "base04",
		Heading: "base0E",
		Key:     "base0C",
		Path:    "base04",
	}
}
