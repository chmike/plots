# Plot packages

The plot package provides functions to create lines or spike plots.

## Installation

To use this package execute `go get "github.com/chmike/plots"` the terminal after you executed the `go mod init` instruction.

## Line plots

A line plot is a conventional plot with lines, lines with points or simply points. It is possible to specify the line color, width, dashing, as well as the point glyph, color and size.

One may use the `plot.Add(p *plot.Plot, l plots.Line)`command to add a line to your plot, or use the `plots.MakeLinePlot(l plots.Lines, fileNames ...string)`command.

You may then generate of the following plots. The test code shows how to set the line properties.

![Line plot with custom size and different styles.](images/linePlot1.png)

![Line plot with default size and other styles.](images/linePlot2.png)

![Sine and cosine line plot.](images/linePlot3.png)

## Spike plots

Spike plots are used for Spiking Neural Networks (SNN) studies.

This package provides the function `MakeSpikePlot` to create such plots. It makes
it possible to also draw bars extending multiple spike lines to emphasis coincidences.

![Simple spike plot illustrating the possibilities.](images/simpleSpikePlot.png)

The following plot is a more realistic use case with Poisson distributed spikes and
20ms refractory period.

![More realistic spike plot.](images/poissonSpikePlot.png)

The following plot shows that one can easily select a sub-time range to display in
the plot.

![Sub-range spike plot.](images/poissonSpikeSubPlot.png)
