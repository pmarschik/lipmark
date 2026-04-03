# lipmark

Markup format strings, themed output, and glyph management for lipgloss CLI apps.

lipmark provides three composable systems that work together:

- **Markup** — `{name}` syntax for inline styled glyphs and text in format strings
- **Themes** — Named style presets (default, nocolor, tinted-theming) with extensible struct-tag registration
- **Glyphs** — Named glyph sets across display modes (unicode, ascii, nerd fonts, emoji)

All three adapt automatically to terminal capabilities (NO_COLOR, TERM, locale) so your output looks right everywhere — no `if tty` branches needed.

## Install

```bash
go get github.com/pmarschik/lipmark@latest
```

## Quick Start

```go
import "github.com/pmarschik/lipmark"

ui := lipmark.NewUI() // auto-detects theme + glyphs

ui.Success("deployed %s", "my-app")
ui.Error("build failed: %s", err)
ui.Stderr("{success} pushed %d commits {arrow} %s", 3, "origin/main")
ui.Stdout("result: {bold %s}", data)
```

## Markup Syntax

lipmark's format strings extend `fmt.Sprintf` with `{placeholder}` syntax:

| Syntax           | Meaning            | Example                             |
| ---------------- | ------------------ | ----------------------------------- |
| `{name}`         | Styled glyph       | `{success}` → styled "✓"            |
| `{name text}`    | Styled text        | `{bold hello}` → bold "hello"       |
| `%s`, `%d`, etc. | Standard fmt verbs | Resolved first, before placeholders |
| `{{` / `}}`      | Literal braces     | `{{name}}` → literal "{name}"       |

Format verbs are resolved first via `fmt.Sprintf`, then `{placeholders}` are resolved against the glyph and theme registries.

```go
ui.Sprintf("{success} repo {bold %s} updated {arrow} %s", name, branch)
// "✓ repo my-repo updated → main" (with styling)
```

The standalone `lipmark.Render()` function processes markup without a UI:

```go
out := lipmark.Render("Hello {bold world}!", glyphSet, themeSet)
```

## UI

The `UI` type wraps resolved glyphs + theme + writers for ergonomic CLI output.

```go
ui := lipmark.NewUI(
    lipmark.WithUIGlyphs(glyphSet),
    lipmark.WithUITheme(themeSet),
    lipmark.WithUIOut(os.Stdout),
    lipmark.WithUIErr(os.Stderr),
)
```

All options are optional — defaults to auto-detected theme, glyphs, stdout, and stderr.

### Output Methods

```go
// Markup-formatted output
ui.Stdout(format, args...)   // to stdout, with newline
ui.Stderr(format, args...)   // to stderr, with newline
ui.Print(format, args...)    // to stdout, no newline
ui.Sprintf(format, args...)  // returns string

// Semantic helpers (write to stderr with glyph prefix)
ui.Success(format, args...)  // {success} prefix
ui.Error(format, args...)    // {error} prefix, message styled red
ui.Warning(format, args...)  // {warning} prefix
ui.Info(format, args...)     // {info} prefix
ui.Note(format, args...)     // {arrow} prefix, message muted
ui.Status(format, args...)   // plain text, no markup processing

// Direct access
ui.Glyph("success")          // raw glyph character: "✓"
ui.StyledGlyph("success")    // styled glyph: "\033[32m✓\033[0m"
ui.Style("heading")          // lipgloss.Style for manual use
```

### Theme Preview

Render a sample showing all theme styles:

```go
lipmark.Preview(os.Stdout, glyphSet, themeSet)
```

## Themes

The `theme/` package is a registry of named style presets. Apps register their own style names alongside the built-in 15.

```go
import "github.com/pmarschik/lipmark/theme"

reg := theme.NewWithOptions(theme.WithDefaults()) // "default" + "nocolor" presets

// Resolve — "auto" picks "nocolor" on NO_COLOR/TERM=dumb, else "default"
styles := reg.Resolve("auto")
styles.Get(theme.Success).Render("ok")
styles.Get("issue_key").Render("PROJ-123") // app-specific style
```

### Extensible via Struct Tags

Define app styles as structs with `style:"name"` tags. Compose with embedding:

```go
type MyTheme struct {
    theme.BaseStyles                       // 15 built-in styles
    IssueKey lipgloss.Style `style:"issue_key"`
    Repo     lipgloss.Style `style:"repo"`
}

var t MyTheme
reg.ResolveInto("auto", &t)
t.Success.Render("ok")        // from BaseStyles
t.IssueKey.Render("PROJ-123") // app-specific
```

Register definitions the same way:

```go
type MyStyleDefs struct {
    theme.BaseStyleDefs
    IssueKey theme.StyleDef `style:"issue_key"`
}

reg.RegisterStruct("default", MyStyleDefs{
    IssueKey: theme.StyleDef{Foreground: "#2563EB", Bold: true},
})
```

### Built-in Style Names

`bold`, `dim`, `italic`, `underline`, `success`, `error`, `warning`, `info`, `muted`, `command`, `flag`, `heading`, `key`, `value`, `path`

### Loading from Files

```go
defs, _ := theme.LoadDefs(reader)         // from io.Reader (YAML)
defs, _ := theme.LoadDefsFile("theme.yaml")
merged := theme.MergeDefs(base, overlay)  // partial override
```

### Custom NO_COLOR

```go
theme.AutoPresetWithEnv([]string{"NO_COLOR", "MY_APP_NO_COLOR"})
theme.IsNoColor("NO_COLOR", "MY_APP_NO_COLOR")
```

## Glyphs

The `glyphs/` package is a registry of named glyph characters across display modes.

```go
import "github.com/pmarschik/lipmark/glyphs"

reg := glyphs.NewWithOptions(glyphs.WithDefaults()) // 34 glyphs x 4 modes
set := reg.Resolve("auto") // detects terminal capabilities
set.Get(glyphs.Success) // "✓" (unicode), "[ok]" (ascii), "" (nerd), "✅" (emoji)
```

### Extensible via Struct Tags

```go
type MyGlyphs struct {
    glyphs.BaseGlyphs               // 34 built-in glyphs
    Dirty    string `glyph:"dirty"`
    Unpushed string `glyph:"unpushed"`
}

reg.RegisterStruct("nerd", MyGlyphs{Dirty: "\uf044", Unpushed: "\uf062"})

var g MyGlyphs
reg.ResolveInto("auto", &g)
```

### Built-in Glyphs (34)

**Status:** success, error, warning, info, check, cross
**Navigation:** arrow, arrow_left, arrow_up, arrow_down, bullet, ellipsis, separator
**Actions:** star, heart, folder, file, lock, unlock, edit, trash, search, refresh, download, upload, clock, play, pause, stop, gear, link, user, plus, minus

### Modes

| Mode        | Detection                      | Example                 |
| ----------- | ------------------------------ | ----------------------- |
| **auto**    | Env vars, locale, TERM_PROGRAM | Auto-selects best mode  |
| **unicode** | Default for UTF-8 terminals    | ✓ ✗ ⚠ → ★               |
| **ascii**   | Non-UTF-8 or explicit          | [ok] [err] [warn] -> \* |
| **nerd**    | NERD_FONT env var              | Nerd Font icons         |
| **emoji**   | Modern terminal apps           | ✅ ❌ ⚠️ ➡️ ⭐            |

## Tinted Themes

The `theme/tinted/` package integrates with [tinted-theming](https://github.com/tinted-theming/schemes) — 300+ color schemes (base16 + base24).

### Embedded Schemes

5 popular schemes are pre-generated and available without any YAML dependency:

```go
import "github.com/pmarschik/lipmark/theme/tinted/embedded"

reg := theme.NewWithOptions(theme.WithDefaults(), embedded.WithEmbedded())
// Adds: catppuccin-mocha, dracula, nord, solarized-dark, solarized-light
```

Access palette colors directly:

```go
embedded.DraculaPalette.Base0C.Foreground // "#8be9fd" (cyan)
embedded.DraculaStyles.Error              // StyleDef{Foreground: dracula red, Bold: true}
```

### Typed Palette

Compose tinted palette with your app theme:

```go
import "github.com/pmarschik/lipmark/theme/tinted"

type MyTheme struct {
    theme.BaseStyles
    tinted.Palette  // base00-base17 as lipgloss.Style fields
    Custom lipgloss.Style `style:"custom"`
}

var t MyTheme
reg.ResolveInto("dracula", &t)
t.Base0C.Render("cyan text")  // from palette
t.Success.Render("ok")        // from base styles
```

### Loading Additional Schemes

For runtime YAML loading (adds `gopkg.in/yaml.v3` dependency):

```go
import "github.com/pmarschik/lipmark/theme/tinted/load"

defs, _ := load.File("tokyo-night.yaml")
reg.Register("tokyo-night", defs)

// Load an entire directory
schemes, _ := load.Dir("./schemes")
for name, defs := range schemes {
    reg.Register(name, defs)
}
```

### Code Generation

Generate Go source files from tinted-theming schemes — no YAML at runtime:

```go
//go:generate go run github.com/pmarschik/lipmark/cmd/lipmark-tinted generate -pkg mythemes -export -o themes_gen.go tokyo-night gruvbox-dark
```

The generator fetches schemes from GitHub (base24 first, falls back to base16), converts them using `DefaultMapping()`, and writes typed struct literals.

Flags: `-pkg` (package name), `-export`/`-no-export` (symbol visibility), `-tinted-pkg` (import path for tinted types), `-o` (output file).

To fetch raw YAML files for runtime loading:

```bash
go run github.com/pmarschik/lipmark/cmd/lipmark-tinted fetch tokyo-night -o themes/
```

### Palette Constants

Use `tinted.Base00` through `tinted.Base17` for typed palette access:

```go
styles.Get(tinted.Base0C) // cyan style from current theme
styles.Get(tinted.Base08) // red style
```

## Package Structure

```
lipmark/
  markup.go         # Render() — {placeholder} markup processor
  ui.go             # UI type — themed output with markup
  preview.go        # Preview() — theme sample rendering
  theme/
    theme.go        # Registry, BaseStyles, BaseStyleDefs, StyleDef
    tinted/
      tinted.go     # Palette, PaletteDefs, Mapping, constants
      embedded/     # 5 pre-generated schemes, WithEmbedded()
      load/         # YAML loading (separate dependency)
  glyphs/
    glyphs.go       # Registry, BaseGlyphs, 34 icons x 4 modes
  cmd/
    lipmark-tinted/ # Code generator + scheme fetcher
```

## Development

### Prerequisites

- [mise](https://mise.jdx.dev/) -- install tools and run tasks

### Setup

```bash
mise run setup
```

### Commands

```bash
mise run check    # all quality gates (format + lint + test)
mise run fmt      # format code
mise run lint     # lint
mise run test     # tests
go build ./...    # build
```

### Regenerating Embedded Schemes

```bash
go generate ./theme/tinted/embedded/
```

## Contributing

- Use [Conventional Commits](https://www.conventionalcommits.org/)
- Run `mise run check` before pushing
- See AGENTS.md for detailed conventions
