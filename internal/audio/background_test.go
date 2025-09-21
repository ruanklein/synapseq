/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package audio

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-audio/wav"
)

func mustReadWavAll(t *testing.T, path string) ([]int, uint32, int, int) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open wav: %v", err)
	}
	defer f.Close()
	dec := wav.NewDecoder(f)
	if !dec.IsValidFile() {
		t.Fatalf("invalid wav: %s", path)
	}
	buf, err := dec.FullPCMBuffer()
	if err != nil {
		t.Fatalf("read wav: %v", err)
	}
	return buf.Data, dec.SampleRate, int(dec.NumChans), int(dec.BitDepth)
}

func TestBackgroundAudio_LoadReadAndLoop(t *testing.T) {
	path := filepath.Join("testdata", "noise.wav")
	data, sr, chans, depth := mustReadWavAll(t, path)

	bg, err := NewBackgroundAudio(path)
	if err != nil {
		t.Fatalf("NewBackgroundAudio: %v", err)
	}
	defer bg.Close()

	if !bg.IsEnabled() {
		t.Fatalf("expected enabled background")
	}
	if bg.sampleRate != int(sr) || bg.channels != chans || bg.bitDepth != depth {
		t.Fatalf("mismatched bg props sr=%d ch=%d bd=%d vs file sr=%d ch=%d bd=%d", bg.sampleRate, bg.channels, bg.bitDepth, sr, chans, depth)
	}

	// Force looping at least once, reading in chunks up to bg.bufferSize
	target := len(data) + 123
	var buf []int
	chunk := bg.bufferSize
	if chunk <= 0 {
		t.Fatalf("invalid bg.bufferSize: %d", chunk)
	}
	tmp := make([]int, chunk)
	total := 0
	for total < target {
		need := target - total
		if need > chunk {
			need = chunk
		}
		n, err := bg.ReadSamples(tmp[:need], need)
		if err != nil {
			t.Fatalf("ReadSamples error: %v", err)
		}
		if n != need {
			t.Fatalf("ReadSamples short read: got %d want %d", n, need)
		}
		buf = append(buf, tmp[:need]...)
		total += n
	}

	// Prefix must match the original file data
	for i := 0; i < len(data) && i < len(buf); i++ {
		if buf[i] != data[i] {
			t.Fatalf("prefix mismatch at %d: got %d want %d", i, buf[i], data[i])
		}
	}

	// After exactly len(data) samples, sequence should restart at beginning
	if total > len(data) {
		if buf[len(data)] != data[0] {
			t.Fatalf("loop restart mismatch: got %d want %d", buf[len(data)], data[0])
		}
	}
}

func TestBackgroundAudio_DisabledAndClose(t *testing.T) {
	bg, err := NewBackgroundAudio("")
	if err != nil {
		t.Fatalf("NewBackgroundAudio empty: %v", err)
	}
	if bg.IsEnabled() {
		t.Fatalf("expected disabled when no path provided")
	}
	buf := make([]int, 256)
	n, err := bg.ReadSamples(buf, len(buf))
	if err != nil {
		t.Fatalf("ReadSamples disabled error: %v", err)
	}
	if n != len(buf) {
		t.Fatalf("ReadSamples disabled count: got %d want %d", n, len(buf))
	}
	for i, v := range buf {
		if v != 0 {
			t.Fatalf("disabled should fill zeros at %d: %d", i, v)
		}
	}

	// Closing should keep it disabled and safe
	if err := bg.Close(); err != nil {
		t.Fatalf("Close disabled: %v", err)
	}
	_, err = bg.ReadSamples(buf, len(buf))
	if err != nil {
		t.Fatalf("ReadSamples after close error: %v", err)
	}
	for i, v := range buf {
		if v != 0 {
			t.Fatalf("after close should fill zeros at %d: %d", i, v)
		}
	}
}

func TestBackgroundAudio_InvalidPath(t *testing.T) {
	if _, err := NewBackgroundAudio(filepath.Join("testdata", "missing.wav")); err == nil {
		t.Fatalf("expected error for missing background file")
	}
}
