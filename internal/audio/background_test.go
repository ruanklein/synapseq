/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package audio

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	bwav "github.com/gopxl/beep/v2/wav"
)

const maxBackgroundFileSize = 10 * 1024 * 1024 // 10MB

func mustReadWavAll(t *testing.T, path string) ([]int, uint32, int, int) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open wav: %v", err)
	}
	defer f.Close()

	s, fmt, err := bwav.Decode(f)
	if err != nil {
		t.Fatalf("decode wav: %v", err)
	}
	defer s.Close()

	// Stream all frames and convert to interleaved int16 samples to match BackgroundAudio
	var data []int
	const scale = 32768.0 // 2^15 for 16-bit
	buf := make([][2]float64, 4096)
	for {
		n, ok := s.Stream(buf)
		for i := 0; i < n; i++ {
			l := int(buf[i][0] * scale)
			r := int(buf[i][1] * scale)
			if l > audioMaxValue {
				l = audioMaxValue
			}
			if l < audioMinValue {
				l = audioMinValue
			}
			if r > audioMaxValue {
				r = audioMaxValue
			}
			if r < audioMinValue {
				r = audioMinValue
			}
			data = append(data, l, r)
		}
		if !ok {
			break
		}
	}
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	// Return data, sample rate, channels, bit depth
	return data, uint32(fmt.SampleRate), fmt.NumChannels, fmt.Precision * 8
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

func TestBackgroundAudio_RemoteWAV(t *testing.T) {
	// Read local WAV file to serve
	path := filepath.Join("testdata", "noise.wav")
	wavData, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test WAV: %v", err)
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/wav")
		w.Write(wavData)
	}))
	defer server.Close()

	// Load from remote URL
	bg, err := NewBackgroundAudio(server.URL)
	if err != nil {
		t.Fatalf("NewBackgroundAudio remote: %v", err)
	}
	defer bg.Close()

	if !bg.IsEnabled() {
		t.Fatalf("expected enabled background for remote")
	}

	// Verify cache was populated
	if bg.cachedData == nil {
		t.Fatalf("expected cachedData to be populated")
	}
	if len(bg.cachedData) != len(wavData) {
		t.Fatalf("cached data size mismatch: got %d want %d", len(bg.cachedData), len(wavData))
	}

	// Read some samples to verify it works
	buf := make([]int, 1024)
	n, err := bg.ReadSamples(buf, len(buf))
	if err != nil {
		t.Fatalf("ReadSamples remote error: %v", err)
	}
	if n != len(buf) {
		t.Fatalf("ReadSamples remote count: got %d want %d", n, len(buf))
	}

	// Verify non-zero samples
	hasNonZero := false
	for _, v := range buf {
		if v != 0 {
			hasNonZero = true
			break
		}
	}
	if !hasNonZero {
		t.Fatalf("expected non-zero samples from remote WAV")
	}
}

func TestBackgroundAudio_Remote10MBLimit(ts *testing.T) {
	// Create a server that serves more than 10MB
	const size = 12 * 1024 * 1024 // 12MB
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/wav")
		// Write a simple WAV header (44 bytes) + data
		header := make([]byte, 44)
		copy(header[0:4], "RIFF")
		copy(header[8:12], "WAVE")
		copy(header[12:16], "fmt ")
		// fmt chunk size
		header[16] = 16
		// PCM format (1)
		header[20] = 1
		// 2 channels
		header[22] = 2
		// 44100 sample rate
		header[24] = 0x44
		header[25] = 0xac
		// byte rate
		header[28] = 0x10
		header[29] = 0xb1
		header[30] = 0x02
		// block align
		header[32] = 4
		// bits per sample
		header[34] = 16
		// data chunk
		copy(header[36:40], "data")
		// data size (size - 44)
		dataSize := size - 44
		header[40] = byte(dataSize)
		header[41] = byte(dataSize >> 8)
		header[42] = byte(dataSize >> 16)
		header[43] = byte(dataSize >> 24)

		w.Write(header)
		// Write more data to exceed 10MB
		chunk := make([]byte, 1024*1024) // 1MB chunks
		for i := 0; i < size-44; i += len(chunk) {
			remaining := size - 44 - i
			if remaining < len(chunk) {
				w.Write(chunk[:remaining])
			} else {
				w.Write(chunk)
			}
		}
	}))
	defer server.Close()

	bg, err := NewBackgroundAudio(server.URL)
	if err != nil {
		ts.Fatalf("NewBackgroundAudio 10MB limit: %v", err)
	}
	defer bg.Close()

	// Verify that only 10MB was read
	if len(bg.cachedData) != maxBackgroundFileSize {
		ts.Fatalf("expected cached data to be limited to %d bytes, got %d", maxBackgroundFileSize, len(bg.cachedData))
	}
}

func TestBackgroundAudio_Local10MBLimit(ts *testing.T) {
	// Create a temporary WAV file larger than 10MB
	tmpDir := ts.TempDir()
	path := filepath.Join(tmpDir, "large.wav")

	f, err := os.Create(path)
	if err != nil {
		ts.Fatalf("failed to create temp file: %v", err)
	}

	const size = 12 * 1024 * 1024 // 12MB
	// Write WAV header
	header := make([]byte, 44)
	copy(header[0:4], "RIFF")
	copy(header[8:12], "WAVE")
	copy(header[12:16], "fmt ")
	header[16] = 16
	header[20] = 1
	header[22] = 2
	header[24] = 0x44
	header[25] = 0xac
	header[28] = 0x10
	header[29] = 0xb1
	header[30] = 0x02
	header[32] = 4
	header[34] = 16
	copy(header[36:40], "data")
	dataSize := size - 44
	header[40] = byte(dataSize)
	header[41] = byte(dataSize >> 8)
	header[42] = byte(dataSize >> 16)
	header[43] = byte(dataSize >> 24)

	if _, err := f.Write(header); err != nil {
		ts.Fatalf("failed to write header: %v", err)
	}

	// Write data
	chunk := make([]byte, 1024*1024) // 1MB chunks
	for i := 0; i < size-44; i += len(chunk) {
		remaining := size - 44 - i
		if remaining < len(chunk) {
			if _, err := f.Write(chunk[:remaining]); err != nil {
				ts.Fatalf("failed to write data: %v", err)
			}
		} else {
			if _, err := f.Write(chunk); err != nil {
				ts.Fatalf("failed to write data: %v", err)
			}
		}
	}
	f.Close()

	bg, err := NewBackgroundAudio(path)
	if err != nil {
		ts.Fatalf("NewBackgroundAudio local 10MB limit: %v", err)
	}
	defer bg.Close()

	// Verify that only 10MB was read
	if len(bg.cachedData) != maxBackgroundFileSize {
		ts.Fatalf("expected cached data to be limited to %d bytes, got %d", maxBackgroundFileSize, len(bg.cachedData))
	}
}

func TestBackgroundAudio_InvalidPath(t *testing.T) {
	if _, err := NewBackgroundAudio(filepath.Join("testdata", "missing.wav")); err == nil {
		t.Fatalf("expected error for missing background file")
	}
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
