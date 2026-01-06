//go:build !wasm

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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gopxl/beep/v2"
	bwav "github.com/gopxl/beep/v2/wav"
	"github.com/synapseq-foundation/synapseq/v3/internal/info"
	s "github.com/synapseq-foundation/synapseq/v3/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

// Use constStreamer from renderer_test.go

func TestWriteAndExtractICMTChunk_Integration(ts *testing.T) {
	tempDir := ts.TempDir()
	wavPath := filepath.Join(tempDir, "test.wav")
	seqPath := filepath.Join(tempDir, "seq.spsq")

	// Create a simple text sequence file
	seqContent := "# Export test\n@samplerate 44100\n@volume 80\nalpha\n  tone 300 binaural 10 amplitude 20\n00:00:00 alpha\n00:01:00 alpha\n"
	if err := os.WriteFile(seqPath, []byte(seqContent), 0o600); err != nil {
		ts.Fatalf("failed to create seq.spsq: %v", err)
	}

	// Create a simple WAV file
	format := beep.Format{SampleRate: 44100, NumChannels: 2, Precision: 3}
	cs := &constStreamer{framesLeft: 44100, val: 0.1}
	wavFile, err := os.Create(wavPath)
	if err != nil {
		ts.Fatalf("failed to create WAV: %v", err)
	}
	if err := bwav.Encode(wavFile, cs, format); err != nil {
		ts.Fatalf("failed to write WAV: %v", err)
	}
	wavFile.Close()

	rawData, err := s.GetFile(seqPath, t.FormatText)
	if err != nil {
		ts.Fatalf("GetFile() error: %v", err)
	}

	metadata, err := info.NewMetadata(rawData)
	if err != nil {
		ts.Fatalf("ReadWAVMetadata error: %v", err)
	}

	// Write the ICMT chunk with the sequence
	if err := WriteICMTChunkFromTextFile(wavPath, metadata); err != nil {
		ts.Fatalf("WriteICMTChunkFromTextFile error: %v", err)
	}

	// Extract the sequence from the ICMT chunk
	content, err := ExtractTextSequenceFromWAV(wavPath)
	if err != nil {
		ts.Fatalf("ExtractTextSequenceFromWAV error: %v", err)
	}

	// Validate that the extracted content contains expected parts
	if !strings.Contains(content, "Export test") {
		ts.Fatalf("Extracted content does not contain expected text: %q", content)
	}
	if !strings.Contains(content, "@samplerate 44100") {
		ts.Fatalf("Extracted content does not contain samplerate: %q", content)
	}
	if !strings.Contains(content, "alpha") {
		ts.Fatalf("Extracted content does not contain preset: %q", content)
	}
}
