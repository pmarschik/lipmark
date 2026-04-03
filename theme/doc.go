// Package theme provides an extensible style registry for CLI applications.
//
// A [Registry] holds style definitions across named presets (e.g. "default",
// "nocolor", "dracula"). Apps register their own style names and per-preset
// definitions on top of the built-in defaults.
//
// # Basic usage with map registration
//
// Use [WithDefaults] to pre-populate the standard styles (bold, success,
// error, etc.) across "default" and "nocolor" presets, then add app-specific
// styles via [Registry.Register]:
//
//	reg := theme.NewWithOptions(theme.WithDefaults())
//	reg.Register("default", map[string]theme.StyleDef{
//	    "issue_key": {Foreground: "#2563EB", Bold: true},
//	})
//	reg.Register("nocolor", map[string]theme.StyleDef{
//	    "issue_key": {},
//	})
//
//	set := reg.Resolve("auto") // detects NO_COLOR/TERM=dumb
//	fmt.Println(set.Get(theme.Success).Render("ok"))
//
// # Struct-based registration
//
// For type-safe style access, define a struct with `style` tags and use
// [Registry.RegisterStruct] and [Registry.ResolveInto]:
//
//	type AppStyles struct {
//	    IssueKey theme.StyleDef `style:"issue_key"`
//	    Repo     theme.StyleDef `style:"repo"`
//	}
//
//	reg := theme.NewWithOptions(theme.WithDefaults())
//	reg.RegisterStruct("default", AppStyles{
//	    IssueKey: theme.StyleDef{Foreground: "#2563EB", Bold: true},
//	    Repo:     theme.StyleDef{Bold: true},
//	})
//
//	type Resolved struct {
//	    IssueKey lipgloss.Style `style:"issue_key"`
//	    Bold     lipgloss.Style `style:"bold"` // base style
//	}
//	var s Resolved
//	_ = reg.ResolveInto("auto", &s)
//
// # Preset detection
//
// [AutoPreset] inspects the NO_COLOR and TERM environment variables to pick
// "nocolor" or "default". Pass "auto" to [Registry.Resolve] or
// [Registry.ResolveInto] to use it.
//
// # Loading from files
//
// [LoadDefs] and [LoadDefsFile] parse YAML files into style definition maps
// that can be registered on a preset. [MergeDefs] overlays partial
// definitions onto a base map.
package theme
