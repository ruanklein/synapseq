/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package audio

import (
	"fmt"
	"os"

	"github.com/gopxl/beep/v2"
	bwav "github.com/gopxl/beep/v2/wav"
)

// RenderWav renders the audio to a WAV file using go-audio/wav
func (r *AudioRenderer) RenderWav(outPath string) error {
	// If in debug mode, render to stdout (nil writer)
	if r.Debug {
		return r.Render(nil)
	}

	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	streamer := newRendererStreamer(r)
	format := beep.Format{
		SampleRate:  beep.SampleRate(r.SampleRate),
		NumChannels: audioChannels,
		Precision:   audioBitDepth / 8, // 24-bit -> 3 bytes
	}

	if err := bwav.Encode(out, streamer, format); err != nil {
		return err
	}
	if streamer.err != nil {
		return streamer.err
	}
	return nil
}
