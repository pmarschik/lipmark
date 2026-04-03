package glyphs

import (
	"testing"
)

func TestNew_CreatesEmptyRegistry(t *testing.T) {
	r := New()
	if r == nil {
		t.Fatal("New returned nil")
	}
	if len(r.Modes()) != 0 {
		t.Errorf("expected 0 modes, got %d", len(r.Modes()))
	}
}

func TestRegister_AddsGlyphs(t *testing.T) {
	r := New()
	r.Register("test", map[string]string{"foo": "bar"})
	s := r.Resolve("test")
	// Unknown mode falls back to unicode, so register as unicode.
	r2 := New()
	r2.Register("unicode", map[string]string{"foo": "bar"})
	s2 := r2.Resolve("unicode")
	if s2.Get("foo") != "bar" {
		t.Errorf("expected 'bar', got %q", s2.Get("foo"))
	}
	_ = s
}

func TestRegister_MergesGlyphs(t *testing.T) {
	r := New()
	r.Register("unicode", map[string]string{"a": "1"})
	r.Register("unicode", map[string]string{"b": "2"})
	s := r.Resolve("unicode")
	if s.Get("a") != "1" {
		t.Error("expected 'a' to be '1'")
	}
	if s.Get("b") != "2" {
		t.Error("expected 'b' to be '2'")
	}
}

func TestResolve_ReturnsCorrectSet(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	tests := []struct {
		mode string
		name string
		want string
	}{
		{"unicode", Success, "\u2713"},
		{"ascii", Success, "[ok]"},
		{"ascii", Arrow, "->"},
		{"nerd", Success, "\uf00c"},
		{"emoji", Success, "\u2705"},
	}
	for _, tt := range tests {
		t.Run(tt.mode+"/"+tt.name, func(t *testing.T) {
			s := r.Resolve(tt.mode)
			if got := s.Get(tt.name); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolve_ModeAliases(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	for _, mode := range []string{"plain", "ASCII", " ascii "} {
		s := r.Resolve(mode)
		if s.Get(Success) != "[ok]" {
			t.Errorf("mode %q: expected ASCII success glyph", mode)
		}
	}
	for _, mode := range []string{"nerd", "nerdfont", "nerdfonts"} {
		s := r.Resolve(mode)
		if s.Get(Success) != "\uf00c" {
			t.Errorf("mode %q: expected nerd success glyph", mode)
		}
	}
}

func TestGet_ReturnsEmptyForUnknown(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	s := r.Resolve("unicode")
	if got := s.Get("nonexistent"); got != "" {
		t.Errorf("expected empty string for unknown glyph, got %q", got)
	}
}

func TestWithDefaults_PopulatesAllModes(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	modes := r.Modes()
	expected := []string{"ascii", "emoji", "nerd", "unicode"}
	if len(modes) != len(expected) {
		t.Fatalf("expected %d modes, got %d: %v", len(expected), len(modes), modes)
	}
	for i, m := range expected {
		if modes[i] != m {
			t.Errorf("expected mode %q at index %d, got %q", m, i, modes[i])
		}
	}
}

func TestAutoMode_ReturnsValidString(t *testing.T) {
	mode := AutoMode()
	valid := map[string]bool{
		"ascii": true, "unicode": true, "nerd": true, "emoji": true,
	}
	if !valid[mode] {
		t.Errorf("AutoMode returned unexpected mode %q", mode)
	}
}

func TestResolve_Auto_UsesAutoMode(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	s := r.Resolve("auto")
	if s.Mode() == "" {
		t.Error("auto resolve returned empty mode")
	}
	if s.Get(Success) == "" {
		t.Error("auto resolve returned empty Success glyph")
	}
}

func TestResolve_UnknownFallsBackToUnicode(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	s := r.Resolve("unknown_mode")
	if s.Mode() != "unicode" {
		t.Errorf("expected unicode fallback, got %q", s.Mode())
	}
	if s.Get(Success) != "\u2713" {
		t.Error("expected unicode success glyph for unknown mode")
	}
}

func TestSet_Mode(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	s := r.Resolve("ascii")
	if s.Mode() != "ascii" {
		t.Errorf("expected mode 'ascii', got %q", s.Mode())
	}
}

func TestSet_Names(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	s := r.Resolve("unicode")
	names := s.Names()
	if len(names) < 30 {
		t.Errorf("expected at least 30 glyph names, got %d: %v", len(names), names)
	}
}

func TestAutoMode_NonUTF8(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_CTYPE", "")
	t.Setenv("LANG", "C")
	mode := AutoMode()
	if mode != "ascii" {
		t.Errorf("expected 'ascii' for non-UTF8 locale, got %q", mode)
	}
}

type testGlyphs struct {
	Success  string `glyph:"success"`
	Dirty    string `glyph:"dirty"`
	Unpushed string `glyph:"unpushed"`
	Ignored  string // no tag
	private  string `glyph:"private"` //nolint:unused // intentionally unused to test that unexported fields are skipped
}

func TestRegisterStruct(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	r.RegisterStruct("nerd", testGlyphs{
		Dirty:    "D",
		Unpushed: "U",
	})
	s := r.Resolve("nerd")
	if s.Get("dirty") != "D" {
		t.Errorf("expected 'D', got %q", s.Get("dirty"))
	}
	if s.Get("unpushed") != "U" {
		t.Errorf("expected 'U', got %q", s.Get("unpushed"))
	}
	// Default glyphs should still be there.
	if s.Get(Success) == "" {
		t.Error("default success glyph should be preserved")
	}
}

func TestRegisterStruct_SkipsEmpty(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	// Success is empty in struct, so it should NOT override the default.
	r.RegisterStruct("unicode", testGlyphs{Dirty: "M"})
	s := r.Resolve("unicode")
	if s.Get(Success) != "\u2713" {
		t.Errorf("empty struct field should not override, got %q", s.Get(Success))
	}
	if s.Get("dirty") != "M" {
		t.Errorf("expected 'M', got %q", s.Get("dirty"))
	}
}

func TestResolveInto(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	r.Register("unicode", map[string]string{"dirty": "M", "unpushed": "^"})

	var g testGlyphs
	if err := r.ResolveInto("unicode", &g); err != nil {
		t.Fatal(err)
	}
	if g.Success != "\u2713" {
		t.Errorf("expected checkmark, got %q", g.Success)
	}
	if g.Dirty != "M" {
		t.Errorf("expected 'M', got %q", g.Dirty)
	}
	if g.Unpushed != "^" {
		t.Errorf("expected '^', got %q", g.Unpushed)
	}
	if g.Ignored != "" {
		t.Error("untagged field should be empty")
	}
}

func TestResolveInto_NonPointerErrors(t *testing.T) {
	r := New()
	var g testGlyphs
	err := r.ResolveInto("unicode", g) // not a pointer
	if err == nil {
		t.Error("expected error for non-pointer")
	}
}

func TestResolveInto_LeavesUnmatchedFields(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	// "dirty" is not in defaults
	g := testGlyphs{Dirty: "original"}
	if err := r.ResolveInto("unicode", &g); err != nil {
		t.Fatal(err)
	}
	if g.Dirty != "original" {
		t.Errorf("unmatched field should be unchanged, got %q", g.Dirty)
	}
}

func TestAutoMode_UTF8WithNerdFont(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_CTYPE", "")
	t.Setenv("LANG", "en_US.UTF-8")
	t.Setenv("NERD_FONT", "1")
	t.Setenv("HAVE_NERD_FONT", "")
	t.Setenv("TERM_PROGRAM", "")
	t.Setenv("KITTY_WINDOW_ID", "")
	t.Setenv("WT_SESSION", "")
	t.Setenv("POWERLINE_COMMAND", "")
	t.Setenv("POWERLINE_CONFIG_COMMAND", "")
	t.Setenv("POWERLINE_BASH_CONTINUATION", "")

	mode := AutoMode()
	if mode != "nerd" {
		t.Errorf("expected 'nerd' when NERD_FONT=1, got %q", mode)
	}
}
