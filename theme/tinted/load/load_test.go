package load

import (
	"strings"
	"testing"

	"github.com/pmarschik/lipmark/theme"
	"github.com/pmarschik/lipmark/theme/tinted"
)

const draculaBase16YAML = `
system: "base16"
name: "Dracula"
author: "Jamy Golden"
variant: "dark"
palette:
  base00: "#282a36"
  base01: "#363447"
  base02: "#44475a"
  base03: "#6272a4"
  base04: "#9ea8c7"
  base05: "#f8f8f2"
  base06: "#f0f1f4"
  base07: "#ffffff"
  base08: "#ff5555"
  base09: "#ffb86c"
  base0A: "#f1fa8c"
  base0B: "#50fa7b"
  base0C: "#8be9fd"
  base0D: "#80bfff"
  base0E: "#ff79c6"
  base0F: "#bd93f9"
`

const draculaBase24YAML = `
system: "base24"
name: "Dracula"
author: "FredHappyface"
variant: "dark"
palette:
  base00: "#282a36"
  base01: "#363447"
  base02: "#44475a"
  base03: "#6272a4"
  base04: "#9ea8c7"
  base05: "#f8f8f2"
  base06: "#f0f1f4"
  base07: "#ffffff"
  base08: "#ff5555"
  base09: "#ffb86c"
  base0A: "#f1fa8c"
  base0B: "#50fa7b"
  base0C: "#8be9fd"
  base0D: "#80bfff"
  base0E: "#ff79c6"
  base0F: "#bd93f9"
  base10: "#21222c"
  base11: "#191a21"
  base12: "#ff6e6e"
  base13: "#ffffa5"
  base14: "#69ff94"
  base15: "#a4ffff"
  base16: "#99b3ff"
  base17: "#ff92df"
`

func TestBytes(t *testing.T) {
	defs, err := Bytes([]byte(draculaBase16YAML))
	if err != nil {
		t.Fatal(err)
	}
	if defs[theme.Error].Foreground != "#ff5555" {
		t.Errorf("expected error=#ff5555, got %q", defs[theme.Error].Foreground)
	}
	if defs[theme.Success].Foreground != "#50fa7b" {
		t.Errorf("expected success=#50fa7b, got %q", defs[theme.Success].Foreground)
	}
	if defs[theme.Command].Foreground != "#8be9fd" {
		t.Errorf("expected command=#8be9fd, got %q", defs[theme.Command].Foreground)
	}
	if defs[theme.Heading].Foreground != "#ff79c6" {
		t.Errorf("expected heading=#ff79c6, got %q", defs[theme.Heading].Foreground)
	}
}

func TestReader(t *testing.T) {
	defs, err := Reader(strings.NewReader(draculaBase16YAML))
	if err != nil {
		t.Fatal(err)
	}
	if len(defs) == 0 {
		t.Error("expected non-empty defs")
	}
}

func TestBytesWithMapping_Custom(t *testing.T) {
	m := tinted.DefaultMapping()
	m.Error = "base0E" // override: use purple for errors
	defs, err := BytesWithMapping([]byte(draculaBase16YAML), m)
	if err != nil {
		t.Fatal(err)
	}
	if defs[theme.Error].Foreground != "#ff79c6" {
		t.Errorf("expected custom error=#ff79c6, got %q", defs[theme.Error].Foreground)
	}
}

func TestBytes_InvalidYAML(t *testing.T) {
	_, err := Bytes([]byte("{{invalid"))
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestParseScheme(t *testing.T) {
	s, err := ParseScheme(strings.NewReader(draculaBase16YAML))
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != "Dracula" {
		t.Errorf("expected name=Dracula, got %q", s.Name)
	}
	if s.Variant != "dark" {
		t.Errorf("expected variant=dark, got %q", s.Variant)
	}
	if len(s.Palette) != 16 {
		t.Errorf("expected 16 palette entries, got %d", len(s.Palette))
	}
}

func TestConvert_SkipsEmptyMapping(t *testing.T) {
	m := tinted.Mapping{} // all empty
	defs, err := BytesWithMapping([]byte(draculaBase16YAML), m)
	if err != nil {
		t.Fatal(err)
	}
	if len(defs) != 0 {
		t.Errorf("expected 0 defs with empty mapping, got %d", len(defs))
	}
}

func TestRegisterWithThemeRegistry(t *testing.T) {
	defs, err := Bytes([]byte(draculaBase16YAML))
	if err != nil {
		t.Fatal(err)
	}
	reg := theme.NewWithOptions(theme.WithDefaults())
	reg.Register("dracula", defs)
	s := reg.Resolve("dracula")
	if s.Get(theme.Error).Render("x") == "" {
		t.Error("expected styled error render")
	}
}

func TestBase24Scheme(t *testing.T) {
	s, err := ParseScheme(strings.NewReader(draculaBase24YAML))
	if err != nil {
		t.Fatal(err)
	}
	if s.System != "base24" {
		t.Errorf("expected system=base24, got %q", s.System)
	}
	if len(s.Palette) != 24 {
		t.Errorf("expected 24 palette entries, got %d", len(s.Palette))
	}

	// base24 schemes work with default mapping (uses only base00-base0F)
	defs, err := Bytes([]byte(draculaBase24YAML))
	if err != nil {
		t.Fatal(err)
	}
	if defs[theme.Error].Foreground != "#ff5555" {
		t.Errorf("expected error=#ff5555, got %q", defs[theme.Error].Foreground)
	}
}

func TestBase24Scheme_CustomMapping(t *testing.T) {
	m := tinted.DefaultMapping()
	m.Error = "base12" // use base24-specific brighter red
	defs, err := BytesWithMapping([]byte(draculaBase24YAML), m)
	if err != nil {
		t.Fatal(err)
	}
	if defs[theme.Error].Foreground != "#ff6e6e" {
		t.Errorf("expected error=#ff6e6e from base12, got %q", defs[theme.Error].Foreground)
	}
}
