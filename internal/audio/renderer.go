/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package audio

import (
	"fmt"
	"math"

	t "github.com/ruanklein/synapseq/internal/types"
)

const (
	audioChannels = 2        // Stereo
	audioBitDepth = 24       // 24-bit audio
	audioBitShift = 8        // 24 Bit shift
	audioMaxValue = 8388607  // 2^23 - 1
	audioMinValue = -8388608 // -2^23
)

// AudioRenderer handle audio generation
type AudioRenderer struct {
	channels        [t.NumberOfChannels]t.Channel
	periods         []t.Period
	waveTables      [4][]int
	noiseGenerator  *NoiseGenerator
	backgroundAudio *BackgroundAudio

	dither0 uint16
	dither1 uint16

	// Embedding options
	*AudioRendererOptions
}

// AudioRendererOptions holds options for the audio renderer
type AudioRendererOptions struct {
	SampleRate     int
	Volume         int
	GainLevel      t.GainLevel
	BackgroundPath string
	Quiet          bool
	Debug          bool
}

// NewAudioRenderer creates a new AudioRenderer instance
func NewAudioRenderer(p []t.Period, ar *AudioRendererOptions) (*AudioRenderer, error) {
	if ar == nil {
		return nil, fmt.Errorf("audio renderer options cannot be nil")
	}

	if ar.SampleRate <= 0 {
		return nil, fmt.Errorf("invalid sample rate: %d", ar.SampleRate)
	}

	if ar.Volume < 0 || ar.Volume > 100 {
		return nil, fmt.Errorf("volume must be between 0 and 100, got %d", ar.Volume)
	}

	if len(p) == 0 {
		return nil, fmt.Errorf("no periods defined in the sequence")
	}

	// Initialize background audio
	backgroundAudio, err := NewBackgroundAudio(ar.BackgroundPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize background audio: %w", err)
	}

	// Validate background audio parameters
	if backgroundAudio.isEnabled {
		bgSampleRate := backgroundAudio.sampleRate
		if bgSampleRate != ar.SampleRate {
			return nil, fmt.Errorf("background audio sample rate (%d Hz) does not match output sample rate (%d Hz)",
				bgSampleRate, ar.SampleRate)
		}
		bgChannels := backgroundAudio.channels
		if bgChannels != audioChannels {
			return nil, fmt.Errorf("background audio must be stereo (%d channels detected)", bgChannels)
		}
		bgBitDepth := backgroundAudio.bitDepth
		if bgBitDepth != audioBitDepth {
			return nil, fmt.Errorf("background audio must be %d-bit (detected %d-bit)", audioBitDepth, bgBitDepth)
		}
	}

	renderer := &AudioRenderer{
		periods:              p,
		waveTables:           InitWaveformTables(),
		noiseGenerator:       NewNoiseGenerator(),
		backgroundAudio:      backgroundAudio,
		dither0:              1,
		dither1:              0,
		AudioRendererOptions: ar,
	}

	return renderer, nil
}

// Render generates the audio and passes buffers to the consume function
func (r *AudioRenderer) Render(consume func(samples []int) error) error {
	// Ensure background audio file is closed if opened
	defer func() {
		if r.backgroundAudio != nil {
			r.backgroundAudio.Close()
		}
	}()

	endMs := r.periods[len(r.periods)-1].Time
	totalFrames := int64(math.Round(float64(endMs) * float64(r.SampleRate) / 1000.0))
	chunkFrames := int64(t.BufferSize)
	framesWritten := int64(0)

	statusReporter := NewStatusReporter(r.Quiet && !r.Debug)
	defer statusReporter.FinalStatus()

	// Stereo: left + right
	samples := make([]int, t.BufferSize*audioChannels)
	periodIdx := 0

	for framesWritten < totalFrames {
		currentTimeMs := int((float64(framesWritten) * 1000.0) / float64(r.SampleRate))
		// Find the correct period for the current time
		for periodIdx+1 < len(r.periods) && currentTimeMs >= r.periods[periodIdx+1].Time {
			periodIdx++
		}

		r.sync(currentTimeMs, periodIdx)
		statusReporter.CheckPeriodChange(r, periodIdx)

		data := r.mix(samples)

		framesToWrite := chunkFrames
		if remain := totalFrames - framesWritten; remain < chunkFrames {
			framesToWrite = remain
			// stereo interleaved
			data = data[:remain*audioChannels]
		}

		if consume != nil {
			if err := consume(data); err != nil {
				return fmt.Errorf("failed to consume audio buffer: %w", err)
			}
		}

		framesWritten += framesToWrite

		if statusReporter.ShouldUpdateStatus() {
			statusReporter.DisplayStatus(r, currentTimeMs)
		}
	}

	return nil
}
