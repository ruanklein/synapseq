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
	"testing"

	t "github.com/ruanklein/synapseq/internal/types"
)

func writeTemp(ts *testing.T, name, content string) string {
	ts.Helper()
	dir := ts.TempDir()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		ts.Fatalf("write temp file %s: %v", name, err)
	}
	return p
}

func assertIncreasing(ts *testing.T, times []int) {
	ts.Helper()
	if len(times) == 0 {
		ts.Fatalf("empty times")
	}
	if times[0] != 0 {
		ts.Fatalf("first time must be 0, got %d", times[0])
	}
	for i := 1; i < len(times); i++ {
		if !(times[i] > times[i-1]) {
			ts.Fatalf("times not strictly increasing at %d: %v", i, times)
		}
	}
}

func periodTimes(periods []t.Period) []int {
	out := make([]int, len(periods))
	for i, p := range periods {
		out[i] = p.Time
	}
	return out
}

func hasTrackIn(tracks [t.NumberOfChannels]t.Track, want t.Track) bool {
	for i := range tracks {
		got := tracks[i]
		if got.Type == want.Type &&
			got.Carrier == want.Carrier &&
			got.Resonance == want.Resonance &&
			got.Amplitude == want.Amplitude &&
			got.Waveform == want.Waveform {
			return true
		}
	}
	return false
}

func TestLoadStructured_JSON_Standalone(ts *testing.T) {
	json := `{
  "description": ["Standalone structured test"],
  "options": { "samplerate": 44100, "volume": 100 },
  "sequence": [
    {
      "time": 0,
      "track": {
        "tones": [
          { "mode": "binaural", "carrier": 250, "resonance": 8, "amplitude": 0, "waveform": "sine" }
        ],
        "noises": [
          { "mode": "pink", "amplitude": 0 }
        ]
      }
    },
    {
      "time": 15000,
      "track": {
        "tones": [
          { "mode": "binaural", "carrier": 250, "resonance": 8, "amplitude": 15, "waveform": "sine" }
        ],
        "noises": [
          { "mode": "pink", "amplitude": 30 }
        ]
      }
    }
  ]
}`
	p := writeTemp(ts, "seq.json", json)

	res, err := LoadStructuredSequence(p, "json")
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(json) error: %v", err)
	}

	if res.Options.SampleRate != 44100 || res.Options.Volume != 100 {
		ts.Fatalf("unexpected options: %+v", *res.Options)
	}
	if len(res.Comments) == 0 {
		ts.Fatalf("expected non-empty description/comments")
	}

	if len(res.Periods) != 2 {
		ts.Fatalf("expected 2 periods, got %d", len(res.Periods))
	}
	assertIncreasing(ts, periodTimes(res.Periods))

	wantTone := t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSine,
	}
	wantNoise := t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(30),
	}
	if !hasTrackIn(res.Periods[1].TrackStart, wantTone) {
		ts.Fatalf("missing expected tone in period[1]")
	}
	if !hasTrackIn(res.Periods[1].TrackStart, wantNoise) {
		ts.Fatalf("missing expected noise in period[1]")
	}
}

func TestLoadStructured_XML_Standalone(ts *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<SynapSeqInput>
  <description>
    <line>Standalone structured test</line>
  </description>
  <options>
    <samplerate>44100</samplerate>
    <volume>100</volume>
  </options>
  <sequence>
    <entry time="0">
      <track>
        <tone mode="binaural" carrier="250" resonance="8" amplitude="0" waveform="sine"/>
        <noise mode="pink" amplitude="0"/>
      </track>
    </entry>
    <entry time="15000">
      <track>
        <tone mode="binaural" carrier="250" resonance="8" amplitude="15" waveform="sine"/>
        <noise mode="pink" amplitude="30"/>
      </track>
    </entry>
  </sequence>
</SynapSeqInput>`
	p := writeTemp(ts, "seq.xml", xml)

	res, err := LoadStructuredSequence(p, "xml")
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(xml) error: %v", err)
	}

	if res.Options.SampleRate != 44100 || res.Options.Volume != 100 {
		ts.Fatalf("unexpected options: %+v", *res.Options)
	}
	if len(res.Comments) == 0 {
		ts.Fatalf("expected non-empty description/comments")
	}

	if len(res.Periods) != 2 {
		ts.Fatalf("expected 2 periods, got %d", len(res.Periods))
	}
	assertIncreasing(ts, periodTimes(res.Periods))

	wantTone := t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSine,
	}
	wantNoise := t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(30),
	}
	if !hasTrackIn(res.Periods[1].TrackStart, wantTone) {
		ts.Fatalf("missing expected tone in period[1]")
	}
	if !hasTrackIn(res.Periods[1].TrackStart, wantNoise) {
		ts.Fatalf("missing expected noise in period[1]")
	}
}

func TestLoadStructured_YAML_Standalone(ts *testing.T) {
	yaml := `description:
  - Standalone structured test
options:
  samplerate: 44100
  volume: 100
sequence:
  - time: 0
    track:
      tones:
        - mode: binaural
          carrier: 250
          resonance: 8
          amplitude: 0
          waveform: sine
      noises:
        - mode: pink
          amplitude: 0
  - time: 15000
    track:
      tones:
        - mode: binaural
          carrier: 250
          resonance: 8
          amplitude: 15
          waveform: sine
      noises:
        - mode: pink
          amplitude: 30
`
	p := writeTemp(ts, "seq.yaml", yaml)

	res, err := LoadStructuredSequence(p, "yaml")
	if err != nil {
		ts.Fatalf("LoadStructuredSequence(yaml) error: %v", err)
	}

	if res.Options.SampleRate != 44100 || res.Options.Volume != 100 {
		ts.Fatalf("unexpected options: %+v", *res.Options)
	}
	if len(res.Comments) == 0 {
		ts.Fatalf("expected non-empty description/comments")
	}

	if len(res.Periods) != 2 {
		ts.Fatalf("expected 2 periods, got %d", len(res.Periods))
	}
	assertIncreasing(ts, periodTimes(res.Periods))

	wantTone := t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSine,
	}
	wantNoise := t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(30),
	}
	if !hasTrackIn(res.Periods[1].TrackStart, wantTone) {
		ts.Fatalf("missing expected tone in period[1]")
	}
	if !hasTrackIn(res.Periods[1].TrackStart, wantNoise) {
		ts.Fatalf("missing expected noise in period[1]")
	}
}
