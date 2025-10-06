/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package sequence

import (
	"encoding/json"
	"fmt"
	"os"

	s "github.com/ruanklein/synapseq/internal/shared"
	t "github.com/ruanklein/synapseq/internal/types"
)

// Helper to initialize tracks
func initializeTracks() [t.NumberOfChannels]t.Track {
	var tracks [t.NumberOfChannels]t.Track
	for ch := range t.NumberOfChannels {
		tracks[ch].Type = t.TrackOff
		tracks[ch].Carrier = 0.0
		tracks[ch].Resonance = 0.0
		tracks[ch].Amplitude = 0.0
		tracks[ch].Waveform = t.WaveformSine
		tracks[ch].Effect = t.Effect{Type: t.EffectOff, Intensity: 0.0}
	}
	return tracks
}

// LoadJSONSequence loads and parses a JSON sequence file
func LoadJSONSequence(filename string) (*LoadResult, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	var jsonInput t.SynapSeqInput
	if err := json.Unmarshal(data, &jsonInput); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	if len(jsonInput.Sequence) == 0 {
		return nil, fmt.Errorf("no sequence data found in JSON")
	}

	// Initialize audio options
	options := &t.Option{
		SampleRate: jsonInput.Options.Samplerate,
		Volume:     jsonInput.Options.Volume,
	}

	if err := options.Validate(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	// Initialize off tracks
	offTracks := initializeTracks()

	var periods []t.Period
	for idx, seq := range jsonInput.Sequence {
		if idx == 0 && seq.Time != 0 {
			return nil, fmt.Errorf("first timeline must start at 0ms (00:00:00)")
		}
		if idx >= 1 && seq.Time <= jsonInput.Sequence[idx-1].Time {
			return nil, fmt.Errorf("timeline %d time must be greater than previous timeline time", idx+1)
		}
		if len(seq.Elements.Tones)+len(seq.Elements.Noises) > t.NumberOfChannels {
			return nil, fmt.Errorf("too many elements defined (max %d)", t.NumberOfChannels)
		}

		tracks := offTracks
		trackIdx := 0

		for _, tone := range seq.Elements.Tones {
			var mode t.TrackType
			// Get mode
			switch tone.Mode {
			case t.KeywordBinaural:
				mode = t.TrackBinauralBeat
			case t.KeywordMonaural:
				mode = t.TrackMonauralBeat
			case t.KeywordIsochronic:
				mode = t.TrackIsochronicBeat
			case t.KeywordPure:
				mode = t.TrackPureTone
			default:
				return nil, fmt.Errorf("invalid tone mode: %s", tone.Mode)
			}

			var waveForm t.WaveformType
			switch tone.Waveform {
			case t.KeywordSine:
				waveForm = t.WaveformSine
			case t.KeywordSquare:
				waveForm = t.WaveformSquare
			case t.KeywordTriangle:
				waveForm = t.WaveformTriangle
			case t.KeywordSawtooth:
				waveForm = t.WaveformSawtooth
			default:
				return nil, fmt.Errorf("invalid waveform type: %s", tone.Waveform)
			}

			tr := t.Track{
				Type:      mode,
				Carrier:   tone.Carrier,
				Resonance: tone.Resonance,
				Amplitude: t.AmplitudePercentToRaw(tone.Amplitude),
				Waveform:  waveForm,
			}

			if err := tr.Validate(); err != nil {
				return nil, fmt.Errorf("%v", err)
			}

			tracks[trackIdx] = tr
			trackIdx++
		}
		for _, noise := range seq.Elements.Noises {
			var mode t.TrackType
			// Get mode
			switch noise.Mode {
			case t.KeywordWhite:
				mode = t.TrackWhiteNoise
			case t.KeywordPink:
				mode = t.TrackPinkNoise
			case t.KeywordBrown:
				mode = t.TrackBrownNoise
			default:
				return nil, fmt.Errorf("invalid noise mode: %s", noise.Mode)
			}

			tr := t.Track{
				Type:      mode,
				Amplitude: t.AmplitudePercentToRaw(noise.Amplitude),
			}

			if err := tr.Validate(); err != nil {
				return nil, fmt.Errorf("%v", err)
			}

			tracks[trackIdx] = tr
			trackIdx++
		}

		// Process Period
		period := t.Period{
			Time:       seq.Time,
			TrackStart: tracks,
			TrackEnd:   tracks,
		}
		// Adjust previous period end if needed
		var lastPeriod *t.Period
		if len(periods) > 0 {
			lastPeriod = &periods[len(periods)-1]
		}
		if lastPeriod != nil {
			if err := s.AdjustPeriods(lastPeriod, &period); err != nil {
				return nil, fmt.Errorf("%v", err)
			}
		}
		periods = append(periods, period)
	}

	description := jsonInput.Description

	return &LoadResult{
		Periods:  periods,
		Options:  options,
		Comments: description,
	}, nil
}
