package plots

import (
	"image/color"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// ColorTable is a table of predefined colors to pick from.
type ColorTable []color.Color

// Id returns color i from color table.
func (t ColorTable) Id(i int) color.Color {
	n := len(t)
	if i < 0 {
		return t[i%n+n]
	}
	return t[i%n]
}

var DarkColors = ColorTable{
	rgb(32, 32, 32),
	rgb(238, 46, 47),
	rgb(24, 90, 169),
	rgb(0, 140, 72),
	rgb(244, 125, 35),
	rgb(102, 44, 145),
	rgb(162, 29, 33),
	rgb(180, 56, 148),
}

var SoftColors = ColorTable{
	rgb(100, 100, 100),
	rgb(241, 90, 96),
	rgb(90, 155, 212),
	rgb(122, 195, 106),
	rgb(250, 167, 91),
	rgb(158, 103, 171),
	rgb(206, 112, 88),
	rgb(215, 127, 180),
}

func rgb(r, g, b uint8) color.RGBA {
	return color.RGBA{r, g, b, 255}
}

// GlyphTable is a table of glyphs one can pick from.
type GlyphTable []draw.GlyphDrawer

// Id returns the glyph i from the GlyphTable.
func (t GlyphTable) Id(i int) draw.GlyphDrawer {
	n := len(t)
	if i < 0 {
		return t[i%n+n]
	}
	return t[i%n]
}

// Glyphs is a table of glyphs on can pick from.
var Glyphs = GlyphTable{
	draw.CircleGlyph{},
	draw.BoxGlyph{},
	draw.PyramidGlyph{},
	draw.RingGlyph{},
	draw.SquareGlyph{},
	draw.TriangleGlyph{},
	draw.CrossGlyph{},
	draw.PlusGlyph{},
}

// DashesTable is a table of dashes one can pic from.
type DashesTable [][]vg.Length

func (t DashesTable) Id(i int) []vg.Length {
	n := len(t)
	if i < 0 {
		return t[i%n+n]
	}
	return t[i%n]
}

// DefaultDashes is a set of dash patterns used by
// the Dashes function.
var Dashes = DashesTable{
	//	{},

	{vg.Points(6), vg.Points(2)},

	{vg.Points(2), vg.Points(2)},

	{vg.Points(1), vg.Points(1)},

	{vg.Points(5), vg.Points(2), vg.Points(1), vg.Points(2)},

	{vg.Points(10), vg.Points(2), vg.Points(2), vg.Points(2),
		vg.Points(2), vg.Points(2), vg.Points(2), vg.Points(2)},

	{vg.Points(10), vg.Points(2), vg.Points(2), vg.Points(2)},

	{vg.Points(5), vg.Points(2), vg.Points(5), vg.Points(2),
		vg.Points(2), vg.Points(2), vg.Points(2), vg.Points(2)},

	{vg.Points(4), vg.Points(2), vg.Points(4), vg.Points(1),
		vg.Points(1), vg.Points(1), vg.Points(1), vg.Points(1),
		vg.Points(1), vg.Points(1)},
}
