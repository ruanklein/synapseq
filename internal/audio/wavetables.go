package audio

import (
	"math"
)

const (
	WaveTableAmplitude = 0x7FFFF // Amplitude of wave in wave-table
)

// InitWaveformTables initializes the waveform tables
func InitWaveformTables() [4][]int {
	var waveTables [4][]int
	for i := range waveTables {
		waveformTable := make([]int, SineTableSize)

		for j := range SineTableSize {
			phase := float64(j) * 2.0 * float64(math.Pi) / float64(SineTableSize)
			var val float64

			switch i { // i is the waveform type
			case int(WaveformSine):
				val = math.Sin(phase)
			case int(WaveformSquare):
				if math.Sin(phase) > 0 {
					val = 1.0
				} else {
					val = -1.0
				}
			case int(WaveformTriangle):
				if phase < math.Pi {
					val = (2.0 * phase / math.Pi) - 1.0
				} else {
					val = 3.0 - (2.0 * phase / math.Pi)
				}
			case int(WaveformSawtooth):
				val = (2.0 * phase / (2.0 * math.Pi)) - 1.0
			default:
				val = math.Sin(phase)
			}

			waveformTable[j] = int(WaveTableAmplitude * val)
		}

		waveTables[i] = waveformTable
	}
	return waveTables
}
