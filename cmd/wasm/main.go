//go:build wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package main

import (
	"fmt"
	"syscall/js"

	"github.com/ruanklein/synapseq/v3/internal/audio"
	"github.com/ruanklein/synapseq/v3/internal/sequence"
)

// generateWav(spsqUint8Array) -> { wav: Uint8Array } or { error: "msg" }
func generateWav(this js.Value, args []js.Value) interface{} {
	promise := js.Global().Get("Promise").New(js.FuncOf(
		func(_ js.Value, pArgs []js.Value) interface{} {

			resolve := pArgs[0]
			reject := pArgs[1]

			go func() {
				defer func() {
					if r := recover(); r != nil {
						reject.Invoke(fmt.Sprintf("panic: %v", r))
					}
				}()

				if len(args) == 0 {
					reject.Invoke("missing SPSQ input buffer")
					return
				}

				input := args[0]
				raw := make([]byte, input.Length())
				js.CopyBytesToGo(raw, input)

				seq, err := sequence.LoadTextSequence(raw)
				if err != nil {
					reject.Invoke(err.Error())
					return
				}

				renderer, err := audio.NewAudioRenderer(seq.Periods, &audio.AudioRendererOptions{
					SampleRate:     seq.Options.SampleRate,
					Volume:         seq.Options.Volume,
					GainLevel:      seq.Options.GainLevel,
					BackgroundPath: seq.Options.BackgroundPath,
				})
				if err != nil {
					reject.Invoke(err.Error())
					return
				}

				wavBytes, err := renderer.RenderWav()
				if err != nil {
					reject.Invoke(err.Error())
					return
				}

				uint8Array := js.Global().Get("Uint8Array").New(len(wavBytes))
				js.CopyBytesToJS(uint8Array, wavBytes)

				result := js.Global().Get("Object").New()
				result.Set("wav", uint8Array)

				resolve.Invoke(result)
			}()

			return nil
		},
	))

	return promise
}

func main() {
	js.Global().Set("synapseqGenerate", js.FuncOf(generateWav))

	select {}
}
