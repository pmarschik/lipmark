// Package embedded provides pre-generated tinted-theming color schemes.
//
// Use [WithEmbedded] to register the bundled schemes (catppuccin-mocha,
// dracula, nord, solarized-dark, solarized-light) on a [theme.Registry]:
//
//	reg := theme.NewWithOptions(theme.WithDefaults(), embedded.WithEmbedded())
package embedded

//go:generate go run github.com/pmarschik/lipmark/cmd/lipmark-tinted generate -pkg embedded -export -tinted-pkg github.com/pmarschik/lipmark/theme/tinted -o schemes_gen.go catppuccin-mocha dracula nord solarized-dark solarized-light

import "github.com/pmarschik/lipmark/theme"

// WithEmbedded returns an option that registers all generated themes
// (catppuccin-mocha, dracula, nord, solarized-dark, solarized-light)
// on a [theme.Registry]. Both the base styles and palette colors are registered.
func WithEmbedded() theme.Option {
	return func(r *theme.Registry) {
		for name, s := range Schemes {
			r.RegisterStruct(name, *s.Styles)
			r.RegisterStruct(name, *s.Palette)
		}
	}
}
