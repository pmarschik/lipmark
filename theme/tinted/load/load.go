// Package load provides YAML-based loading of tinted-theming color schemes.
//
// It converts base16 and base24 scheme YAML files into [theme.StyleDef] maps
// using the mapping from the [tinted] package.
//
// # Loading from file
//
//	defs, err := load.File("dracula.yaml")
//	themeReg.Register("dracula", defs)
//
// # Loading from bytes
//
//	defs, err := load.Bytes(yamlData)
package load

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/pmarschik/lipmark/theme"
	"github.com/pmarschik/lipmark/theme/tinted"
	"gopkg.in/yaml.v3"
)

// Scheme represents a parsed tinted-theming color scheme (base16 or base24)
// with YAML tags for deserialization.
type Scheme struct {
	Palette map[string]string `yaml:"palette"`
	System  string            `yaml:"system"`
	Name    string            `yaml:"name"`
	Author  string            `yaml:"author"`
	Variant string            `yaml:"variant"`
}

// File loads a scheme YAML file and converts it to StyleDefs.
func File(path string) (map[string]theme.StyleDef, error) {
	return FileWithMapping(path, tinted.DefaultMapping())
}

// FileWithMapping loads a scheme YAML file with a custom mapping.
func FileWithMapping(path string, m tinted.Mapping) (map[string]theme.StyleDef, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return decodeAndConvert(f, m)
}

// Bytes loads a scheme from YAML bytes.
func Bytes(data []byte) (map[string]theme.StyleDef, error) {
	return BytesWithMapping(data, tinted.DefaultMapping())
}

// BytesWithMapping loads a scheme from YAML bytes with a custom mapping.
func BytesWithMapping(data []byte, m tinted.Mapping) (map[string]theme.StyleDef, error) {
	var s Scheme
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("tinted: parse scheme: %w", err)
	}
	return convert(s, m), nil
}

// Reader loads a scheme from a reader.
func Reader(r io.Reader) (map[string]theme.StyleDef, error) {
	return decodeAndConvert(r, tinted.DefaultMapping())
}

// ReaderWithMapping loads a scheme from a reader with a custom mapping.
func ReaderWithMapping(r io.Reader, m tinted.Mapping) (map[string]theme.StyleDef, error) {
	return decodeAndConvert(r, m)
}

// ParseScheme parses the raw scheme without converting to StyleDefs.
// Useful for inspecting the palette or building custom mappings.
func ParseScheme(r io.Reader) (*Scheme, error) {
	var s Scheme
	if err := yaml.NewDecoder(r).Decode(&s); err != nil {
		return nil, fmt.Errorf("tinted: parse scheme: %w", err)
	}
	return &s, nil
}

// Dir loads all .yaml files from a directory as schemes.
// Returns a map of scheme name (filename without extension) to StyleDefs.
func Dir(dir string) (map[string]map[string]theme.StyleDef, error) {
	return DirWithMapping(dir, tinted.DefaultMapping())
}

// DirWithMapping loads all .yaml files from a directory with a custom mapping.
func DirWithMapping(dir string, m tinted.Mapping) (map[string]map[string]theme.StyleDef, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("tinted: read dir: %w", err)
	}
	schemes := make(map[string]map[string]theme.StyleDef)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".yaml")
		path := dir + "/" + e.Name()
		defs, loadErr := FileWithMapping(path, m)
		if loadErr != nil {
			return nil, fmt.Errorf("tinted: load %s: %w", path, loadErr)
		}
		schemes[name] = defs
	}
	return schemes, nil
}

// DirFS loads all .yaml files from an [fs.FS] under dir as schemes.
func DirFS(fsys fs.FS, dir string) (map[string]map[string]theme.StyleDef, error) {
	return DirFSWithMapping(fsys, dir, tinted.DefaultMapping())
}

// DirFSWithMapping loads all .yaml files from an [fs.FS] under dir with a custom mapping.
func DirFSWithMapping(fsys fs.FS, dir string, m tinted.Mapping) (map[string]map[string]theme.StyleDef, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("tinted: read fs: %w", err)
	}
	schemes := make(map[string]map[string]theme.StyleDef)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".yaml")
		fpath := dir + "/" + e.Name()
		f, openErr := fsys.Open(fpath)
		if openErr != nil {
			return nil, fmt.Errorf("tinted: open %s: %w", fpath, openErr)
		}
		defs, loadErr := decodeAndConvert(f, m)
		f.Close()
		if loadErr != nil {
			return nil, fmt.Errorf("tinted: load %s: %w", fpath, loadErr)
		}
		schemes[name] = defs
	}
	return schemes, nil
}

func decodeAndConvert(r io.Reader, m tinted.Mapping) (map[string]theme.StyleDef, error) {
	s, err := ParseScheme(r)
	if err != nil {
		return nil, err
	}
	return convert(*s, m), nil
}

func convert(s Scheme, m tinted.Mapping) map[string]theme.StyleDef {
	defs := make(map[string]theme.StyleDef)
	mapStyle := func(name, paletteKey string, bold, faint, italic bool) {
		if paletteKey == "" {
			return
		}
		color, ok := s.Palette[paletteKey]
		if !ok {
			return
		}
		defs[name] = theme.StyleDef{
			Foreground: color,
			Bold:       bold,
			Faint:      faint,
			Italic:     italic,
		}
	}

	mapStyle(theme.Bold, m.Bold, true, false, false)
	mapStyle(theme.Dim, m.Dim, false, true, false)
	mapStyle(theme.Italic, m.Italic, false, false, true)
	if m.Underline != "" {
		if color, ok := s.Palette[m.Underline]; ok {
			defs[theme.Underline] = theme.StyleDef{Foreground: color, Underline: true}
		}
	}
	mapStyle(theme.Success, m.Success, true, false, false)
	mapStyle(theme.Error, m.Error, true, false, false)
	mapStyle(theme.Warning, m.Warning, true, false, false)
	mapStyle(theme.Info, m.Info, false, false, false)
	mapStyle(theme.Muted, m.Muted, false, true, false)
	mapStyle(theme.Command, m.Command, false, false, false)
	mapStyle(theme.Flag, m.Flag, false, false, false)
	mapStyle(theme.Heading, m.Heading, true, false, false)
	mapStyle(theme.Key, m.Key, false, false, false)
	mapStyle(theme.Value, m.Value, false, false, false)
	mapStyle(theme.Path, m.Path, false, true, false)

	return defs
}
