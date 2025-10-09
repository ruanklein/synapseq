/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package audio

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	s "github.com/ruanklein/synapseq/internal/shared"
	t "github.com/ruanklein/synapseq/internal/types"
)

func captureStderr(fn func()) string {
	orig := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	fn()
	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	os.Stderr = orig
	return buf.String()
}

func TestStatusReporter_DisplayPeriodChange_PrintsStartAndDash(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000

	start := t.Track{Type: t.TrackBinauralBeat, Carrier: 100, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	endEqual := start
	p0.TrackStart[0] = start
	p0.TrackEnd[0] = endEqual

	r := &AudioRenderer{periods: []t.Period{p0, p1}, AudioRendererOptions: &AudioRendererOptions{}}

	sr := NewStatusReporter(false)
	out := captureStderr(func() { sr.DisplayPeriodChange(r, 0) })

	if !strings.Contains(out, "- "+p0.TimeString()+" -> "+p1.TimeString()+" ("+p0.Transition.String()+")") {
		ts.Fatalf("missing start time line: %q", out)
	}
	// We no longer print the end time when start==end
	// if !strings.Contains(out, "  "+p1.TimeString()) {
	// 	ts.Fatalf("missing end time line: %q", out)
	// }
	if !strings.Contains(out, start.String()) {
		ts.Fatalf("missing start track string in output: %q", out)
	}
	// if !strings.Contains(out, "\n       --") {
	// 	ts.Fatalf("expected '--' marker when start==end: %q", out)
	// }
}

func TestStatusReporter_DisplayPeriodChange_ShowsEndTrackWhenChanged(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000

	start := t.Track{Type: t.TrackBinauralBeat, Carrier: 100, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	endChanged := start
	endChanged.Amplitude = t.AmplitudePercentToRaw(20)
	// sanity: ensure IsTrackEqual detects difference
	if s.IsTrackEqual(&start, &endChanged) {
		ts.Fatalf("precondition failed: start and end should not be equal")
	}
	p0.TrackStart[0] = start
	p0.TrackEnd[0] = endChanged

	r := &AudioRenderer{periods: []t.Period{p0, p1}, AudioRendererOptions: &AudioRendererOptions{}}
	sr := NewStatusReporter(false)
	out := captureStderr(func() { sr.DisplayPeriodChange(r, 0) })

	if strings.Contains(out, "\n       --") {
		ts.Fatalf("did not expect '--' when start!=end: %q", out)
	}
	if !strings.Contains(out, endChanged.String()) {
		ts.Fatalf("missing end track string when changed: %q", out)
	}
}

func TestStatusReporter_Quiet_SuppressesOutput(ts *testing.T) {
	var p0, p1 t.Period
	p0.Time = 0
	p1.Time = 1000
	tr := t.Track{Type: t.TrackBinauralBeat, Carrier: 100, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	p0.TrackStart[0] = tr
	p0.TrackEnd[0] = tr

	r := &AudioRenderer{periods: []t.Period{p0, p1}, AudioRendererOptions: &AudioRendererOptions{}}
	sr := NewStatusReporter(true)
	out := captureStderr(func() { sr.DisplayPeriodChange(r, 0) })
	if out != "" {
		ts.Fatalf("expected no output in quiet mode, got: %q", out)
	}
}

func TestStatusReporter_CheckPeriodChange_DetectsTransitions(ts *testing.T) {
	var p0, p1, p2 t.Period
	p0.Time = 0
	p1.Time = 1000
	p2.Time = 2000
	tr := t.Track{Type: t.TrackBinauralBeat, Carrier: 100, Resonance: 5, Amplitude: t.AmplitudePercentToRaw(10), Waveform: t.WaveformSine}
	p0.TrackStart[0] = tr
	p0.TrackEnd[0] = tr
	p1.TrackStart[0] = tr
	p1.TrackEnd[0] = tr

	r := &AudioRenderer{periods: []t.Period{p0, p1, p2}, AudioRendererOptions: &AudioRendererOptions{}}
	sr := NewStatusReporter(false)

	out1 := captureStderr(func() { sr.CheckPeriodChange(r, 0) })
	if !strings.Contains(out1, "- "+p0.TimeString()) {
		ts.Fatalf("expected period 0 output on first check: %q", out1)
	}

	out2 := captureStderr(func() { sr.CheckPeriodChange(r, 0) })
	if out2 != "" {
		ts.Fatalf("expected no output when period index unchanged, got: %q", out2)
	}

	out3 := captureStderr(func() { sr.CheckPeriodChange(r, 1) })
	if !strings.Contains(out3, "- "+p1.TimeString()) {
		ts.Fatalf("expected period 1 output after change: %q", out3)
	}
}
