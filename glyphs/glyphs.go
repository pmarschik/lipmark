package glyphs

import (
	"fmt"
	"maps"
	"os"
	"reflect"
	"sort"
	"strings"
)

// Registry holds glyph mappings per mode.
// Apps register their own glyph names and per-mode characters.
type Registry struct {
	modes map[string]map[string]string // mode -> name -> glyph char
}

// New creates an empty Registry.
func New() *Registry {
	return &Registry{
		modes: make(map[string]map[string]string),
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

// Register adds or overwrites glyph mappings for a mode.
// Existing glyphs in that mode are preserved; new ones are added/overwritten.
func (r *Registry) Register(mode string, glyphs map[string]string) {
	if _, ok := r.modes[mode]; !ok {
		r.modes[mode] = make(map[string]string)
	}
	maps.Copy(r.modes[mode], glyphs)
}

// Resolve returns a Set for the given mode. If mode is "auto", it detects
// the best mode from environment. Falls back to "unicode" for unknown modes.
func (r *Registry) Resolve(mode string) Set {
	mode = strings.ToLower(strings.TrimSpace(mode))
	switch mode {
	case "ascii", "plain":
		mode = "ascii"
	case "nerd", "nerdfont", "nerdfonts":
		mode = "nerd"
	case "emoji":
		mode = "emoji"
	case "unicode":
		mode = "unicode"
	case "", "auto":
		mode = AutoMode()
		return r.Resolve(mode)
	default:
		mode = "unicode"
	}
	glyphs := make(map[string]string)
	if m, ok := r.modes[mode]; ok {
		maps.Copy(glyphs, m)
	}
	return Set{mode: mode, glyphs: glyphs}
}

// Modes returns the list of registered mode names.
func (r *Registry) Modes() []string {
	modes := make([]string, 0, len(r.modes))
	for m := range r.modes {
		modes = append(modes, m)
	}
	sort.Strings(modes)
	return modes
}

// Set is a resolved glyph set for a specific mode.
type Set struct {
	glyphs map[string]string
	mode   string
}

// Get returns the glyph character for the given name, or "" if not defined.
func (s Set) Get(name string) string {
	return s.glyphs[name]
}

// Mode returns the mode name of this set.
func (s Set) Mode() string {
	return s.mode
}

// Names returns all glyph names in this set.
func (s Set) Names() []string {
	names := make([]string, 0, len(s.glyphs))
	for k := range s.glyphs {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// Standard glyph name constants.
const (
	// Status indicators.
	Success = "success"
	Error   = "error"
	Warning = "warning"
	Info    = "info"
	Check   = "check"
	Cross   = "cross"

	// Navigation / structure.
	Arrow     = "arrow"
	ArrowLeft = "arrow_left"
	ArrowUp   = "arrow_up"
	ArrowDown = "arrow_down"
	Bullet    = "bullet"
	Ellipsis  = "ellipsis"
	Separator = "separator"

	// Common actions / objects.
	Star     = "star"
	Heart    = "heart"
	Folder   = "folder"
	File     = "file"
	Lock     = "lock"
	Unlock   = "unlock"
	Edit     = "edit"
	Trash    = "trash"
	Search   = "search"
	Refresh  = "refresh"
	Download = "download"
	Upload   = "upload"
	Clock    = "clock"
	Play     = "play"
	Pause    = "pause"
	Stop     = "stop"
	Gear     = "gear"
	Link     = "link"
	User     = "user"
	Plus     = "plus"
	Minus    = "minus"
)

// BaseGlyphs provides typed access to all standard resolved glyphs.
// Embed this in your app's glyph struct for type-safe access:
//
//	type MyGlyphs struct {
//	    glyphs.BaseGlyphs
//	    Dirty    string `glyph:"dirty"`
//	    Unpushed string `glyph:"unpushed"`
//	}
type BaseGlyphs struct {
	Success   string `glyph:"success"`
	Error     string `glyph:"error"`
	Warning   string `glyph:"warning"`
	Info      string `glyph:"info"`
	Check     string `glyph:"check"`
	Cross     string `glyph:"cross"`
	Arrow     string `glyph:"arrow"`
	ArrowLeft string `glyph:"arrow_left"`
	ArrowUp   string `glyph:"arrow_up"`
	ArrowDown string `glyph:"arrow_down"`
	Bullet    string `glyph:"bullet"`
	Ellipsis  string `glyph:"ellipsis"`
	Separator string `glyph:"separator"`
	Star      string `glyph:"star"`
	Heart     string `glyph:"heart"`
	Folder    string `glyph:"folder"`
	File      string `glyph:"file"`
	Lock      string `glyph:"lock"`
	Unlock    string `glyph:"unlock"`
	Edit      string `glyph:"edit"`
	Trash     string `glyph:"trash"`
	Search    string `glyph:"search"`
	Refresh   string `glyph:"refresh"`
	Download  string `glyph:"download"`
	Upload    string `glyph:"upload"`
	Clock     string `glyph:"clock"`
	Play      string `glyph:"play"`
	Pause     string `glyph:"pause"`
	Stop      string `glyph:"stop"`
	Gear      string `glyph:"gear"`
	Link      string `glyph:"link"`
	User      string `glyph:"user"`
	Plus      string `glyph:"plus"`
	Minus     string `glyph:"minus"`
}

// WithDefaults registers the base glyphs across all four built-in modes
// (unicode, ascii, nerd, emoji).
func WithDefaults() Option {
	return func(r *Registry) {
		r.Register("unicode", map[string]string{
			Success: "\u2713", Error: "\u2717", Warning: "\u26A0", Info: "\u2139",
			Check: "\u2713", Cross: "\u2717",
			Arrow: "\u2192", ArrowLeft: "\u2190", ArrowUp: "\u2191", ArrowDown: "\u2193",
			Bullet: "\u2022", Ellipsis: "\u2026", Separator: "\u2502",
			Star: "\u2605", Heart: "\u2665",
			Folder: "\u25B8", File: "\u25AA", Lock: "\u25CF", Unlock: "\u25CB",
			Edit: "\u270E", Trash: "\u2716", Search: "\u25C9",
			Refresh: "\u21BB", Download: "\u2193", Upload: "\u2191",
			Clock: "\u25F7", Play: "\u25B6", Pause: "\u23F8", Stop: "\u25A0",
			Gear: "\u2699", Link: "\u2197", User: "\u25C8",
			Plus: "+", Minus: "\u2212",
		})
		r.Register("ascii", map[string]string{
			Success: "[ok]", Error: "[err]", Warning: "[warn]", Info: "[info]",
			Check: "[ok]", Cross: "[err]",
			Arrow: "->", ArrowLeft: "<-", ArrowUp: "^", ArrowDown: "v",
			Bullet: "*", Ellipsis: "...", Separator: "|",
			Star: "*", Heart: "<3",
			Folder: ">", File: "-", Lock: "[x]", Unlock: "[ ]",
			Edit: "[~]", Trash: "[x]", Search: "[?]",
			Refresh: "[r]", Download: "[v]", Upload: "[^]",
			Clock: "[t]", Play: "[>]", Pause: "[||]", Stop: "[#]",
			Gear: "[*]", Link: "[>]", User: "[@]",
			Plus: "+", Minus: "-",
		})
		r.Register("nerd", map[string]string{
			Success: "\uf00c", Error: "\uf00d", Warning: "\uf071", Info: "\uf05a",
			Check: "\uf00c", Cross: "\uf00d",
			Arrow: "\uf061", ArrowLeft: "\uf060", ArrowUp: "\uf062", ArrowDown: "\uf063",
			Bullet: "\uf111", Ellipsis: "\uf141", Separator: "\ue0b1",
			Star: "\uf005", Heart: "\uf004",
			Folder: "\uf07b", File: "\uf15b", Lock: "\uf023", Unlock: "\uf09c",
			Edit: "\uf044", Trash: "\uf1f8", Search: "\uf002",
			Refresh: "\uf021", Download: "\uf019", Upload: "\uf093",
			Clock: "\uf017", Play: "\uf04b", Pause: "\uf04c", Stop: "\uf04d",
			Gear: "\uf013", Link: "\uf0c1", User: "\uf007",
			Plus: "\uf067", Minus: "\uf068",
		})
		r.Register("emoji", map[string]string{
			Success: "\u2705", Error: "\u274C", Warning: "\u26A0\uFE0F", Info: "\u2139\uFE0F",
			Check: "\u2705", Cross: "\u274C",
			Arrow: "\u27A1\uFE0F", ArrowLeft: "\u2B05\uFE0F", ArrowUp: "\u2B06\uFE0F", ArrowDown: "\u2B07\uFE0F",
			Bullet: "\u25AA\uFE0F", Ellipsis: "\u2026", Separator: "\u2502",
			Star: "\u2B50", Heart: "\u2764\uFE0F",
			Folder: "\U0001F4C1", File: "\U0001F4C4", Lock: "\U0001F512", Unlock: "\U0001F513",
			Edit: "\u270F\uFE0F", Trash: "\U0001F5D1\uFE0F", Search: "\U0001F50D",
			Refresh: "\U0001F504", Download: "\u2B07\uFE0F", Upload: "\u2B06\uFE0F",
			Clock: "\U0001F552", Play: "\u25B6\uFE0F", Pause: "\u23F8\uFE0F", Stop: "\u23F9\uFE0F",
			Gear: "\u2699\uFE0F", Link: "\U0001F517", User: "\U0001F464",
			Plus: "\u2795", Minus: "\u2796",
		})
	}
}

// RegisterStruct registers glyph mappings for a mode from a struct.
// Fields must be exported strings with a `glyph:"name"` tag.
// Only non-empty field values are registered.
//
//	type MyGlyphs struct {
//	    Dirty    string `glyph:"dirty"`
//	    Unpushed string `glyph:"unpushed"`
//	}
//	reg.RegisterStruct("nerd", MyGlyphs{Dirty: "", Unpushed: ""})
func (r *Registry) RegisterStruct(mode string, v any) {
	m := structToMap(v)
	r.Register(mode, m)
}

// ResolveInto resolves glyphs for the given mode and populates the target struct.
// Target must be a pointer to a struct with `glyph:"name"` tags on string fields.
// Fields without a matching glyph in the registry are left unchanged.
func (r *Registry) ResolveInto(mode string, target any) error {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("glyphs: ResolveInto requires a pointer to struct, got %T", target)
	}

	set := r.Resolve(mode)
	resolveGlyphFields(rv.Elem(), set)
	return nil
}

func resolveGlyphFields(rv reflect.Value, set Set) {
	rt := rv.Type()

	for i := range rt.NumField() {
		f := rt.Field(i)
		// Recurse into embedded structs.
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			resolveGlyphFields(rv.Field(i), set)
			continue
		}
		if !f.IsExported() || f.Type.Kind() != reflect.String {
			continue
		}
		name := f.Tag.Get("glyph")
		if name == "" || name == "-" {
			continue
		}
		if g := set.Get(name); g != "" {
			rv.Field(i).SetString(g)
		}
	}
}

func structToMap(v any) map[string]string {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil
	}

	m := make(map[string]string)
	collectGlyphDefs(rv, m)
	return m
}

func collectGlyphDefs(rv reflect.Value, m map[string]string) {
	rt := rv.Type()

	for i := range rt.NumField() {
		f := rt.Field(i)
		// Recurse into embedded structs.
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			collectGlyphDefs(rv.Field(i), m)
			continue
		}
		if !f.IsExported() || f.Type.Kind() != reflect.String {
			continue
		}
		name := f.Tag.Get("glyph")
		if name == "" || name == "-" {
			continue
		}
		val := rv.Field(i).String()
		if val != "" {
			m[name] = val
		}
	}
}

// AutoMode detects the best glyph mode from environment variables.
// Checks: NERD_FONT, POWERLINE_*, TERM_PROGRAM, locale for UTF-8.
func AutoMode() string {
	if !isUTF8Locale() {
		return "ascii"
	}
	if hasTruthyEnv("NERD_FONT", "HAVE_NERD_FONT") {
		return "nerd"
	}
	if hasTruthyEnv("POWERLINE_COMMAND", "POWERLINE_CONFIG_COMMAND", "POWERLINE_BASH_CONTINUATION") {
		return "unicode"
	}
	if supportsEmoji() {
		return "emoji"
	}
	return "unicode"
}

func supportsEmoji() bool {
	termProgram := strings.ToLower(os.Getenv("TERM_PROGRAM"))
	switch termProgram {
	case "iterm.app", "apple_terminal", "wezterm", "vscode", "warpterminal", "ghostty", "hyper", "rio":
		return true
	}
	return os.Getenv("KITTY_WINDOW_ID") != "" || os.Getenv("WT_SESSION") != ""
}

func hasTruthyEnv(keys ...string) bool {
	for _, k := range keys {
		v := strings.TrimSpace(strings.ToLower(os.Getenv(k)))
		if v == "" {
			continue
		}
		switch v {
		case "0", "false", "no", "off":
			continue
		default:
			return true
		}
	}
	return false
}

func isUTF8Locale() bool {
	locale := strings.ToUpper(strings.Join([]string{
		os.Getenv("LC_ALL"),
		os.Getenv("LC_CTYPE"),
		os.Getenv("LANG"),
	}, " "))
	return strings.Contains(locale, "UTF-8") || strings.Contains(locale, "UTF8")
}
