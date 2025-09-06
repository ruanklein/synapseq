package audio

import (
	"fmt"
	"math"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	t "github.com/ruanklein/synapseq/internal/types"
)

// AudioRenderer handle audio generation
type AudioRenderer struct {
	channels       [t.NumberOfChannels]t.Channel
	periods        []t.Period
	periodIdx      int      // Current period index
	waveTables     [4][]int // Waveform tables for different waveforms
	noiseGenerator *NoiseGenerator
	sampleRate     int
	volume         int // Volume level (0-100)
}

// NewAudioRenderer creates a new AudioRenderer instance
func NewAudioRenderer(periods []t.Period, option *t.Option) (*AudioRenderer, error) {
	if option == nil {
		return nil, fmt.Errorf("audio renderer options are required")
	}

	renderer := &AudioRenderer{
		periods:        periods,
		waveTables:     InitWaveformTables(),
		noiseGenerator: NewNoiseGenerator(),
		sampleRate:     option.SampleRate,
		volume:         option.Volume,
		periodIdx:      0,
	}

	return renderer, nil
}

// RenderToWAV renders the audio to a WAV file using go-audio/wav
func (r *AudioRenderer) RenderToWAV(outPath string) error {
	out := os.Stdout // Use standard output as default

	if outPath != "-" {
		var err error
		out, err = os.Create(outPath)
		if err != nil {
			return fmt.Errorf("create output: %w", err)
		}
		defer out.Close()
	}

	enc := wav.NewEncoder(out, r.sampleRate, 16, 2, 1)
	defer enc.Close()

	endMs := r.periods[len(r.periods)-1].Time

	totalFrames := int64(math.Round(float64(endMs) * float64(r.sampleRate) / 1000.0))
	chunkFrames := int64(t.BufferSize)
	framesWritten := int64(0)

	statusReporter := NewStatusReporter(false)
	defer statusReporter.FinalStatus()

	samples := make([]int, t.BufferSize*2) // Stereo: left + right
	audioBuf := &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 2,
			SampleRate:  r.sampleRate,
		},
		Data:           samples,
		SourceBitDepth: 16,
	}

	for framesWritten < totalFrames {
		currentTimeMs := int((float64(framesWritten) * 1000.0) / float64(r.sampleRate))
		r.sync(currentTimeMs)

		statusReporter.CheckPeriodChange(r)

		audioBuf.Data = r.mix(samples)

		framesToWrite := chunkFrames
		if remain := totalFrames - framesWritten; remain < chunkFrames {
			framesToWrite = remain
			audioBuf.Data = audioBuf.Data[:remain*2] // stereo interleaved
		}

		if err := enc.Write(audioBuf); err != nil {
			enc.Close()
			return fmt.Errorf("write wav: %w", err)
		}

		framesWritten += framesToWrite

		if statusReporter.ShouldUpdateStatus() {
			statusReporter.DisplayStatus(r, currentTimeMs)
		}
	}

	return nil
}
