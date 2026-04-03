package tinted

import (
	"testing"
)

func TestDefaultMapping(t *testing.T) {
	m := DefaultMapping()
	if m.Error != "base08" {
		t.Errorf("expected Error=base08, got %q", m.Error)
	}
	if m.Success != "base0B" {
		t.Errorf("expected Success=base0B, got %q", m.Success)
	}
}

func TestPaletteConstants(t *testing.T) {
	if Base00 != "base00" {
		t.Errorf("expected Base00=base00, got %q", Base00)
	}
	if Base0F != "base0F" {
		t.Errorf("expected Base0F=base0F, got %q", Base0F)
	}
	if Base17 != "base17" {
		t.Errorf("expected Base17=base17, got %q", Base17)
	}
}
