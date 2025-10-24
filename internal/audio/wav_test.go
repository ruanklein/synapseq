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
	"strings"
	"testing"

	"github.com/gopxl/beep/v2"
	bwav "github.com/gopxl/beep/v2/wav"
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

	// Write the ICMT chunk with the sequence
	if err := WriteICMTChunkFromTextFile(wavPath, seqPath); err != nil {
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
