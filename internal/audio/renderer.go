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
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
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
		AudioRendererOptions: ar,
	}

	return renderer, nil
}

// RenderToWAV renders the audio to a WAV file using go-audio/wav
func (r *AudioRenderer) RenderToWAV(outPath string) error {
	var out *os.File
	var enc *wav.Encoder

	if !r.Debug {
		var err error
		out, err = os.Create(outPath)
		if err != nil {
			return fmt.Errorf("create output: %w", err)
		}
		defer out.Close()

		enc = wav.NewEncoder(out, r.SampleRate, audioBitDepth, audioChannels, 1)
		defer enc.Close()
	}

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

	samples := make([]int, t.BufferSize*audioChannels) // Stereo: left + right
	audioBuf := &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: audioChannels,
			SampleRate:  r.SampleRate,
		},
		Data:           samples,
		SourceBitDepth: audioBitDepth,
	}

	periodIdx := 0
	for framesWritten < totalFrames {
		currentTimeMs := int((float64(framesWritten) * 1000.0) / float64(r.SampleRate))
		// Find the correct period for the current time
		for periodIdx+1 < len(r.periods) && currentTimeMs >= r.periods[periodIdx+1].Time {
			periodIdx++
		}

		r.sync(currentTimeMs, periodIdx)
		statusReporter.CheckPeriodChange(r, periodIdx)

		audioBuf.Data = r.mix(samples)

		framesToWrite := chunkFrames
		if remain := totalFrames - framesWritten; remain < chunkFrames {
			framesToWrite = remain
			audioBuf.Data = audioBuf.Data[:remain*audioChannels] // stereo interleaved
		}

		if !r.Debug {
			if err := enc.Write(audioBuf); err != nil {
				enc.Close()
				return fmt.Errorf("write wav: %w", err)
			}
		}

		framesWritten += framesToWrite

		if statusReporter.ShouldUpdateStatus() {
			statusReporter.DisplayStatus(r, currentTimeMs)
		}
	}

	return nil
}
