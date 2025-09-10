/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package audio

import (
	"math"

	t "github.com/ruanklein/synapseq/internal/types"
)

// InitWaveformTables initializes the waveform tables
func InitWaveformTables() [4][]int {
	var waveTables [4][]int
	for i := range waveTables {
		waveformTable := make([]int, t.SineTableSize)

		for j := range t.SineTableSize {
			phase := float64(j) * 2.0 * float64(math.Pi) / float64(t.SineTableSize)
			var val float64

			switch i { // i is the waveform type
			case int(t.WaveformSine):
				val = math.Sin(phase)
			case int(t.WaveformSquare):
				if math.Sin(phase) > 0 {
					val = 1.0
				} else {
					val = -1.0
				}
			case int(t.WaveformTriangle):
				if phase < math.Pi {
					val = (2.0 * phase / math.Pi) - 1.0
				} else {
					val = 3.0 - (2.0 * phase / math.Pi)
				}
			case int(t.WaveformSawtooth):
				val = (2.0 * phase / (2.0 * math.Pi)) - 1.0
			default:
				val = math.Sin(phase)
			}

			waveformTable[j] = int(t.WaveTableAmplitude * val)
		}

		waveTables[i] = waveformTable
	}
	return waveTables
}
