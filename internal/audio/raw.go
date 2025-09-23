/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package audio

import (
	"bufio"
	"io"

	"github.com/go-audio/audio"
	t "github.com/ruanklein/synapseq/internal/types"
)

// RenderRaw renders the audio to a raw PCM stream (24-bit little-endian)
func (r *AudioRenderer) RenderRaw(w io.Writer) error {
	origQuiet := r.Quiet
	r.Quiet = true
	defer func() { r.Quiet = origQuiet }()

	bw := bufio.NewWriter(w)
	// 3 bytes per sample (24-bit)
	out := make([]byte, t.BufferSize*audioChannels*3)

	err := r.Render(func(buf *audio.IntBuffer) error {
		need := len(buf.Data) * 3
		if cap(out) < need {
			out = make([]byte, need)
		}
		b := out[:need]

		j := 0
		for _, s := range buf.Data {
			if s > audioMaxValue {
				s = audioMaxValue
			} else if s < audioMinValue {
				s = audioMinValue
			}
			v := int32(s)
			b[j] = byte(v)         // LSB
			b[j+1] = byte(v >> 8)  // Mid
			b[j+2] = byte(v >> 16) // MSB (signed)
			j += 3
		}
		_, err := bw.Write(b)
		return err
	})
	if err != nil {
		return err
	}
	return bw.Flush()
}
