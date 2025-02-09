package plots

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"
	"os"
	"testing"

	"gonum.org/v1/plot/vg"
)

func TestSpikePlot(t *testing.T) {
	spikes := SpikeLines{
		Title: "Example with extended lines",
		Lines: []SpikeLine{
			{
				Label:  "line 0",
				Spikes: []float64{0.5, 2, 2.7, 3, 5.8, 6.9, 8},
			},
			{
				Label:  "line 1",
				Spikes: []float64{0.1, .4, 1.3, 1.2, 2.4, 2.1, 3.5, 4.1, 5.8, 6.9},
			},
			{
				Label:  "line 2",
				Spikes: []float64{0, .5, 1., 1.5, 2, 2.5, 3., 4, 5, 7},
				Property: SpikeLineProperty{
					Color: color.RGBA{50, 50, 50, 255},
				},
			},
			{
				Label:  "line 3",
				Spikes: []float64{.2, .7, 1.2, 1.7, 2.2, 2.7, 3.2, 4.2, 5.2, 7.2},
				ZIndex: -1,
				Property: SpikeLineProperty{
					Extend:      2,
					ExtendColor: color.RGBA{255, 150, 150, 255},
					ExtendWidth: vg.Points(3),
					Color:       color.RGBA{0, 0, 255, 255},
				},
			},
			{
				Label:  "line 4",
				Spikes: []float64{},
			},
			{
				Label:  "line 5",
				Spikes: []float64{0.5, 2, 2.7, 3, 5.8, 6.9, 8},
			},
			{
				Label:  "line 6",
				Spikes: nil,
			},
		},
	}
	os.MkdirAll("tests", 0766)
	err := MakeSpikePlot(spikes, "tests/simpleSpikePlot.png", "tests/simpleSpikePlot.svg")
	if err != nil {
		t.Fatalf("failed saving image: %s", err)
	}
}

var rng = rand.New(rand.NewPCG(1, 2))

// GeneratePoissonDistributedSpikes returns a sequence of poisson distributed spikes
// with the given rate and in the given duration time span.
func GeneratePoissonDistributedSpikes(duration, rate, refractory float64) []float64 {
	if rate <= 0 || duration <= 0 {
		return []float64{}
	}
	expectedEvents := rate * duration
	stdDev := math.Sqrt(expectedEvents)
	capacity := int(expectedEvents + 3*stdDev)
	spikes := make([]float64, 0, capacity)
	var t float64
	for t < duration {
		interval := rng.ExpFloat64() / rate
		for interval <= refractory {
			interval = rng.ExpFloat64() / rate
		}
		t += interval
		spikes = append(spikes, t)
	}
	return spikes[:len(spikes)-1]
}

// CopySpikes into a new spike array.
func CopySpikes(spikes []float64) []float64 {
	return append([]float64(nil), spikes...)
}

// AddPoissonNoise adds poisson noise to the given spikes, ensuring that the
// spike interval is not smaller than refractory.
func AddPoissonNoise(spikes []float64, duration, rate, refractory float64) []float64 {
	if rate <= 0 || duration <= 0 {
		return CopySpikes(spikes)
	}
	expectedEvents := rate * duration
	stdDev := math.Sqrt(expectedEvents)
	capacity := int(expectedEvents + 3*stdDev)
	res := make([]float64, 0, capacity+len(spikes))
	var nextSpike int
	var t float64
	for t < duration && nextSpike < len(spikes) {
		// get new spike interval bigger than refractory
		interval := rng.ExpFloat64() / rate
		for interval <= refractory {
			interval = rng.ExpFloat64() / rate
		}
		// if it fits before the next spike, accept it
		if t+interval+refractory < spikes[nextSpike] {
			t += interval
			res = append(res, t)
			continue
		}
		// output the spike
		t = spikes[nextSpike]
		res = append(res, t)
		nextSpike++
	}
	for t < duration {
		// get new spike interval bigger than refractory
		interval := rng.ExpFloat64() / rate
		for interval <= refractory {
			interval = rng.ExpFloat64() / rate
		}
		t += interval
		res = append(res, t)
	}
	for len(res) > 0 && res[len(res)-1] > duration {
		res = res[:len(res)-1]
	}
	return res
}

func TestManySpikeLines(t *testing.T) {

	os.MkdirAll("tests", 0766)
	duration := 4.     // seconds
	rate := 8.         // default noise rate (Hz)
	commonRate := 3.   // common spike rate (Hz)
	refractory := 0.02 // refractory period is 20ms
	nLines := 30       // number of lines to generated
	ratio := .5        // percentage of random spikes

	// generation the common spikes.
	commonSpikes := GeneratePoissonDistributedSpikes(duration, commonRate, refractory)

	var spikes SpikeLines
	for i := 0; i < nLines; i++ {
		if rng.Float64() < ratio {
			spikes.Lines = append(spikes.Lines, SpikeLine{
				Label:  fmt.Sprintf("%d", i+1),
				Spikes: GeneratePoissonDistributedSpikes(duration, rate, refractory),
			})
		} else {
			spikes.Lines = append(spikes.Lines, SpikeLine{
				Label:  fmt.Sprintf("%d", i+1),
				Spikes: AddPoissonNoise(commonSpikes, duration, rate-commonRate, refractory),
			})
		}
	}
	spikes.Lines = append(spikes.Lines, SpikeLine{
		Label:  "output",
		Spikes: commonSpikes,
		ZIndex: -1,
		Property: SpikeLineProperty{
			Extend:      nLines,
			ExtendWidth: vg.Points(3),
			ExtendColor: color.RGBA{200, 200, 200, 255},
		},
	})
	spikes.Title = "Poisson distributed spikes"

	err := MakeSpikePlot(spikes, "tests/poissonSpikePlot.png", "tests/poissonSpikePlot.svg")
	if err != nil {
		t.Fatalf("failed saving image: %s", err)
	}
	spikes.XMin = 1
	spikes.XMax = duration / 2

	err = MakeSpikePlot(spikes, "tests/poissonSpikeSubPlot.png", "tests/poissonSpikeSubPlot.svg")
	if err != nil {
		t.Fatalf("failed saving image: %s", err)
	}
}
