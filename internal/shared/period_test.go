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

func TestAdjustPeriods_NormalCopy(ts *testing.T) {
	var last, next t.Period

	last.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(10),
		Waveform:  t.WaveformSine,
	}
	last.TrackEnd[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   310,
		Resonance: 11,
		Amplitude: t.AmplitudePercentToRaw(12),
		Waveform:  t.WaveformSine,
	}
	next.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   350,
		Resonance: 12,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSine,
	}

	if err := AdjustPeriods(&last, &next); err != nil {
		ts.Fatalf("unexpected error: %v", err)
	}
	if last.TrackEnd[0] != next.TrackStart[0] {
		ts.Fatalf("carry-forward mismatch: last.TrackEnd != next.TrackStart\nlast=%+v\nnext=%+v", last.TrackEnd[0], next.TrackStart[0])
	}
}

func TestAdjustPeriods_FadeInFromSilence(ts *testing.T) {
	var last, next t.Period

	last.TrackStart[0] = t.Track{
		Type:     t.TrackSilence,
		Waveform: t.WaveformSquare,
	}
	last.TrackEnd[0] = t.Track{
		Type:      t.TrackSilence,
		Amplitude: 0,
	}
	next.TrackStart[0] = t.Track{
		Type:      t.TrackMonauralBeat,
		Carrier:   200,
		Resonance: 6,
		Amplitude: t.AmplitudePercentToRaw(25),
		Waveform:  t.WaveformTriangle,
	}

	if err := AdjustPeriods(&last, &next); err != nil {
		ts.Fatalf("unexpected error: %v", err)
	}

	got := last.TrackStart[0]
	if got.Type != t.TrackMonauralBeat || got.Amplitude != 0 || got.Carrier != 200 || got.Resonance != 6 || got.Waveform != t.WaveformTriangle {
		ts.Fatalf("fade-in not applied as expected: %+v", got)
	}
	if last.TrackEnd[0] != next.TrackStart[0] {
		ts.Fatalf("carry-forward mismatch after fade-in: last.TrackEnd != next.TrackStart\nlast=%+v\nnext=%+v", last.TrackEnd[0], next.TrackStart[0])
	}
}

func TestAdjustPeriods_FadeOutToSilence(ts *testing.T) {
	var last, next t.Period

	last.TrackStart[0] = t.Track{
		Type:      t.TrackBackground,
		Carrier:   200,
		Resonance: 5,
		Amplitude: t.AmplitudePercentToRaw(50),
		Waveform:  t.WaveformSquare,
		Effect:    t.Effect{Type: t.EffectSpin, Intensity: t.IntensityPercentToRaw(70)},
	}
	last.TrackEnd[0] = last.TrackStart[0]
	next.TrackStart[0] = t.Track{
		Type:      t.TrackSilence,
		Amplitude: 0,
	}

	if err := AdjustPeriods(&last, &next); err != nil {
		ts.Fatalf("unexpected error: %v", err)
	}

	// Fade-out should copy carrier/resonance/intensity into the silence start
	if next.TrackStart[0].Type != t.TrackSilence ||
		next.TrackStart[0].Carrier != 200 ||
		next.TrackStart[0].Resonance != 5 ||
		next.TrackStart[0].Intensity != t.IntensityPercentToRaw(70) {
		ts.Fatalf("fade-out not applied as expected: %+v", next.TrackStart[0])
	}
	if last.TrackEnd[0] != next.TrackStart[0] {
		ts.Fatalf("carry-forward mismatch after fade-out: last.TrackEnd != next.TrackStart\nlast=%+v\nnext=%+v", last.TrackEnd[0], next.TrackStart[0])
	}
}

func TestAdjustPeriods_Errors(ts *testing.T) {
	makePer := func(tr0, tr1, tr2 t.Track) (t.Period, t.Period) {
		var last, next t.Period
		last.TrackStart[0] = tr0
		last.TrackEnd[0] = tr1
		next.TrackStart[0] = tr2
		return last, next
	}

	tests := []struct {
		name string
		tr0  t.Track
		tr1  t.Track
		tr2  t.Track
	}{
		{
			name: "turn off directly",
			tr1:  t.Track{Type: t.TrackBinauralBeat, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine},
			tr2:  t.Track{Type: t.TrackOff},
		},
		{
			name: "turn on directly",
			tr1:  t.Track{Type: t.TrackOff},
			tr2:  t.Track{Type: t.TrackBinauralBeat, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine},
		},
		{
			name: "change type while on",
			tr1:  t.Track{Type: t.TrackBinauralBeat, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine},
			tr2:  t.Track{Type: t.TrackMonauralBeat, Amplitude: t.AmplitudePercentToRaw(12), Waveform: t.WaveformSine},
		},
		{
			name: "change waveform while on",
			tr1:  t.Track{Type: t.TrackBinauralBeat, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine},
			tr2:  t.Track{Type: t.TrackBinauralBeat, Amplitude: t.AmplitudePercentToRaw(12), Waveform: t.WaveformTriangle},
		},
		{
			name: "change effect type while on (background)",
			tr1:  t.Track{Type: t.TrackBackground, Amplitude: t.AmplitudePercentToRaw(20), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectSpin}},
			tr2:  t.Track{Type: t.TrackBackground, Amplitude: t.AmplitudePercentToRaw(25), Waveform: t.WaveformSine, Effect: t.Effect{Type: t.EffectPulse}},
		},
	}

	for _, tc := range tests {
		last, next := makePer(tc.tr0, tc.tr1, tc.tr2)
		if err := AdjustPeriods(&last, &next); err == nil {
			ts.Fatalf("%s: expected error, got nil", tc.name)
		}
	}
}
