package audio

const (
	NoiseShift     = 12                               // NoiseShift is the bit shift for noise generation
	NoiseDither    = 16                               // NoiseDither is the bit depth for noise dithering
	NoiseAmplitude = WaveTableAmplitude << NoiseShift // NoiseAmplitude is the amplitude for noise generation
	NoiseBands     = 9                                // NoiseBands is the number of bands for noise generation
)

// Noise represents a noise generator state
type Noise struct {
	Value     int // Current output value
	Increment int // Increment
}
