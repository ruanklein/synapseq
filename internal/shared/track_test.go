/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package shared

import (
	"testing"

	t "github.com/ruanklein/synapseq/internal/types"
)

func TestIsTrackEqual(ts *testing.T) {
	base := &t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(20),
		Waveform:  t.WaveformSine,
		Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)},
	}

	tests := []struct {
		name string
		a, b *t.Track
		eq   bool
	}{
		{
			name: "identical",
			a:    base,
			b: &t.Track{
				Type:      t.TrackBinauralBeat,
				Carrier:   300,
				Resonance: 10,
				Amplitude: t.AmplitudePercentToRaw(20),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)},
			},
			eq: true,
		},
		{
			name: "different amplitude",
			a:    base,
			b: &t.Track{
				Type:      t.TrackBinauralBeat,
				Carrier:   300,
				Resonance: 10,
				Amplitude: t.AmplitudePercentToRaw(30),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)},
			},
			eq: false,
		},
		{
			name: "different carrier",
			a:    base,
			b: &t.Track{
				Type:      t.TrackBinauralBeat,
				Carrier:   320,
				Resonance: 10,
				Amplitude: t.AmplitudePercentToRaw(20),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)},
			},
			eq: false,
		},
		{
			name: "different resonance",
			a:    base,
			b: &t.Track{
				Type:      t.TrackBinauralBeat,
				Carrier:   300,
				Resonance: 12,
				Amplitude: t.AmplitudePercentToRaw(20),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)},
			},
			eq: false,
		},
		{
			name: "different waveform",
			a:    base,
			b: &t.Track{
				Type:      t.TrackBinauralBeat,
				Carrier:   300,
				Resonance: 10,
				Amplitude: t.AmplitudePercentToRaw(20),
				Waveform:  t.WaveformTriangle,
				Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)},
			},
			eq: false,
		},
		{
			name: "different intensity",
			a:    base,
			b: &t.Track{
				Type:      t.TrackBinauralBeat,
				Carrier:   300,
				Resonance: 10,
				Amplitude: t.AmplitudePercentToRaw(20),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(50)},
			},
			eq: false,
		},
		{
			name: "different type",
			a:    base,
			b: &t.Track{
				Type:      t.TrackMonauralBeat,
				Carrier:   300,
				Resonance: 10,
				Amplitude: t.AmplitudePercentToRaw(20),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectOff, Intensity: t.IntensityPercentToRaw(25)},
			},
			eq: false,
		},
		{
			name: "background effect type ignored",
			a: &t.Track{
				Type:      t.TrackBackground,
				Carrier:   200,
				Resonance: 5,
				Amplitude: t.AmplitudePercentToRaw(40),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectSpin, Intensity: t.IntensityPercentToRaw(60)},
			},
			b: &t.Track{
				Type:      t.TrackBackground,
				Carrier:   200,
				Resonance: 5,
				Amplitude: t.AmplitudePercentToRaw(40),
				Waveform:  t.WaveformSine,
				Effect:    t.Effect{Type: t.EffectPulse, Intensity: t.IntensityPercentToRaw(60)},
			},
			eq: true,
		},
	}

	for _, tc := range tests {
		got := IsTrackEqual(tc.a, tc.b)
		if got != tc.eq {
			ts.Errorf("%s: expected %v, got %v\nA=%+v\nB=%+v", tc.name, tc.eq, got, *tc.a, *tc.b)
		}
	}
}
