package theme

import (
	"fmt"
	"io"
	"maps"
	"os"
	"reflect"
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
	"gopkg.in/yaml.v3"
)

// Standard style name constants.
const (
	Bold      = "bold"
	Dim       = "dim"
	Italic    = "italic"
	Underline = "underline"
	Success   = "success"
	Error     = "error"
	Warning   = "warning"
	Info      = "info"
	Muted     = "muted"
	Command   = "command"
	Flag      = "flag"
	Heading   = "heading"
	Key       = "key"
	Value     = "value"
	Path      = "path"
)

// BaseStyles provides typed access to the standard resolved styles.
// Embed this in your app's theme struct for type-safe style access:
//
//	type MyTheme struct {
//	    theme.BaseStyles
//	    IssueKey lipgloss.Style `style:"issue_key"`
//	}
//	var t MyTheme
//	reg.ResolveInto("auto", &t)
//	t.Success.Render("ok")  // from BaseStyles
//	t.IssueKey.Render("X")  // app-specific
type BaseStyles struct {
	Bold      lipgloss.Style `style:"bold"`
	Dim       lipgloss.Style `style:"dim"`
	Italic    lipgloss.Style `style:"italic"`
	Underline lipgloss.Style `style:"underline"`
	Success   lipgloss.Style `style:"success"`
	Error     lipgloss.Style `style:"error"`
	Warning   lipgloss.Style `style:"warning"`
	Info      lipgloss.Style `style:"info"`
	Muted     lipgloss.Style `style:"muted"`
	Command   lipgloss.Style `style:"command"`
	Flag      lipgloss.Style `style:"flag"`
	Heading   lipgloss.Style `style:"heading"`
	Key       lipgloss.Style `style:"key"`
	Value     lipgloss.Style `style:"value"`
	Path      lipgloss.Style `style:"path"`
}

// BaseStyleDefs provides typed access to the standard style definitions.
// Embed this in your app's definition struct for type-safe registration:
//
//	type MyDefs struct {
//	    theme.BaseStyleDefs
//	    IssueKey theme.StyleDef `style:"issue_key"`
//	}
type BaseStyleDefs struct {
	Bold      StyleDef `style:"bold"`
	Dim       StyleDef `style:"dim"`
	Italic    StyleDef `style:"italic"`
	Underline StyleDef `style:"underline"`
	Success   StyleDef `style:"success"`
	Error     StyleDef `style:"error"`
	Warning   StyleDef `style:"warning"`
	Info      StyleDef `style:"info"`
	Muted     StyleDef `style:"muted"`
	Command   StyleDef `style:"command"`
	Flag      StyleDef `style:"flag"`
	Heading   StyleDef `style:"heading"`
	Key       StyleDef `style:"key"`
	Value     StyleDef `style:"value"`
	Path      StyleDef `style:"path"`
}

// StyleDef is a serializable style definition that can be loaded from config files.
type StyleDef struct {
	Foreground string `yaml:"foreground,omitempty" json:"foreground,omitempty"`
	Background string `yaml:"background,omitempty" json:"background,omitempty"`
	Bold       bool   `yaml:"bold,omitempty" json:"bold,omitempty"`
	Faint      bool   `yaml:"faint,omitempty" json:"faint,omitempty"`
	Italic     bool   `yaml:"italic,omitempty" json:"italic,omitempty"`
	Underline  bool   `yaml:"underline,omitempty" json:"underline,omitempty"`
}

// ToStyle converts a StyleDef to a lipgloss.Style.
func (d StyleDef) ToStyle() lipgloss.Style {
	s := lipgloss.NewStyle()
	if d.Foreground != "" {
		s = s.Foreground(lipgloss.Color(d.Foreground))
	}
	if d.Background != "" {
		s = s.Background(lipgloss.Color(d.Background))
	}
	if d.Bold {
		s = s.Bold(true)
	}
	if d.Faint {
		s = s.Faint(true)
	}
	if d.Italic {
		s = s.Italic(true)
	}
	if d.Underline {
		s = s.Underline(true)
	}
	return s
}

// Registry holds style definitions across named presets.
type Registry struct {
	presets map[string]map[string]StyleDef // preset -> name -> StyleDef
}

// New creates an empty Registry.
func New() *Registry {
	return &Registry{
		presets: make(map[string]map[string]StyleDef),
	}
}

// Option configures a Registry.
type Option func(*Registry)

// NewWithOptions creates a Registry with options applied.
func NewWithOptions(opts ...Option) *Registry {
	r := New()
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Register adds or overwrites style defs for a preset.
// Existing styles in that preset are preserved; new ones are added/overwritten.
func (r *Registry) Register(preset string, styles map[string]StyleDef) {
	if _, ok := r.presets[preset]; !ok {
		r.presets[preset] = make(map[string]StyleDef)
	}
	maps.Copy(r.presets[preset], styles)
}

// RegisterStruct registers styles from a struct with `style:"name"` tags.
// Only non-zero StyleDef fields are registered.
func (r *Registry) RegisterStruct(preset string, v any) {
	m := structToStyleMap(v)
	r.Register(preset, m)
}

// Resolve returns a Set for the given preset.
// "auto" detects NO_COLOR/TERM=dumb -> "nocolor", otherwise "default".
func (r *Registry) Resolve(preset string) Set {
	if strings.EqualFold(preset, "auto") || preset == "" {
		preset = AutoPreset()
	}
	styles := make(map[string]lipgloss.Style)
	if m, ok := r.presets[preset]; ok {
		for k, v := range m {
			styles[k] = v.ToStyle()
		}
	}
	return Set{preset: preset, styles: styles}
}

// ResolveInto resolves a preset into a target struct with `style:"name"` tags
// on lipgloss.Style fields.
func (r *Registry) ResolveInto(preset string, target any) error {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("theme: ResolveInto requires a pointer to struct, got %T", target)
	}

	set := r.Resolve(preset)
	resolveStyleFields(rv.Elem(), set)
	return nil
}

func resolveStyleFields(rv reflect.Value, set Set) {
	rt := rv.Type()
	styleType := reflect.TypeFor[lipgloss.Style]()

	for i := range rt.NumField() {
		f := rt.Field(i)
		// Recurse into embedded structs.
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			resolveStyleFields(rv.Field(i), set)
			continue
		}
		if !f.IsExported() || f.Type != styleType {
			continue
		}
		name := f.Tag.Get("style")
		if name == "" || name == "-" {
			continue
		}
		rv.Field(i).Set(reflect.ValueOf(set.Get(name)))
	}
}

// Presets returns registered preset names.
func (r *Registry) Presets() []string {
	presets := make([]string, 0, len(r.presets))
	for p := range r.presets {
		presets = append(presets, p)
	}
	sort.Strings(presets)
	return presets
}

// Set is a resolved set of lipgloss styles.
type Set struct {
	styles map[string]lipgloss.Style
	preset string
}

// Get returns the style for the given name, or a plain style if not found.
func (s Set) Get(name string) lipgloss.Style {
	if st, ok := s.styles[name]; ok {
		return st
	}
	return lipgloss.NewStyle()
}

// Preset returns the preset name of this set.
func (s Set) Preset() string {
	return s.preset
}

// Names returns all style names in this set.
func (s Set) Names() []string {
	names := make([]string, 0, len(s.styles))
	for k := range s.styles {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// DefaultNoColorEnvVars is the default set of environment variables checked for no-color mode.
var DefaultNoColorEnvVars = []string{"NO_COLOR"}

// AutoPreset detects the best preset: "nocolor" if any no-color env var is set or TERM=dumb,
// else "default". Pass noColor=true to force no-color regardless of environment.
func AutoPreset(noColor ...bool) string {
	return AutoPresetWithEnv(DefaultNoColorEnvVars, noColor...)
}

// AutoPresetWithEnv is like AutoPreset but checks the given env var names instead of the default.
// This allows apps to use custom env vars (e.g. "MYAPP_NO_COLOR") alongside or instead of "NO_COLOR".
func AutoPresetWithEnv(envVars []string, noColor ...bool) string {
	if len(noColor) > 0 && noColor[0] {
		return "nocolor"
	}
	if IsNoColor(envVars...) {
		return "nocolor"
	}
	return "default"
}

// IsNoColor returns true if color should be disabled based on the given env var names
// (or DefaultNoColorEnvVars if none given) and TERM=dumb.
func IsNoColor(envVars ...string) bool {
	if len(envVars) == 0 {
		envVars = DefaultNoColorEnvVars
	}
	for _, v := range envVars {
		if os.Getenv(v) != "" {
			return true
		}
	}
	return os.Getenv("TERM") == "dumb"
}

// WithDefaults registers "default" (charm colors) and "nocolor" (plain) presets
// with the standard style names.
func WithDefaults() Option {
	return func(r *Registry) {
		r.Register("default", map[string]StyleDef{
			Bold:      {Bold: true},
			Dim:       {Faint: true},
			Italic:    {Italic: true},
			Underline: {Underline: true},

			Success: {Foreground: "#16A34A", Bold: true},
			Error:   {Foreground: "#DC2626", Bold: true},
			Warning: {Foreground: "#CA8A04", Bold: true},
			Info:    {Foreground: "#2563EB"},
			Muted:   {Faint: true},

			Command: {Foreground: "#06B6D4"},
			Flag:    {Foreground: "#6B7280"},
			Heading: {Foreground: "#7C3AED", Bold: true},
			Key:     {Foreground: "#06B6D4"},
			Value:   {},
			Path:    {Faint: true},
		})
		r.Register("nocolor", map[string]StyleDef{
			Bold: {}, Dim: {}, Italic: {}, Underline: {},
			Success: {}, Error: {}, Warning: {}, Info: {}, Muted: {},
			Command: {}, Flag: {}, Heading: {}, Key: {},
			Value: {}, Path: {},
		})
	}
}

// LoadDefs reads style definitions from a reader (YAML format).
func LoadDefs(r io.Reader) (map[string]StyleDef, error) {
	var raw map[string]StyleDef
	dec := yaml.NewDecoder(r)
	if err := dec.Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// LoadDefsFile loads style definitions from a YAML file path.
func LoadDefsFile(path string) (map[string]StyleDef, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadDefs(f)
}

// MergeDefs overlays non-zero fields from overlay onto base.
// This lets users provide partial theme files that override specific styles.
func MergeDefs(base, overlay map[string]StyleDef) map[string]StyleDef {
	result := make(map[string]StyleDef, len(base))
	maps.Copy(result, base)
	maps.Copy(result, overlay)
	return result
}

func structToStyleMap(v any) map[string]StyleDef {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil
	}

	m := make(map[string]StyleDef)
	collectStyleDefs(rv, m)
	return m
}

func collectStyleDefs(rv reflect.Value, m map[string]StyleDef) {
	rt := rv.Type()
	zero := StyleDef{}

	for i := range rt.NumField() {
		f := rt.Field(i)
		// Recurse into embedded structs.
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			collectStyleDefs(rv.Field(i), m)
			continue
		}
		if !f.IsExported() || f.Type != reflect.TypeFor[StyleDef]() {
			continue
		}
		name := f.Tag.Get("style")
		if name == "" || name == "-" {
			continue
		}
		val, ok := rv.Field(i).Interface().(StyleDef)
		if ok && val != zero {
			m[name] = val
		}
	}
}
