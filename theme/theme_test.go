package theme

import (
	"slices"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestNew_CreatesEmptyRegistry(t *testing.T) {
	r := New()
	if r == nil {
		t.Fatal("New returned nil")
	}
	if len(r.Presets()) != 0 {
		t.Errorf("expected 0 presets, got %d", len(r.Presets()))
	}
}

func TestRegister_AddsStyles(t *testing.T) {
	r := New()
	r.Register("default", map[string]StyleDef{
		"bold": {Bold: true},
	})
	s := r.Resolve("default")
	if s.Get("bold").Render("x") == "" {
		t.Error("expected bold style to render non-empty")
	}
}

func TestRegister_MergesStyles(t *testing.T) {
	r := New()
	r.Register("default", map[string]StyleDef{"a": {Bold: true}})
	r.Register("default", map[string]StyleDef{"b": {Italic: true}})
	s := r.Resolve("default")
	if len(s.Names()) != 2 {
		t.Errorf("expected 2 styles, got %d", len(s.Names()))
	}
}

func TestRegisterStruct(t *testing.T) {
	type appStyles struct {
		IssueKey StyleDef `style:"issue_key"`
		Repo     StyleDef `style:"repo"`
	}
	r := NewWithOptions(WithDefaults())
	r.RegisterStruct("default", appStyles{
		IssueKey: StyleDef{Foreground: "#2563EB", Bold: true},
		Repo:     StyleDef{Bold: true},
	})
	s := r.Resolve("default")
	if s.Get("issue_key").Render("x") == "" {
		t.Error("expected issue_key style to render non-empty")
	}
	if s.Get("repo").Render("x") == "" {
		t.Error("expected repo style to render non-empty")
	}
}

func TestRegisterStruct_SkipsZero(t *testing.T) {
	type appStyles struct {
		IssueKey StyleDef `style:"issue_key"`
		Repo     StyleDef `style:"repo"`
	}
	r := New()
	r.RegisterStruct("default", appStyles{
		IssueKey: StyleDef{Foreground: "#2563EB"},
	})
	s := r.Resolve("default")
	names := s.Names()
	if slices.Contains(names, "repo") {
		t.Error("zero-value style should not be registered")
	}
}

func TestResolve_ReturnsCorrectStyles(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	s := r.Resolve("default")
	// Bold should render something styled.
	got := s.Get(Bold).Render("test")
	if got == "" {
		t.Error("Bold style returned empty render")
	}
}

func TestResolve_NoColorPreset(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	s := r.Resolve("nocolor")
	// All nocolor styles should render plain text.
	names := []string{Bold, Dim, Italic, Underline, Success, Error, Warning, Info, Muted, Command, Flag, Heading, Key, Value, Path}
	for _, name := range names {
		got := s.Get(name).Render("hello")
		if got != "hello" {
			t.Errorf("nocolor %q: expected plain %q, got %q", name, "hello", got)
		}
	}
}

func TestResolve_UnknownPreset(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	s := r.Resolve("nonexistent")
	// Should return empty set with no styles.
	if len(s.Names()) != 0 {
		t.Errorf("expected 0 styles for unknown preset, got %d", len(s.Names()))
	}
}

func TestResolve_Auto(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "")
	r := NewWithOptions(WithDefaults())
	s := r.Resolve("auto")
	if s.Preset() == "" {
		t.Error("auto resolve returned empty preset")
	}
}

func TestResolveInto(t *testing.T) {
	type appStyles struct {
		IssueKey StyleDef `style:"issue_key"`
	}
	r := NewWithOptions(WithDefaults())
	r.RegisterStruct("default", appStyles{
		IssueKey: StyleDef{Foreground: "#2563EB", Bold: true},
	})

	type resolved struct {
		IssueKey lipgloss.Style `style:"issue_key"`
		Bold     lipgloss.Style `style:"bold"`
	}
	var s resolved
	if err := r.ResolveInto("default", &s); err != nil {
		t.Fatal(err)
	}
	if s.IssueKey.Render("x") == "" {
		t.Error("IssueKey should be resolved")
	}
	if s.Bold.Render("x") == "" {
		t.Error("Bold should be resolved from base")
	}
}

func TestResolveInto_NonPointerErrors(t *testing.T) {
	r := New()
	type s struct {
		Bold lipgloss.Style `style:"bold"`
	}
	var v s
	err := r.ResolveInto("default", v) // not a pointer
	if err == nil {
		t.Error("expected error for non-pointer")
	}
}

func TestAutoPreset_Default(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "")
	if got := AutoPreset(); got != "default" {
		t.Errorf("expected 'default', got %q", got)
	}
}

func TestAutoPreset_NoColorParam(t *testing.T) {
	if got := AutoPreset(true); got != "nocolor" {
		t.Errorf("expected 'nocolor', got %q", got)
	}
}

func TestAutoPreset_EnvNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	if got := AutoPreset(); got != "nocolor" {
		t.Errorf("expected 'nocolor' with NO_COLOR set, got %q", got)
	}
}

func TestAutoPreset_EnvTermDumb(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "dumb")
	if got := AutoPreset(); got != "nocolor" {
		t.Errorf("expected 'nocolor' with TERM=dumb, got %q", got)
	}
}

func TestWithDefaults_PopulatesBothPresets(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	presets := r.Presets()
	if len(presets) != 2 {
		t.Fatalf("expected 2 presets, got %d: %v", len(presets), presets)
	}
	if !slices.Contains(presets, "default") {
		t.Error("missing 'default' preset")
	}
	if !slices.Contains(presets, "nocolor") {
		t.Error("missing 'nocolor' preset")
	}
}

func TestWithDefaults_DefaultPresetHasStandardStyles(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	s := r.Resolve("default")
	expected := []string{Bold, Dim, Italic, Underline, Success, Error, Warning, Info, Muted, Command, Flag, Heading, Key, Value, Path}
	for _, name := range expected {
		if !slices.Contains(s.Names(), name) {
			t.Errorf("default preset missing style %q", name)
		}
	}
}

func TestStyleDef_ToStyle(t *testing.T) {
	d := StyleDef{
		Foreground: "#FF0000",
		Bold:       true,
		Italic:     true,
	}
	s := d.ToStyle()
	rendered := s.Render("x")
	if rendered == "" {
		t.Error("ToStyle returned style that renders empty")
	}
}

func TestStyleDef_ToStyle_Empty(t *testing.T) {
	d := StyleDef{}
	s := d.ToStyle()
	got := s.Render("hello")
	if got != "hello" {
		t.Errorf("empty StyleDef should render plain, got %q", got)
	}
}

func TestLoadDefs_FromYAMLReader(t *testing.T) {
	yamlData := `
bold:
  bold: true
success:
  foreground: "#00FF00"
  bold: true
command:
  foreground: "#00FFFF"
`
	r := strings.NewReader(yamlData)
	defs, err := LoadDefs(r)
	if err != nil {
		t.Fatalf("LoadDefs failed: %v", err)
	}
	if len(defs) != 3 {
		t.Errorf("expected 3 defs, got %d", len(defs))
	}
	if defs["success"].Foreground != "#00FF00" {
		t.Errorf("expected success foreground #00FF00, got %q", defs["success"].Foreground)
	}
}

func TestLoadDefs_InvalidYAML(t *testing.T) {
	r := strings.NewReader("{{invalid yaml")
	_, err := LoadDefs(r)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestMergeDefs_OverlaysPartial(t *testing.T) {
	base := map[string]StyleDef{
		"bold":    {Bold: true},
		"success": {Foreground: "#16A34A"},
	}
	overlay := map[string]StyleDef{
		"success": {Foreground: "#00FF00"},
	}
	merged := MergeDefs(base, overlay)

	if merged["success"].Foreground != "#00FF00" {
		t.Errorf("expected overlay foreground, got %q", merged["success"].Foreground)
	}
	if !merged["bold"].Bold {
		t.Error("expected base bold to remain")
	}
}

func TestMergeDefs_EmptyOverlay(t *testing.T) {
	base := map[string]StyleDef{
		"success": {Foreground: "#16A34A"},
	}
	merged := MergeDefs(base, map[string]StyleDef{})
	if merged["success"].Foreground != base["success"].Foreground {
		t.Error("empty overlay should not change base")
	}
}

func TestSet_Get_Unknown(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	s := r.Resolve("default")
	got := s.Get("nonexistent").Render("x")
	if got != "x" {
		t.Errorf("unknown style should render plain, got %q", got)
	}
}

func TestPresets_Sorted(t *testing.T) {
	r := NewWithOptions(WithDefaults())
	presets := r.Presets()
	if !slices.IsSorted(presets) {
		t.Errorf("presets should be sorted, got %v", presets)
	}
}
