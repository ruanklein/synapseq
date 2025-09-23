/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package audio

import (
	"bytes"
	"errors"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	t "github.com/ruanklein/synapseq/internal/types"
)

func TestAudioRenderer_RenderWav_Integration(ts *testing.T) {
	// Create test periods (2 seconds total) with different track types
	var p0, p1, p2 t.Period

	// Period 0: 0-500ms binaural beat
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(20),
		Waveform:  t.WaveformSine,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]

	// Period 1: 500-1000ms monaural beat with interpolation
	p1.Time = 500
	p1.TrackStart[0] = t.Track{
		Type:      t.TrackMonauralBeat,
		Carrier:   250,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformTriangle,
	}
	p1.TrackEnd[0] = t.Track{
		Type:      t.TrackMonauralBeat,
		Carrier:   280,
		Resonance: 12,
		Amplitude: t.AmplitudePercentToRaw(25),
		Waveform:  t.WaveformTriangle,
	}

	// Period 2: 1000-2000ms with multiple tracks (noise + isochronic)
	p2.Time = 1000
	p2.TrackStart[0] = t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(10),
		Waveform:  t.WaveformSine,
	}
	p2.TrackStart[1] = t.Track{
		Type:      t.TrackIsochronicBeat,
		Carrier:   40,
		Resonance: 2.5,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSawtooth,
	}
	p2.TrackEnd[0] = p2.TrackStart[0]
	p2.TrackEnd[1] = p2.TrackStart[1]

	// End period at 2s
	var pEnd t.Period
	pEnd.Time = 2000

	periods := []t.Period{p0, p1, p2, pEnd}

	options := &AudioRendererOptions{
		SampleRate:     44100,
		Volume:         80,
		GainLevel:      t.GainLevelMedium,
		BackgroundPath: "",
		Quiet:          true,
		Debug:          false,
	}

	renderer, err := NewAudioRenderer(periods, options)
	if err != nil {
		ts.Fatalf("NewAudioRenderer failed: %v", err)
	}

	// Create temp directory and output file
	tempDir := ts.TempDir()
	outPath := filepath.Join(tempDir, "test_output.wav")

	// Render to WAV
	if err := renderer.RenderWav(outPath); err != nil {
		ts.Fatalf("RenderWav failed: %v", err)
	}

	// Validate the generated WAV file
	file, err := os.Open(outPath)
	if err != nil {
		ts.Fatalf("Failed to open generated WAV: %v", err)
	}
	defer file.Close()

	decoder := wav.NewDecoder(file)
	if !decoder.IsValidFile() {
		ts.Fatalf("Generated file is not a valid WAV")
	}

	if decoder.SampleRate != uint32(options.SampleRate) {
		ts.Fatalf("Sample rate mismatch: got %d, want %d", decoder.SampleRate, options.SampleRate)
	}

	if decoder.NumChans != audioChannels {
		ts.Fatalf("Channel count mismatch: got %d, want %d", decoder.NumChans, audioChannels)
	}

	if decoder.BitDepth != audioBitDepth {
		ts.Fatalf("Bit depth mismatch: got %d, want %d", decoder.BitDepth, audioBitDepth)
	}

	// Verify file size is reasonable for 2 seconds of audio
	stat, err := file.Stat()
	if err != nil {
		ts.Fatalf("Failed to stat file: %v", err)
	}

	expectedMinSize := int64(2 * options.SampleRate * audioChannels * audioBitDepth / 8) // ~2 seconds of raw PCM
	if stat.Size() < expectedMinSize/2 {                                                 // Allow some margin for headers/compression
		ts.Fatalf("Generated file too small: got %d bytes, expected at least %d", stat.Size(), expectedMinSize/2)
	}

	// Read and verify some audio data exists (non-zero samples)
	decoder.Rewind()
	audioBuf, err := decoder.FullPCMBuffer()
	if err != nil {
		ts.Fatalf("Failed to read audio data: %v", err)
	}

	if len(audioBuf.Data) == 0 {
		ts.Fatalf("Audio buffer is empty")
	}

	// Check that we have non-zero samples (basic sanity check)
	hasNonZero := false
	for _, sample := range audioBuf.Data[:1000] { // Check first 1000 samples
		if sample != 0 {
			hasNonZero = true
			break
		}
	}
	if !hasNonZero {
		ts.Fatalf("All samples are zero - audio generation may be broken")
	}
}

func TestAudioRenderer_RenderWav_WithBackground(ts *testing.T) {
	// Create a simple test WAV file as background
	tempDir := ts.TempDir()
	bgPath := filepath.Join(tempDir, "background.wav")

	// Create a minimal background WAV file
	bgFile, err := os.Create(bgPath)
	if err != nil {
		ts.Fatalf("Failed to create background file: %v", err)
	}

	bgEnc := wav.NewEncoder(bgFile, 44100, audioBitDepth, audioChannels, 1)

	// Generate 1 second of simple background audio
	bgSamples := make([]int, 44100*audioChannels) // 1 second stereo
	for i := range bgSamples {
		bgSamples[i] = 1000 // Simple constant value
	}

	bgBuf := &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: audioChannels,
			SampleRate:  44100,
		},
		Data:           bgSamples,
		SourceBitDepth: audioBitDepth,
	}

	if err := bgEnc.Write(bgBuf); err != nil {
		ts.Fatalf("Failed to write background: %v", err)
	}
	bgEnc.Close()
	bgFile.Close()

	// Create test period with background track
	var p0, pEnd t.Period
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackBackground,
		Amplitude: t.AmplitudePercentToRaw(30),
		Waveform:  t.WaveformSine,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]

	pEnd.Time = 1000 // 1 second
	periods := []t.Period{p0, pEnd}

	options := &AudioRendererOptions{
		SampleRate:     44100,
		Volume:         100,
		GainLevel:      t.GainLevelMedium,
		BackgroundPath: bgPath,
		Quiet:          true,
		Debug:          false,
	}

	renderer, err := NewAudioRenderer(periods, options)
	if err != nil {
		ts.Fatalf("NewAudioRenderer with background failed: %v", err)
	}

	outPath := filepath.Join(tempDir, "test_with_bg.wav")
	if err := renderer.RenderWav(outPath); err != nil {
		ts.Fatalf("RenderWav with background failed: %v", err)
	}

	// Basic validation
	if _, err := os.Stat(outPath); err != nil {
		ts.Fatalf("Output file not created: %v", err)
	}
}

func TestAudioRenderer_RenderWav_DebugMode(ts *testing.T) {
	var p0, pEnd t.Period
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   100,
		Resonance: 5,
		Amplitude: t.AmplitudePercentToRaw(10),
		Waveform:  t.WaveformSine,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]
	pEnd.Time = 500 // 0.5 seconds

	periods := []t.Period{p0, pEnd}

	options := &AudioRendererOptions{
		SampleRate: 44100,
		Volume:     100,
		GainLevel:  t.GainLevelMedium,
		Quiet:      true,
		Debug:      true, // Debug mode - no file should be created
	}

	renderer, err := NewAudioRenderer(periods, options)
	if err != nil {
		ts.Fatalf("NewAudioRenderer failed: %v", err)
	}

	tempDir := ts.TempDir()
	outPath := filepath.Join(tempDir, "debug_test.wav")

	// In debug mode, this should not create a file
	if err := renderer.RenderWav(outPath); err != nil {
		ts.Fatalf("RenderWav in debug mode failed: %v", err)
	}

	// File should not exist in debug mode
	if _, err := os.Stat(outPath); err == nil {
		ts.Fatalf("File should not be created in debug mode")
	}
}

func TestAudioRenderer_Render_CallbacksAndSizes(ts *testing.T) {
	sr := 44100

	endMs := 1234
	totalFrames := int64(math.Round(float64(endMs) * float64(sr) / 1000.0))

	var p0, pEnd t.Period
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackMonauralBeat,
		Carrier:   220,
		Resonance: 5,
		Amplitude: t.AmplitudePercentToRaw(20),
		Waveform:  t.WaveformSine,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]
	pEnd.Time = endMs

	periods := []t.Period{p0, pEnd}

	opts := &AudioRendererOptions{
		SampleRate: sr,
		Volume:     80,
		GainLevel:  t.GainLevelMedium,
		Quiet:      true,
		Debug:      false,
	}

	r, err := NewAudioRenderer(periods, opts)
	if err != nil {
		ts.Fatalf("NewAudioRenderer failed: %v", err)
	}

	var lens []int
	calls := 0
	var formatChecked bool

	consume := func(buf *audio.IntBuffer) error {
		if !formatChecked {
			if buf.Format == nil {
				ts.Fatalf("Buffer format is nil")
			}
			if buf.Format.SampleRate != sr {
				ts.Fatalf("SampleRate mismatch: got %d, want %d", buf.Format.SampleRate, sr)
			}
			if buf.Format.NumChannels != audioChannels {
				ts.Fatalf("NumChannels mismatch: got %d, want %d", buf.Format.NumChannels, audioChannels)
			}
			if buf.SourceBitDepth != audioBitDepth {
				ts.Fatalf("BitDepth mismatch: got %d, want %d", buf.SourceBitDepth, audioBitDepth)
			}
			formatChecked = true
		}
		lens = append(lens, len(buf.Data))
		calls++
		return nil
	}

	if err := r.Render(consume); err != nil {
		ts.Fatalf("Render failed: %v", err)
	}

	chunk := int64(t.BufferSize)
	full := int(totalFrames / chunk)
	rem := int(totalFrames % chunk)

	expected := make([]int, 0, full+1)
	for i := 0; i < full; i++ {
		expected = append(expected, t.BufferSize*audioChannels)
	}
	if rem > 0 {
		expected = append(expected, rem*audioChannels)
	}

	if calls != len(expected) {
		ts.Fatalf("Expected %d callbacks, got %d", len(expected), calls)
	}
	for i := range expected {
		if lens[i] != expected[i] {
			ts.Fatalf("Chunk %d size mismatch: got %d, want %d", i, lens[i], expected[i])
		}
	}
}

func TestAudioRenderer_Render_PropagatesError(ts *testing.T) {
	sr := 44100
	endMs := 2000

	var p0, pEnd t.Period
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   200,
		Resonance: 7,
		Amplitude: t.AmplitudePercentToRaw(10),
		Waveform:  t.WaveformSine,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]
	pEnd.Time = endMs
	periods := []t.Period{p0, pEnd}

	opts := &AudioRendererOptions{
		SampleRate: sr,
		Volume:     100,
		GainLevel:  t.GainLevelMedium,
		Quiet:      true,
		Debug:      false,
	}

	r, err := NewAudioRenderer(periods, opts)
	if err != nil {
		ts.Fatalf("NewAudioRenderer failed: %v", err)
	}

	targetErr := errors.New("sink failure")
	consume := func(buf *audio.IntBuffer) error {
		return targetErr
	}

	err = r.Render(consume)
	if err == nil {
		ts.Fatalf("Expected error from consumer, got nil")
	}
	if !errors.Is(err, targetErr) {
		ts.Fatalf("Expected wrapped target error, got: %v", err)
	}
}

func TestAudioRenderer_Render_NilConsumer(ts *testing.T) {
	sr := 44100
	endMs := 2000

	var p0, pEnd t.Period
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackIsochronicBeat,
		Carrier:   10,
		Resonance: 2,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformTriangle,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]
	pEnd.Time = endMs
	periods := []t.Period{p0, pEnd}

	opts := &AudioRendererOptions{
		SampleRate: sr,
		Volume:     90,
		GainLevel:  t.GainLevelMedium,
		Quiet:      true,
		Debug:      false,
	}

	r, err := NewAudioRenderer(periods, opts)
	if err != nil {
		ts.Fatalf("NewAudioRenderer failed: %v", err)
	}

	if err := r.Render(nil); err != nil {
		ts.Fatalf("Render with nil consumer failed: %v", err)
	}
}

func TestRenderRaw_WritesExpectedBytesAndRestoresQuiet(ts *testing.T) {
	sr := 44100
	endMs := 200

	var p0, pEnd t.Period
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackMonauralBeat,
		Carrier:   220,
		Resonance: 7,
		Amplitude: t.AmplitudePercentToRaw(20),
		Waveform:  t.WaveformSine,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]
	pEnd.Time = endMs

	periods := []t.Period{p0, pEnd}

	opts := &AudioRendererOptions{
		SampleRate: sr,
		Volume:     90,
		GainLevel:  t.GainLevelMedium,
		Quiet:      false,
		Debug:      false,
	}

	r, err := NewAudioRenderer(periods, opts)
	if err != nil {
		ts.Fatalf("NewAudioRenderer failed: %v", err)
	}

	var buf bytes.Buffer
	origQuiet := r.Quiet

	if err := r.RenderRaw(&buf); err != nil {
		ts.Fatalf("RenderRaw failed: %v", err)
	}

	if r.Quiet != origQuiet {
		ts.Fatalf("Quiet should be restored to %v, got %v", origQuiet, r.Quiet)
	}

	totalFrames := int64(math.Round(float64(endMs) * float64(sr) / 1000.0))
	expectedBytes := totalFrames * int64(audioChannels*3)
	got := int64(buf.Len())
	if got != expectedBytes {
		ts.Fatalf("byte length mismatch: got %d, want %d", got, expectedBytes)
	}

	raw := buf.Bytes()
	const maxSamplesToCheck = 2000
	foundNonZero := false
	limit := maxSamplesToCheck * 3
	if limit > len(raw) {
		limit = len(raw)
	}
	for i := 0; i+2 < limit; i += 3 {
		v := int32(raw[i]) | int32(raw[i+1])<<8 | int32(raw[i+2])<<16
		// Sign-extend 24-bit
		if v&0x00800000 != 0 {
			v |= ^int32(0x00FFFFFF)
		}
		if v != 0 {
			foundNonZero = true
			break
		}
	}
	if !foundNonZero {
		ts.Fatalf("no non-zero samples found in first %d samples", maxSamplesToCheck)
	}
}

type failingWriter struct{}

func (f *failingWriter) Write(p []byte) (int, error) {
	return 0, errors.New("sink failure")
}

func TestRenderRaw_PropagatesWriteError(ts *testing.T) {
	sr := 44100
	endMs := 50

	var p0, pEnd t.Period
	p0.Time = 0
	p0.TrackStart[0] = t.Track{
		Type:      t.TrackIsochronicBeat,
		Carrier:   10,
		Resonance: 2,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformTriangle,
	}
	p0.TrackEnd[0] = p0.TrackStart[0]
	pEnd.Time = endMs

	opts := &AudioRendererOptions{
		SampleRate: sr,
		Volume:     80,
		GainLevel:  t.GainLevelMedium,
		Quiet:      false,
		Debug:      false,
	}

	r, err := NewAudioRenderer([]t.Period{p0, pEnd}, opts)
	if err != nil {
		ts.Fatalf("NewAudioRenderer failed: %v", err)
	}

	err = r.RenderRaw(&failingWriter{})
	if err == nil {
		ts.Fatalf("expected error from writer, got nil")
	}
}
