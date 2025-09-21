/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package audio

import (
	"math"
	"testing"

	t "github.com/ruanklein/synapseq/internal/types"
)

func genWhiteExpected(n int) []int {
	// Replicates the LCG and scaling used in generateWhiteNoise
	s := 2
	scale := int(t.WaveTableAmplitude / 65535)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		s = (s * randMult) % 131074
		out[i] = (s - 65535) * scale
	}
	return out
}

func meanAbsDelta(vals []int) float64 {
	if len(vals) < 2 {
		return 0
	}
	var sum float64
	for i := 1; i < len(vals); i++ {
		d := float64(vals[i] - vals[i-1])
		if d < 0 {
			d = -d
		}
		sum += d
	}
	return sum / float64(len(vals)-1)
}

func TestNoise_White_DeterministicFirstValues(ts *testing.T) {
	ng := NewNoiseGenerator()
	exp := genWhiteExpected(8)
	got := make([]int, len(exp))
	for i := range got {
		got[i] = ng.Generate(t.TrackWhiteNoise)
	}
	for i := range exp {
		if got[i] != exp[i] {
			ts.Fatalf("white[%d]: got %d want %d", i, got[i], exp[i])
		}
	}
}

func TestNoise_BoundsAndSign(ts *testing.T) {
	ng := NewNoiseGenerator()
	N := 16384
	types := []t.TrackType{t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise}
	for _, typ := range types {
		pos, neg := 0, 0
		for i := 0; i < N; i++ {
			v := ng.Generate(typ)
			if v > int(t.WaveTableAmplitude) || v < -int(t.WaveTableAmplitude) {
				ts.Fatalf("%v: out of bounds value %d at i=%d", typ, v, i)
			}
			if v > 0 {
				pos++
			} else if v < 0 {
				neg++
			}
		}
		// Ensure both signs appear reasonably
		if pos == 0 || neg == 0 {
			ts.Fatalf("%v: expected both positive and negative values (pos=%d neg=%d)", typ, pos, neg)
		}
	}
}

func TestNoise_SmoothnessComparison(ts *testing.T) {
	N := 16384

	// White
	ngW := NewNoiseGenerator()
	white := make([]int, N)
	for i := 0; i < N; i++ {
		white[i] = ngW.Generate(t.TrackWhiteNoise)
	}

	// Pink
	ngP := NewNoiseGenerator()
	pink := make([]int, N)
	for i := 0; i < N; i++ {
		pink[i] = ngP.Generate(t.TrackPinkNoise)
	}

	// Brown
	ngB := NewNoiseGenerator()
	brown := make([]int, N)
	for i := 0; i < N; i++ {
		brown[i] = ngB.Generate(t.TrackBrownNoise)
	}

	dW := meanAbsDelta(white)
	dP := meanAbsDelta(pink)
	dB := meanAbsDelta(brown)

	// Brown should be smoothest (lowest delta), then Pink, then White highest
	if !(dB < dP && dP < dW) {
		ts.Fatalf("unexpected smoothness order: dB=%.2f dP=%.2f dW=%.2f", dB, dP, dW)
	}

	// Also ensure deltas are non-trivial
	if math.IsNaN(dW) || math.IsNaN(dP) || math.IsNaN(dB) {
		ts.Fatalf("NaN deltas: dB=%.2f dP=%.2f dW=%.2f", dB, dP, dW)
	}
}

func TestNoise_DeterminismAcrossInstances(ts *testing.T) {
	N := 2048
	types := []t.TrackType{t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise}
	for _, typ := range types {
		ng1 := NewNoiseGenerator()
		ng2 := NewNoiseGenerator()
		for i := 0; i < N; i++ {
			v1 := ng1.Generate(typ)
			v2 := ng2.Generate(typ)
			if v1 != v2 {
				ts.Fatalf("%v: non-deterministic at i=%d: %d vs %d", typ, i, v1, v2)
			}
		}
	}
}
