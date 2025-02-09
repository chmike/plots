package plot

import (
	"math"
	"os"
	"testing"

	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func TestLinePlot1(t *testing.T) {
	os.MkdirAll("tests", 0766)
	lines := Lines{
		Title:  "lines example 1",
		XLabel: "X",
		YLabel: "Y",
		XDim:   15 * vg.Centimeter,
		YDim:   9 * vg.Centimeter,
		Lines: []Line{
			{
				Label:  "line 1",
				Points: XYs([]float64{.2, .7, 1.2, 1.7, 2.2, 2.7, 3.2, 4.2, 5.2, 7.2}),
				Color:  DarkColors.Id(1),
				Glyph:  Glyphs.Id(1),
			},
			{
				Label:  "line 2",
				Points: XYs([]float64{0.5, 2, 2.7, 3, 5.8, 6.9, 8}),
				Color:  DarkColors.Id(3),
				Glyph:  Glyphs.Id(0),
			},
			{
				Label:  "line 3",
				Points: XYs([]float64{0.1, .2, .7, 2, 3.8, 4.9, 5, 7}),
				Color:  DarkColors.Id(2),
				Dashes: Dashes.Id(1),
				Glyph:  Glyphs.Id(2),
			},
			{
				Label:       "points",
				Points:      XYs([]float64{1, 3, 3.8, 4.6, 5.4, 6.2, 7, 7.8}),
				Glyph:       Glyphs.Id(5),
				GlyphRadius: vg.Points(2),
			},
		},
	}
	err := MakeLinePlot(lines, "tests/linePlot1.png", "tests/linePlot1.svg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLinePlot2(t *testing.T) {
	os.MkdirAll("tests", 0766)
	lines := Lines{
		Title: "lines example 2",
		Lines: []Line{
			{
				Points: XYs([]float64{.2, .7, 1.2, 1.7, 2.2, 2.7, 3.2, 4.2, 5.2, 7.2}),
			},
			{
				Label:  "line 2",
				Points: XYs([]float64{0.5, 2, 2.7, 3, 5.8, 6.9, 8}),
				Color:  DarkColors.Id(1),
			},
			{
				Label:  "line 3",
				Points: XYs([]float64{0.1, .2, .7, 2, 3.8, 4.9, 5, 7}),
				Dashes: Dashes.Id(0),
			},
			{
				Label:       "points",
				Points:      XYs([]float64{1, 3, 3.8, 4.6, 5.4, 6.2, 7, 7.8}),
				GlyphRadius: vg.Points(2),
			},
		},
	}
	err := MakeLinePlot(lines, "tests/linePlot2.png", "tests/linePlot2.svg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLinePlot3(t *testing.T) {
	os.MkdirAll("tests", 0666)
	var cosXYs, sinXYs plotter.XYs
	for x := -math.Pi; x < math.Pi; x += 2 * math.Pi / 256. {
		cosXYs = append(cosXYs, plotter.XY{X: x, Y: math.Cos(x)})
		sinXYs = append(sinXYs, plotter.XY{X: x, Y: math.Sin(x)})
	}
	lines := Lines{
		Title: "lines example 3",
		Lines: []Line{
			{
				Points: plotter.XYs{
					{X: -math.Pi, Y: 0},
					{X: math.Pi, Y: 0},
				},
				Dashes: Dashes.Id(0),
			},
			{
				Label:  "Sin",
				Points: sinXYs,
				Color:  DarkColors.Id(1),
			},
			{
				Label:  "Cos",
				Points: cosXYs,
				Color:  DarkColors.Id(2),
			},
		},
	}
	err := MakeLinePlot(lines, "tests/linePlot3.png", "tests/linePlot3.svg")
	if err != nil {
		t.Fatal(err)
	}
}
