package plots

import (
	"fmt"
	"image/color"
	"math"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// SpikeLines is a plot of spikes lines drawn from top to bottom.
type SpikeLines struct {
	Title  string      // Title
	Lines  []SpikeLine // Spike lines.
	XLimit *Limit      // Spike time range limit.
	XDim   vg.Length   // X dimension of saved plot, use default if 0.
	YDim   vg.Length   // Y dimension of saved plot, use default if 0.

}

// SpikeLine is a labeled sequence of spike event time values.
// Requires that spike event times are sorted in increasing order.
type SpikeLine struct {
	Label    string            // Spike sequence label.
	Spikes   []float64         // Spike event time values.
	ZIndex   int               // Drawing order in increasing value order.
	Property SpikeLineProperty // Spike line property (use default if nil).
}

// SpikeLineProperty is the the graphical property of spikes.
type SpikeLineProperty struct {
	Width       vg.Length   // Spike stroke width (default = vg.Points(1)).
	Color       color.Color // Spike color (default = black).
	LWidth      vg.Length   // Horizontal line thickness (default = vg.Points(1)).
	LColor      color.Color // Horizontal line color (default = black).
	Extend      int         // Extend spike over n lines above (default = 0).
	ExtendColor color.Color // Extend spike color (default = grey).
	ExtendWidth vg.Length   // Extend spike width (default = vg.Points(1)).
}

// DefaultSpikeProperty returns a default spike line property.
func DefaultSpikeProperty() SpikeLineProperty {
	return SpikeLineProperty{
		Width:       vg.Points(1),
		Color:       color.RGBA{0, 0, 0, 255},
		LWidth:      vg.Points(1),
		LColor:      color.RGBA{0, 0, 0, 255},
		Extend:      0,
		ExtendColor: color.RGBA{160, 160, 160, 255},
		ExtendWidth: vg.Points(1),
	}
}

// Plot draws the spike lines.
func (s SpikeLines) Plot(canvas draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&canvas)

	// determine drawing order
	drawingOrder := make([]int, len(s.Lines))
	for i := range drawingOrder {
		drawingOrder[i] = i
	}
	sort.Slice(drawingOrder, func(i, j int) bool {
		return s.Lines[drawingOrder[i]].ZIndex < s.Lines[drawingOrder[j]].ZIndex
	})

	// horizontal line coordinates
	xMinPx := trX(plt.X.Min)
	xMaxPx := trX(plt.X.Max)

	// spike line height
	dy := (plt.Y.Max - plt.Y.Min) / float64(len(s.Lines))

	// draw lines
	for _, i := range drawingOrder {
		l := &s.Lines[i]

		// get property by using default values
		spikeProperty := DefaultSpikeProperty()
		if l.Property.Extend != 0 {
			spikeProperty.Extend = l.Property.Extend
			if l.Property.ExtendWidth != 0 {
				spikeProperty.ExtendWidth = l.Property.ExtendWidth
			}
			if l.Property.ExtendColor != nil {
				spikeProperty.ExtendColor = l.Property.ExtendColor
			}
		}
		if l.Property.Width != 0 {
			spikeProperty.Width = l.Property.Width
		}
		if l.Property.Color != nil {
			spikeProperty.Color = l.Property.Color
		}
		if l.Property.LWidth != 0 {
			spikeProperty.LWidth = l.Property.LWidth
		}
		if l.Property.LColor != nil {
			spikeProperty.LColor = l.Property.LColor
		}

		// find spikes to draw
		beg := sort.SearchFloat64s(l.Spikes, plt.X.Min)
		end := sort.SearchFloat64s(l.Spikes, plt.X.Max)
		spikes := l.Spikes[beg:end]

		// compute horizontal line y coordinate from top to bottom
		yMin := plt.Y.Min + dy*float64(len(s.Lines)-i-1)
		yMinPx := trY(yMin)

		// draw extend spikes if any
		if spikeProperty.Extend != 0 {
			yMax := yMin + dy*float64(spikeProperty.Extend+1) - dy/4
			yMaxPx := trY(yMax)
			canvas.SetLineStyle(draw.LineStyle{
				Color: spikeProperty.ExtendColor,
				Width: spikeProperty.ExtendWidth,
			})
			for _, v := range spikes {
				var path vg.Path
				xPixel := trX(v)
				path.Move(vg.Point{X: xPixel, Y: yMinPx})
				path.Line(vg.Point{X: xPixel, Y: yMaxPx})
				canvas.Stroke(path)
			}
		}

		// draw spikes
		yMax := yMin + dy/2
		yMaxPx := trY(yMax)
		canvas.SetLineStyle(draw.LineStyle{
			Color: spikeProperty.Color,
			Width: spikeProperty.Width,
		})
		for _, v := range spikes {
			var path vg.Path
			xPixel := trX(v)
			path.Move(vg.Point{X: xPixel, Y: yMinPx})
			path.Line(vg.Point{X: xPixel, Y: yMaxPx})
			canvas.Stroke(path)
		}

		// draw horizontal line
		canvas.SetLineStyle(draw.LineStyle{
			Color: spikeProperty.LColor,
			Width: spikeProperty.LWidth,
		})
		var path vg.Path
		path.Move(vg.Point{X: xMinPx, Y: yMinPx})
		path.Line(vg.Point{X: xMaxPx, Y: yMinPx})
		canvas.Stroke(path)
	}
}

// Ticks generates the Y axis ticks with the line labels.
func (s SpikeLines) Ticks(min, max float64) []plot.Tick {
	dy := (max - min) / float64(len(s.Lines))
	ticks := make([]plot.Tick, 0, len(s.Lines))
	for i := range s.Lines {
		ticks = append(ticks, plot.Tick{
			Value: min + dy*float64(len(s.Lines)-i-1),
			Label: s.Lines[i].Label,
		})
	}
	return ticks
}

// MakeSpikePlot generates the spike plot for a given time range.
func MakeSpikePlot(spikeLines SpikeLines, fileNames ...string) error {
	if len(fileNames) == 0 {
		return nil
	}

	p := plot.New()
	p.Title.Text = spikeLines.Title
	p.X.Label.Text = "Time (s)"
	if spikeLines.XLimit == nil {
		xMin, xMax := math.Inf(1), math.Inf(-1)
		for i := range spikeLines.Lines {
			for _, v := range spikeLines.Lines[i].Spikes {
				if v < xMin {
					xMin = v
				}
				if v > xMax {
					xMax = v
				}
			}
		}
		p.X.Min = xMin
		p.X.Max = xMax
	} else {
		p.X.Min = spikeLines.XLimit.Min
		p.X.Max = spikeLines.XLimit.Max
	}
	p.Y.Tick.Marker = spikeLines
	p.Add(spikeLines)
	xDim, yDim := spikeLines.XDim, spikeLines.YDim
	if xDim == 0 {
		xDim = 15 * vg.Centimeter
	}
	if yDim == 0 {
		yDim = 15 * vg.Centimeter
	}
	for _, fileName := range fileNames {
		err := p.Save(xDim, yDim, fileName)
		if err != nil {
			return fmt.Errorf("spike plot: %w", err)
		}
	}
	return nil
}
