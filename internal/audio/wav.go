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

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

// RenderToWAV renders the audio to a WAV file using go-audio/wav
func (r *AudioRenderer) RenderToWAV(outPath string) error {
	// If in debug mode, render to stdout (nil writer)
	if r.Debug {
		return r.Render(nil)
	}

	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	enc := wav.NewEncoder(out, r.SampleRate, audioBitDepth, audioChannels, 1)
	defer enc.Close()

	return r.Render(func(buf *audio.IntBuffer) error {
		return enc.Write(buf)
	})
}
