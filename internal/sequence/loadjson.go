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
func initializeTracks(typ t.TrackType) [t.NumberOfChannels]t.Track {
	var tracks [t.NumberOfChannels]t.Track
	for ch := range t.NumberOfChannels {
		tracks[ch].Type = typ
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

	gainLevel := t.GainLevelMedium
	if jsonInput.Options.GainLevel != "" {
		switch jsonInput.Options.GainLevel {
		case t.KeywordOptionGainLevelVeryLow:
			gainLevel = t.GainLevelVeryLow
		case t.KeywordOptionGainLevelLow:
			gainLevel = t.GainLevelLow
		case t.KeywordOptionGainLevelMedium:
			gainLevel = t.GainLevelMedium
		case t.KeywordOptionGainLevelHigh:
			gainLevel = t.GainLevelHigh
		case t.KeywordOptionGainLevelVeryHigh:
			gainLevel = t.GainLevelVeryHigh
		default:
			return nil, fmt.Errorf("invalid gain level: %s", jsonInput.Options.GainLevel)
		}
	}

	// Initialize audio options
	options := &t.Option{
		SampleRate:     jsonInput.Options.Samplerate,
		Volume:         jsonInput.Options.Volume,
		BackgroundPath: jsonInput.Options.BackgroundPath,
		PresetPath:     jsonInput.Options.PresetPath,
		GainLevel:      gainLevel,
	}

	if err := options.Validate(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	// Initialize built-in silence
	silenceTracks := initializeTracks(t.TrackSilence)
	// Initialize off tracks
	offTracks := initializeTracks(t.TrackOff)

	var periods []t.Period
	for idx, seq := range jsonInput.Sequence {
		if idx == 0 && seq.Time != 0 {
			return nil, fmt.Errorf("first timeline must start at 0ms (00:00:00)")
		}
		if idx >= 1 && seq.Time <= jsonInput.Sequence[idx-1].Time {
			return nil, fmt.Errorf("timeline %d time must be greater than previous timeline time", idx+1)
		}

		if len(seq.Elements) == 0 {
			period := t.Period{
				Time:       seq.Time,
				TrackStart: silenceTracks,
				TrackEnd:   silenceTracks,
			}
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
			continue
		}

		tracks := offTracks
		for ch, el := range seq.Elements {
			if el.Kind == "" {
				continue
			}

			var mode t.TrackType
			// Get mode
			switch el.Mode {
			case t.KeywordBinaural:
				mode = t.TrackBinauralBeat
			case t.KeywordMonaural:
				mode = t.TrackMonauralBeat
			case t.KeywordIsochronic:
				mode = t.TrackIsochronicBeat
			case t.KeywordWhite:
				mode = t.TrackWhiteNoise
			case t.KeywordPink:
				mode = t.TrackPinkNoise
			case t.KeywordBrown:
				mode = t.TrackBrownNoise
			case t.KeywordPure: // Pure tone
				mode = t.TrackPureTone
			default:
				return nil, fmt.Errorf("invalid tone mode: %s", el.Mode)
			}

			waveForm := t.WaveformSine
			if el.Kind == t.KeywordTone {
				switch el.Waveform {
				case t.KeywordSine:
					waveForm = t.WaveformSine
				case t.KeywordSquare:
					waveForm = t.WaveformSquare
				case t.KeywordTriangle:
					waveForm = t.WaveformTriangle
				case t.KeywordSawtooth:
					waveForm = t.WaveformSawtooth
				default:
					return nil, fmt.Errorf("invalid waveform type: %s", el.Waveform)
				}

				if el.Carrier < 0 {
					return nil, fmt.Errorf("carrier frequency must be greater than 0Hz")
				}
				if el.Resonance < 0 {
					return nil, fmt.Errorf("resonance frequency cannot be negative")
				}
			}

			if el.Amplitude < 0 || el.Amplitude > 100 {
				return nil, fmt.Errorf("amplitude must be between 0 and 100")
			}

			tracks[ch].Type = mode
			tracks[ch].Carrier = el.Carrier
			tracks[ch].Resonance = el.Resonance
			tracks[ch].Amplitude = t.AmplitudePercentToRaw(el.Amplitude)
			tracks[ch].Waveform = waveForm
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

	return &LoadResult{
		Periods:  periods,
		Options:  options,
		Comments: nil,
	}, nil
}
