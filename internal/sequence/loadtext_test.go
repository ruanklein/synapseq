/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package sequence

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	t "github.com/ruanklein/synapseq/internal/types"
)

func writeSeqFile(tst *testing.T, content string) string {
	tst.Helper()
	dir := tst.TempDir()
	p := filepath.Join(dir, "seq.spsq")
	if err := os.WriteFile(p, []byte(strings.TrimSpace(content)+"\n"), 0o600); err != nil {
		tst.Fatalf("write temp sequence: %v", err)
	}
	return p
}

func eqTrackGotWant(got, want t.Track) bool {
	if got.Type != want.Type {
		return false
	}
	if got.Carrier != want.Carrier {
		return false
	}
	if got.Resonance != want.Resonance {
		return false
	}
	if got.Amplitude != want.Amplitude {
		return false
	}
	if got.Waveform != want.Waveform {
		return false
	}
	if got.Effect.Type != want.Effect.Type {
		return false
	}
	if got.Effect.Intensity != want.Effect.Intensity {
		return false
	}
	return true
}

func hasTrack(tracks [t.NumberOfChannels]t.Track, want t.Track) bool {
	for ch := range t.NumberOfChannels {
		if eqTrackGotWant(tracks[ch], want) {
			return true
		}
	}
	return false
}

func TestLoadTextSequence_Success(ts *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		ts.Fatalf("getwd: %v", err)
	}
	bgPath := filepath.Join(wd, "testdata", "noise.wav")

	seq := `
# Options
@samplerate 48000
@volume 80
@background testdata/noise.wav
@gainlevel high

# Presets
alpha
  noise brown amplitude 40
  tone 300 binaural 10 amplitude 20

beta
  noise brown amplitude 40
  tone 300 binaural 14 amplitude 15

# Timeline
00:00:00 alpha
00:01:00 beta
`
	path := writeSeqFile(ts, seq)

	result, err := LoadTextSequence(path)
	if err != nil {
		ts.Fatalf("LoadTextSequence error: %v", err)
	}
	opts := result.Options
	if opts.SampleRate != 48000 || opts.Volume != 80 || opts.GainLevel != t.GainLevelHigh {
		ts.Fatalf("unexpected options: %+v", *opts)
	}
	if opts.BackgroundPath != bgPath {
		ts.Fatalf("unexpected background path: got %q want %q", opts.BackgroundPath, bgPath)
	}

	periods := result.Periods
	if len(periods) != 2 {
		ts.Fatalf("expected 2 periods, got %d", len(periods))
	}
	if periods[0].Time != 0 || periods[1].Time != 60_000 {
		ts.Fatalf("unexpected period times: %+v", []int{periods[0].Time, periods[1].Time})
	}

	wantNoise := t.Track{
		Type:      t.TrackBrownNoise,
		Amplitude: t.AmplitudePercentToRaw(40),
	}
	wantToneAlpha := t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(20),
	}
	wantToneBeta := t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 14,
		Amplitude: t.AmplitudePercentToRaw(15),
	}

	if !hasTrack(periods[0].TrackStart, wantNoise) || !hasTrack(periods[0].TrackStart, wantToneAlpha) {
		ts.Fatalf("missing alpha tracks in period[0]: %+v", periods[0].TrackStart)
	}
	if !hasTrack(periods[1].TrackStart, wantToneBeta) {
		ts.Fatalf("missing beta tracks in period[1]: %+v", periods[1].TrackStart)
	}
}

func TestLoadTextSequence_Error_TimelineBeforePreset(ts *testing.T) {
	seq := `
# Timeline first
00:00:00 alpha
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for timeline before preset")
	}
}

func TestLoadTextSequence_Error_OptionsAfterPreset(ts *testing.T) {
	seq := `
alpha
  tone 100 binaural 1 amplitude 1
@volume 80
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for options after preset")
	}
}

func TestLoadTextSequence_Error_FirstTimelineNotZero(ts *testing.T) {
	seq := `
alpha
  tone 100 binaural 1 amplitude 1
00:00:15 alpha
00:01:00 alpha
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for first timeline not 00:00:00")
	}
}

func TestLoadTextSequence_Error_OverlappingTimeline(ts *testing.T) {
	seq := `
alpha
  tone 100 binaural 1 amplitude 1
00:00:00 alpha
00:01:00 alpha
00:00:30 alpha
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for overlapping timeline")
	}
}

func TestLoadTextSequence_Error_BackgroundWithoutFile(ts *testing.T) {
	seq := `
alpha
  background amplitude 10
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for background track without background option")
	}
}

func TestLoadTextSequence_Error_DuplicatePreset(ts *testing.T) {
	seq := `
alpha
alpha
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for duplicate preset")
	}
}

func TestLoadTextSequence_Error_MaxPresets(ts *testing.T) {
	var b strings.Builder
	for i := 1; i <= t.MaxPresets-1; i++ {
		b.WriteString("p")
		b.WriteString(strings.Repeat("x", i%5+1))
		b.WriteString("\n")
	}
	b.WriteString("overflow\n")
	path := writeSeqFile(ts, b.String())
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for maximum number of presets reached")
	}
}

func TestLoadTextSequence_Comments(ts *testing.T) {
	seq := `
## This is a comment on at the top

# Presets
alpha
  tone 100 binaural 1 amplitude 1

## Another comment line

# Timeline
00:00:00 alpha
##Another comment line in between
00:01:00 alpha
`
	path := writeSeqFile(ts, seq)
	result, err := LoadTextSequence(path)
	if err != nil {
		ts.Fatalf("LoadTextSequence with comments error: %v", err)
	}

	cmms := result.Comments
	if cmms == nil || len(cmms) != 3 {
		ts.Fatalf("expected 2 comment blocks, got %d", len(cmms))
	}
	if cmms[0] != "This is a comment on at the top" {
		ts.Fatalf("unexpected first comment block: %q", cmms[0])
	}
	if cmms[1] != "Another comment line" {
		ts.Fatalf("unexpected second comment block: %q", cmms[1])
	}
	if cmms[2] != "Another comment line in between" {
		ts.Fatalf("unexpected third comment block: %q", cmms[2])
	}
}

func TestLoadTextSequence_Error_PresetEmpty(ts *testing.T) {
	seq := `
alpha
00:00:00 alpha
00:01:00 alpha
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for empty preset")
	}
}

func TestLoadTextSequence_Error_MultipleBackgrounds(ts *testing.T) {
	wd, _ := os.Getwd()
	bg := filepath.Join(wd, "testdata", "noise.wav")
	if _, err := os.Stat(bg); err != nil {
		ts.Fatalf("missing test bg file: %v", err)
	}
	seq := `
@background testdata/noise.wav
alpha
  background amplitude 20
  background amplitude 30
00:00:00 alpha
00:01:00 alpha
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for multiple background tracks in a preset")
	}
}

func TestLoadTextSequence_Error_AtLeastTwoPeriods(ts *testing.T) {
	seq := `
alpha
  tone 100 binaural 1 amplitude 1
00:00:00 alpha
`
	path := writeSeqFile(ts, seq)
	_, err := LoadTextSequence(path)
	if err == nil {
		ts.Fatalf("expected error for less than two periods")
	}
}
