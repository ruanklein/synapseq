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

func TestFindPreset(ts *testing.T) {
	var presets []t.Preset
	presets = append(presets, *t.NewBuiltinSilencePreset())

	alpha, err := t.NewPreset("alpha", false, nil)
	if err != nil {
		ts.Fatalf("unexpected error creating preset alpha: %v", err)
	}
	beta, err := t.NewPreset("beta", false, nil)
	if err != nil {
		ts.Fatalf("unexpected error creating preset beta: %v", err)
	}
	presets = append(presets, *alpha, *beta)

	got := FindPreset("alpha", presets)
	if got == nil || got != &presets[1] || got.String() != "alpha" {
		ts.Fatalf("FindPreset('alpha') failed, got=%v, want address of presets[1]", got)
	}
	got = FindPreset("silence", presets)
	if got == nil || got != &presets[0] || got.String() != "silence" {
		ts.Fatalf("FindPreset('silence') failed, got=%v, want address of presets[0]", got)
	}
	got = FindPreset("missing", presets)
	if got != nil {
		ts.Fatalf("FindPreset('missing') should be nil, got=%v", got)
	}
}

func TestAllocateTrack(ts *testing.T) {
	p, err := t.NewPreset("alpha", false, nil)
	if err != nil {
		ts.Fatalf("unexpected error: %v", err)
	}

	for i := range t.NumberOfChannels {
		idx, err := AllocateTrack(p)
		if err != nil {
			ts.Fatalf("AllocateTrack failed at i=%d: %v", i, err)
		}
		if idx != i {
			ts.Fatalf("AllocateTrack index mismatch: got %d, want %d", idx, i)
		}
		p.Track[idx].Type = t.TrackBinauralBeat
	}

	if _, err := AllocateTrack(p); err == nil {
		ts.Fatalf("AllocateTrack should fail when no free tracks")
	}
}

func TestIsPresetEmpty(ts *testing.T) {
	p, err := t.NewPreset("alpha", false, nil)
	if err != nil {
		ts.Fatalf("unexpected error: %v", err)
	}
	if !IsPresetEmpty(p) {
		ts.Fatalf("new preset should be empty")
	}

	p.Track[0].Type = t.TrackWhiteNoise
	if IsPresetEmpty(p) {
		ts.Fatalf("preset with one active track should not be empty")
	}

	sil := t.NewBuiltinSilencePreset()
	if IsPresetEmpty(sil) {
		ts.Fatalf("silence preset should not be considered empty")
	}
}

func TestNumBackgroundTracks(ts *testing.T) {
	p, err := t.NewPreset("alpha", false, nil)
	if err != nil {
		ts.Fatalf("unexpected error: %v", err)
	}
	if NumBackgroundTracks(p) != 0 {
		ts.Fatalf("expected 0 background tracks initially")
	}

	p.Track[0].Type = t.TrackBackground
	if NumBackgroundTracks(p) != 1 {
		ts.Fatalf("expected 1 background track")
	}

	if len(p.Track) > 2 {
		p.Track[2].Type = t.TrackBackground
	}
	want := 2
	if len(p.Track) <= 2 {
		want = 1
	}
	if NumBackgroundTracks(p) != want {
		ts.Fatalf("unexpected background track count: got %d, want %d", NumBackgroundTracks(p), want)
	}

	sil := t.NewBuiltinSilencePreset()
	if NumBackgroundTracks(sil) != 0 {
		ts.Fatalf("silence preset should have 0 background tracks")
	}
}
