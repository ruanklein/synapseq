package audio

import (
	"fmt"
	"math"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	t "github.com/ruanklein/synapseq/internal/types"
)

// AudioRenderer handle audio generation
type AudioRenderer struct {
	channels    [t.NumberOfChannels]t.Channel
	periods     []t.Period
	currentTime int      // Current playback time in milliseconds
	periodIdx   int      // Current period index
	waveTables  [4][]int // Waveform tables for different waveforms
	sampleRate  int
	volume      int // Volume level (0-100)
}

// NewAudioRenderer creates a new AudioRenderer instance
func NewAudioRenderer(periods []t.Period, option *t.Option) (*AudioRenderer, error) {
	if option == nil {
		return nil, fmt.Errorf("audio renderer options are required")
	}

	renderer := &AudioRenderer{
		periods:    periods,
		waveTables: InitWaveformTables(),
		sampleRate: option.SampleRate,
		volume:     option.Volume,
	}

	return renderer, nil
}

// UpdateChannelStates updates the state of each audio channel based on the current playback time
func (r *AudioRenderer) UpdateChannelStates(timeMs int) {
	r.currentTime = timeMs

	// Find the correct period for the current time
	for r.periodIdx+1 < len(r.periods) && timeMs >= r.periods[r.periodIdx+1].Time {
		r.periodIdx++
	}

	if r.periodIdx >= len(r.periods) {
		return
	}

	period := r.periods[r.periodIdx]
	nextTime := timeMs + 1000 // Default next time
	if r.periodIdx+1 < len(r.periods) {
		nextTime = r.periods[r.periodIdx+1].Time
	}

	// Calculate interpolation factor (0.0 to 1.0)
	t0 := period.Time
	t1 := nextTime
	progress := 0.0
	if t1 > t0 {
		progress = float64(timeMs-t0) / float64(t1-t0)
	}
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	// Update each channel
	for ch := range t.NumberOfChannels {
		r.updateSingleChannel(ch, period, progress)
	}
}

// GenerateAudioChunk generates a buffer of audio samples
func (r *AudioRenderer) GenerateAudioChunk() *audio.IntBuffer {
	samples := make([]int, t.BufferSize*2) // Stereo: left + right

	for i := range t.BufferSize {
		left, right := r.generateStereoSample()

		// Apply global volume
		if r.volume != 100 {
			left = int(int64(left) * int64(r.volume) / 100)
			right = int(int64(right) * int64(r.volume) / 100)
		}

		// Clipping to 16-bit range
		if left > 32767 {
			left = 32767
		}
		if left < -32768 {
			left = -32768
		}
		if right > 32767 {
			right = 32767
		}
		if right < -32768 {
			right = -32768
		}

		samples[i*2] = left    // Left channel
		samples[i*2+1] = right // Right channel
	}

	return &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 2,
			SampleRate:  r.sampleRate,
		},
		Data:           samples,
		SourceBitDepth: 16,
	}
}

// generateStereoSample generates a single stereo sample
func (r *AudioRenderer) generateStereoSample() (int, int) {
	left, right := 0, 0

	for ch := range t.NumberOfChannels {
		channel := &r.channels[ch]
		l, r := r.generateChannelSample(channel)
		left += l
		right += r
	}

	return left >> 16, right >> 16 // Scale down to prevent overflow
}

// generateChannelSample generates sample for a specific channel
func (r *AudioRenderer) generateChannelSample(ch *t.Channel) (int, int) {
	switch ch.Voice.Type {
	case t.VoiceBinauralBeat:
		return r.generateBinauralSample(ch)
	default:
		return 0, 0 // Silence for unimplemented types
	}
}

// generateBinauralSample generates binaural beat sample
func (r *AudioRenderer) generateBinauralSample(ch *t.Channel) (int, int) {
	// Advance offset for each ear
	ch.Offset[0] += ch.Increment[0]
	ch.Offset[0] &= (t.SineTableSize << 16) - 1

	ch.Offset[1] += ch.Increment[1]
	ch.Offset[1] &= (t.SineTableSize << 16) - 1

	// Generate samples using waveform table
	waveIdx := int(ch.Voice.Waveform) % 4
	if waveIdx >= len(r.waveTables) {
		waveIdx = 0 // Default to sine wave
	}

	leftSample := ch.Amplitude[0] * r.waveTables[waveIdx][ch.Offset[0]>>16]
	rightSample := ch.Amplitude[1] * r.waveTables[waveIdx][ch.Offset[1]>>16]

	return leftSample, rightSample
}

// updateSingleChannel updates the state of a single audio channel
func (r *AudioRenderer) updateSingleChannel(chIdx int, period t.Period, progress float64) {
	ch := &r.channels[chIdx]
	v0 := period.VoiceStart[chIdx]
	v1 := period.VoiceEnd[chIdx]

	// Interpolate voice parameters
	ch.Voice = r.interpolateVoice(v0, v1, progress)

	// If the type changed, reset offsets
	if ch.Type != ch.Voice.Type {
		ch.Type = ch.Voice.Type
		ch.Offset[0] = 0
		ch.Offset[1] = 0
	}

	// Configure specific channel parameters based on type
	switch ch.Voice.Type {
	case t.VoiceOff: // Type 0
		ch.Amplitude[0] = 0
		ch.Amplitude[1] = 0
	case t.VoiceBinauralBeat: // Type 1
		freq1 := ch.Voice.Carrier + ch.Voice.Resonance/2
		freq2 := ch.Voice.Carrier - ch.Voice.Resonance/2
		ch.Amplitude[0] = int(ch.Voice.Amplitude)
		ch.Amplitude[1] = int(ch.Voice.Amplitude)
		ch.Increment[0] = int(freq1 / float64(r.sampleRate) * t.SineTableSize * 65536)
		ch.Increment[1] = int(freq2 / float64(r.sampleRate) * t.SineTableSize * 65536)

		// case t.VoiceWhiteNoise, t.VoicePinkNoise, t.VoiceBrownNoise: // Types 2, 9, 10
		// 	ch.Amplitude[0] = int(ch.Voice.Amplitude)

		// case t.VoiceMonauralBeat: // Type 3
		// 	freqHigh := ch.Voice.Carrier + ch.Voice.Resonance/2
		// 	freqLow := ch.Voice.Carrier - ch.Voice.Resonance/2
		// 	ch.Amplitude[0] = int(ch.Voice.Amplitude)
		// 	ch.Increment[0] = int(freqHigh / float64(r.sampleRate) * t.SineTableSize * 65536)
		// 	ch.Increment[1] = int(freqLow / float64(r.sampleRate) * t.SineTableSize * 65536)

		// case t.SpinPink: // Type 4
		// 	ch.Amplitude[0] = int(ch.Voice.Amplitude)
		// 	ch.Increment[0] = int(ch.Voice.Resonance / float64(r.sampleRate) * ST_SIZ * 65536)
		// 	ch.Increment[1] = int(ch.Voice.Carrier * 1e-6 * float64(r.sampleRate) * (1 << 24) / ST_AMP)

		// case t.Background: // Type 5
		// 	ch.Amplitude[0] = int(ch.Voice.Amplitude)

		// case t.Isochronic: // Type 8
		// 	ch.Amplitude[0] = int(ch.Voice.Amplitude)
		// 	ch.Increment[0] = int(ch.Voice.Carrier / float64(r.sampleRate) * ST_SIZ * 65536)
		// 	ch.Increment[1] = int(ch.Voice.Resonance / float64(r.sampleRate) * ST_SIZ * 65536)
		//
	default:
		panic("voice not implemented")
	}
}

func (r *AudioRenderer) interpolateVoice(v0, v1 t.Voice, progress float64) t.Voice {
	return t.Voice{
		Type:      v0.Type, // Type does not interpolate
		Amplitude: t.AmplitudeType(float64(v0.Amplitude)*(1-progress) + float64(v1.Amplitude)*progress),
		Carrier:   v0.Carrier*(1-progress) + v1.Carrier*progress,
		Resonance: v0.Resonance*(1-progress) + v1.Resonance*progress,
		Waveform:  v0.Waveform, // Waveform does not interpolate
		Intensity: t.IntensityType(float64(v0.Intensity)*(1-progress) + float64(v1.Intensity)*progress),
	}
}

// RenderToWAV renders the audio to a WAV file using go-audio/wav
func (r *AudioRenderer) RenderToWAV(outPath string) error {
	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	enc := wav.NewEncoder(out, r.sampleRate, 16, 2, 1)

	endMs := r.periods[len(r.periods)-1].Time

	totalFrames := int64(math.Round(float64(endMs) * float64(r.sampleRate) / 1000.0))
	framesWritten := int64(0)
	chunkFrames := int64(t.BufferSize)

	r.currentTime = 0
	r.periodIdx = 0
	for i := 0; i < t.NumberOfChannels; i++ {
		r.channels[i] = t.Channel{}
	}

	for framesWritten < totalFrames {
		currentTimeMs := int((float64(framesWritten) * 1000.0) / float64(r.sampleRate))
		r.UpdateChannelStates(currentTimeMs)

		buf := r.GenerateAudioChunk()

		framesToWrite := chunkFrames
		if remain := totalFrames - framesWritten; remain < chunkFrames {
			framesToWrite = remain
			buf.Data = buf.Data[:remain*2] // stereo interleaved
		}

		if err := enc.Write(buf); err != nil {
			_ = enc.Close()
			return fmt.Errorf("write wav: %w", err)
		}

		framesWritten += framesToWrite

		if framesWritten%int64(r.sampleRate) == 0 {
			pct := float64(framesWritten) / float64(totalFrames) * 100.0
			secs := float64(framesWritten) / float64(r.sampleRate)
			fmt.Printf("Progress: %.1f%% (%d/%d frames, %.1fs)\n", pct, framesWritten, totalFrames, secs)
		}
	}

	if err := enc.Close(); err != nil {
		return fmt.Errorf("close encoder: %w", err)
	}
	if err := out.Sync(); err != nil {
		return fmt.Errorf("sync file: %w", err)
	}

	secs := float64(framesWritten) / float64(r.sampleRate)
	fmt.Printf("Audio rendering complete: %d frames written (%.2f seconds)\n", framesWritten, secs)
	return nil
}
