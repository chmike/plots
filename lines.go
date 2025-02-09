package plots

import (
	"fmt"
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// XYs converts the floats to XYs
func XYs(points []float64) plotter.XYs {
	pts := make(plotter.XYs, len(points))
	for i := range points {
		pts[i].X = float64(i)
		pts[i].Y = points[i]
	}
	return pts
}

// LineProperty are the line properties. To draw lines, specify its
// color, width or dashes. Default default values are used when needed.
// To draw points, specify its Glyph, color or radius. Default values are
// used when needed.
// type LineProperty struct {
// 	Width       vg.Length        // Line width.
// 	Color       color.Color      // Line color.
// 	Dashes      []vg.Length      // Line dashes.
// 	DashOffs    vg.Length        // Dashes offset.
// 	Glyph       draw.GlyphDrawer // Glyph to draw.
// 	GlyphColor  color.Color      // Glyph color.
// 	GlyphRadius vg.Length        // Glyph size.
// }

// Line is a line to be drawn with its properties. To draw lines, specify its
// color, width or dashes. Default default values are used when needed.
// To draw points, specify its Glyph, color or radius. Default values are
// used when needed.
type Line struct {
	Label       string           // Line label for the legend, not in legend if empty.
	Points      plotter.XYer     // Line points.
	Width       vg.Length        // Line width.
	Color       color.Color      // Line color.
	Dashes      []vg.Length      // Line dashes.
	DashOffs    vg.Length        // Dashes offset.
	Glyph       draw.GlyphDrawer // Glyph to draw.
	GlyphColor  color.Color      // Glyph color.
	GlyphRadius vg.Length        // Glyph size.
}

type Limit struct {
	Min, Max float64 // Axis limit values.
}

// Lines is a set of lines to be drawn.
type Lines struct {
	Title  string    // Line plot title.
	XLabel string    // X axis label, none if empty.
	YLabel string    // Y axis label, none if empty.
	XLimit *Limit    // XLimit values or nil if none.
	YLimit *Limit    // YLimit values or nil if none.
	Lines  []Line    // Lines to draw in plot.
	XDim   vg.Length // X dimension of saved plot, use default if 0.
	YDim   vg.Length // Y dimension of saved plot, use default if 0.
}

// MakeLinePlot generates the line plot.
func MakeLinePlot(lines Lines, fileNames ...string) error {
	p := plot.New()
	p.Title.Text = lines.Title
	p.X.Label.Text = lines.XLabel
	p.Y.Label.Text = lines.YLabel
	if lines.XLimit != nil {
		p.X.Min = lines.XLimit.Min
		p.X.Max = lines.XLimit.Max
	}
	if lines.YLimit != nil {
		p.Y.Min = lines.YLimit.Min
		p.Y.Max = lines.YLimit.Max
	}
	for i := range lines.Lines {
		err := Add(p, lines.Lines[i])
		if err != nil {
			return fmt.Errorf("line plot '%s': %w", lines.Lines[i].Label, err)
		}
	}
	xDim, yDim := lines.XDim, lines.YDim
	if xDim == 0 {
		xDim = 15 * vg.Centimeter
	}
	if yDim == 0 {
		yDim = 15 * vg.Centimeter
	}
	for _, fileName := range fileNames {
		err := p.Save(xDim, yDim, fileName)
		if err != nil {
			return fmt.Errorf("line plot: %w", err)
		}
	}
	return nil
}

// Add adds the points to the plot using the given style options.
func Add(plt *plot.Plot, line Line) error {
	var hasProperty bool
	if line.Color != nil || line.Width != 0 ||
		line.Dashes != nil || line.DashOffs != 0 {
		hasProperty = true
		if line.Color == nil {
			line.Color = rgb(20, 20, 20)
		}
		if line.Width == 0 {
			line.Width = vg.Points(1)
		}
		if line.DashOffs != 0 {
			if line.Dashes == nil {
				line.Dashes = Dashes.Id(0)
			}
		}
	}
	if line.Glyph != nil || line.GlyphColor != nil || line.GlyphRadius != 0 {
		hasProperty = true
		if line.Glyph == nil {
			line.Glyph = Glyphs.Id(0)
		}
		if line.GlyphColor == nil {
			line.GlyphColor = line.Color
			if line.GlyphColor == nil {
				line.GlyphColor = rgb(20, 20, 20)
			}
			if line.GlyphRadius == 0 {
				line.GlyphRadius = line.Width + vg.Points(1)
			}
		}
	}
	if !hasProperty {
		line.Color = rgb(20, 20, 20)
		line.Width = vg.Points(1)
	}

	var l *plotter.Line
	var s *plotter.Scatter
	xys, err := plotter.CopyXYs(line.Points)
	if err != nil {
		return err
	}
	if line.Width != 0 {
		l = &plotter.Line{
			XYs: xys,
			LineStyle: draw.LineStyle{
				Color:    line.Color,
				Width:    line.Width,
				Dashes:   line.Dashes,
				DashOffs: line.DashOffs,
			},
		}
	}
	if line.GlyphRadius != 0 {
		s = &plotter.Scatter{
			XYs: xys,
			GlyphStyle: draw.GlyphStyle{
				Shape:  line.Glyph,
				Color:  line.GlyphColor,
				Radius: line.GlyphRadius,
			},
		}
	}
	switch {
	case l != nil && s != nil:
		plt.Add(l, s)
	case l != nil && s == nil:
		plt.Add(l)
	case l == nil && s != nil:
		plt.Add(s)
	}
	if line.Label != "" {
		switch {
		case l != nil && s != nil:
			plt.Legend.Add(line.Label, l, s)
		case l != nil && s == nil:
			plt.Legend.Add(line.Label, l)
		case l == nil && s != nil:
			plt.Legend.Add(line.Label, s)
		}
	}
	return nil
}
