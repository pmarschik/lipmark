package glyphs_test

import (
	"fmt"

	"github.com/pmarschik/lipmark/glyphs"
)

func Example() {
	reg := glyphs.NewWithOptions(glyphs.WithDefaults())
	set := reg.Resolve("ascii")

	fmt.Println(set.Get(glyphs.Success))
	fmt.Println(set.Get(glyphs.Error))
	fmt.Println(set.Get(glyphs.Arrow))
	// Output:
	// [ok]
	// [err]
	// ->
}

func Example_appSpecific() {
	reg := glyphs.NewWithOptions(glyphs.WithDefaults())
	reg.Register("ascii", map[string]string{
		"dirty":    "M",
		"unpushed": "^",
	})
	reg.Register("unicode", map[string]string{
		"dirty":    "M",
		"unpushed": "⇡",
	})

	set := reg.Resolve("ascii")
	fmt.Println(set.Get("dirty"))
	fmt.Println(set.Get("unpushed"))
	fmt.Println(set.Get(glyphs.Success)) // base glyphs still available
	// Output:
	// M
	// ^
	// [ok]
}

// AppGlyphs demonstrates a typed glyph struct with glyph tags.
type AppGlyphs struct {
	Success  string `glyph:"success"`
	Dirty    string `glyph:"dirty"`
	Unpushed string `glyph:"unpushed"`
}

func Example_structBased() {
	reg := glyphs.NewWithOptions(glyphs.WithDefaults())
	reg.RegisterStruct("ascii", AppGlyphs{Dirty: "M", Unpushed: "^"})
	reg.RegisterStruct("unicode", AppGlyphs{Dirty: "M", Unpushed: "⇡"})

	var g AppGlyphs
	if err := reg.ResolveInto("ascii", &g); err != nil {
		panic(err)
	}

	fmt.Println(g.Success)
	fmt.Println(g.Dirty)
	fmt.Println(g.Unpushed)
	// Output:
	// [ok]
	// M
	// ^
}
