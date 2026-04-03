// Package glyphs provides an extensible glyph registry for CLI applications.
//
// A [Registry] holds glyph character mappings across multiple display modes
// (unicode, ascii, nerd fonts, emoji). Apps register their own glyph names
// and per-mode characters on top of the built-in defaults.
//
// # Basic usage with map registration
//
// Use [WithDefaults] to pre-populate the standard glyphs (success, error,
// warning, etc.) across all four built-in modes, then add app-specific
// glyphs via [Registry.Register]:
//
//	reg := glyphs.NewWithOptions(glyphs.WithDefaults())
//	reg.Register("nerd", map[string]string{"dirty": "", "unpushed": ""})
//	reg.Register("emoji", map[string]string{"dirty": "🟠", "unpushed": "⬆"})
//	reg.Register("ascii", map[string]string{"dirty": "M", "unpushed": "^"})
//	reg.Register("unicode", map[string]string{"dirty": "M", "unpushed": "⇡"})
//
//	set := reg.Resolve("auto") // detects terminal capabilities
//	fmt.Println(set.Get("dirty"))
//
// # Struct-based registration
//
// For type-safe glyph access, define a struct with `glyph` tags and use
// [Registry.RegisterStruct] and [Registry.ResolveInto]:
//
//	type AppGlyphs struct {
//	    Success  string `glyph:"success"`
//	    Error    string `glyph:"error"`
//	    Dirty    string `glyph:"dirty"`
//	    Unpushed string `glyph:"unpushed"`
//	}
//
//	reg := glyphs.NewWithOptions(glyphs.WithDefaults())
//	reg.RegisterStruct("nerd", AppGlyphs{Dirty: "", Unpushed: ""})
//	reg.RegisterStruct("unicode", AppGlyphs{Dirty: "M", Unpushed: "⇡"})
//
//	var g AppGlyphs
//	_ = reg.ResolveInto("auto", &g)
//	fmt.Println(g.Dirty) // "" on nerd font terminals, "M" on unicode
//
// [RegisterStruct] only registers non-empty fields, so base glyphs from
// [WithDefaults] are preserved. [ResolveInto] leaves struct fields unchanged
// if no matching glyph is found in the resolved mode.
//
// # Mode detection
//
// [AutoMode] inspects environment variables (NERD_FONT, POWERLINE_*,
// TERM_PROGRAM, locale) to pick the best mode for the current terminal.
// Pass "auto" to [Registry.Resolve] or [Registry.ResolveInto] to use it.
package glyphs
