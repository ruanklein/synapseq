//go:build wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package main

import (
	"syscall/js"

	"github.com/ruanklein/synapseq/v3/internal/audio"
	"github.com/ruanklein/synapseq/v3/internal/sequence"
)

// generateWav(spsqUint8Array) -> { wav: Uint8Array } or { error: "msg" }
func generateWav(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		return map[string]interface{}{
			"error": "missing SPSQ input buffer",
		}
	}

	spsqJS := args[0]
	raw := make([]byte, spsqJS.Length())
	js.CopyBytesToGo(raw, spsqJS)

	seq, err := sequence.LoadTextSequence(raw)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	renderer, err := audio.NewAudioRenderer(seq.Periods, &audio.AudioRendererOptions{
		SampleRate:     seq.Options.SampleRate,
		Volume:         seq.Options.Volume,
		GainLevel:      seq.Options.GainLevel,
		BackgroundPath: seq.Options.BackgroundPath,
	})
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	wavBytes, err := renderer.RenderWav()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	uint8Array := js.Global().Get("Uint8Array").New(len(wavBytes))
	js.CopyBytesToJS(uint8Array, wavBytes)

	return map[string]interface{}{
		"wav": uint8Array,
		"ok":  true,
	}
}

func main() {
	js.Global().Set("synapseqGenerate", js.FuncOf(generateWav))

	select {}
}
