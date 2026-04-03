package embedded

import (
	"slices"
	"testing"

	"github.com/pmarschik/lipmark/theme"
)

func TestWithEmbedded(t *testing.T) {
	reg := theme.NewWithOptions(theme.WithDefaults(), WithEmbedded())
	presets := reg.Presets()
	expected := SchemeNames()
	for _, name := range expected {
		if !slices.Contains(presets, name) {
			t.Errorf("embedded scheme %q not found in registry presets", name)
		}
	}
}

func TestWithEmbedded_SchemesLoadCorrectly(t *testing.T) {
	reg := theme.NewWithOptions(WithEmbedded())
	for _, name := range SchemeNames() {
		s := reg.Resolve(name)
		if len(s.Names()) == 0 {
			t.Errorf("scheme %q resolved to empty style set", name)
		}
	}
}

func TestSchemeNames(t *testing.T) {
	names := SchemeNames()
	if len(names) != 5 {
		t.Errorf("expected 5 embedded schemes, got %d", len(names))
	}
	if !slices.IsSorted(names) {
		t.Errorf("embedded scheme names should be sorted, got %v", names)
	}
}
