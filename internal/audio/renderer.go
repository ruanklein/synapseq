package audio

import (
	"fmt"
	"math"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/ruanklein/synapseq/internal/generator"
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
func (r *AudioRenderer) GenerateAudioChunk(samples []int) []int {
	for i := range t.BufferSize {
		left, right := r.generateStereoSample()

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

		samples[i*2] = left
		samples[i*2+1] = right
	}

	return samples
}

// generateStereoSample generates a single stereo sample
func (r *AudioRenderer) generateStereoSample() (int, int) {
	left, right := 0, 0

	for ch := range t.NumberOfChannels {
		channel := &r.channels[ch]
		l, r := generator.GenerateSample(channel, r.waveTables)
		left += l
		right += r
	}

	if r.volume != 100 {
		left = int(int64(left) * int64(r.volume) / 100)
		right = int(int64(right) * int64(r.volume) / 100)
	}

	return left >> 16, right >> 16 // Scale down to prevent overflow
}

// updateSingleChannel updates the state of a single audio channel
func (r *AudioRenderer) updateSingleChannel(chIdx int, period t.Period, progress float64) {
	if chIdx >= len(r.channels) || chIdx >= len(period.VoiceStart) {
		return // Bounds protection
	}

	ch := &r.channels[chIdx]
	v0 := period.VoiceStart[chIdx]
	v1 := period.VoiceEnd[chIdx]

	ch.Voice.Type = v0.Type
	ch.Voice.Amplitude = t.AmplitudeType(float64(v0.Amplitude)*(1-progress) + float64(v1.Amplitude)*progress)
	ch.Voice.Carrier = v0.Carrier*(1-progress) + v1.Carrier*progress
	ch.Voice.Resonance = v0.Resonance*(1-progress) + v1.Resonance*progress
	ch.Voice.Waveform = v0.Waveform
	ch.Voice.Intensity = t.IntensityType(float64(v0.Intensity)*(1-progress) + float64(v1.Intensity)*progress)

	if ch.Type != ch.Voice.Type {
		ch.Type = ch.Voice.Type
		ch.Offset[0] = 0
		ch.Offset[1] = 0
	}

	generator.UpdateChannel(ch, r.sampleRate)
}

// RenderToWAV renders the audio to a WAV file using go-audio/wav
func (r *AudioRenderer) RenderToWAV(outPath string) error {
	out := os.Stdout // Use standard output as default

	if outPath != "-" {
		var err error
		out, err = os.Create(outPath)
		if err != nil {
			return fmt.Errorf("create output: %w", err)
		}
		defer out.Close()
	}

	enc := wav.NewEncoder(out, r.sampleRate, 16, 2, 1)
	defer enc.Close()

	endMs := r.periods[len(r.periods)-1].Time

	totalFrames := int64(math.Round(float64(endMs) * float64(r.sampleRate) / 1000.0))
	chunkFrames := int64(t.BufferSize)
	framesWritten := int64(0)

	r.currentTime = 0
	r.periodIdx = 0

	for i := range t.NumberOfChannels {
		r.channels[i] = t.Channel{}
	}

	statusReporter := NewStatusReporter(false)
	defer statusReporter.FinalStatus()

	samples := make([]int, t.BufferSize*2) // Stereo: left + right
	audioBuf := &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 2,
			SampleRate:  r.sampleRate,
		},
		Data:           samples,
		SourceBitDepth: 16,
	}

	for framesWritten < totalFrames {
		currentTimeMs := int((float64(framesWritten) * 1000.0) / float64(r.sampleRate))
		r.UpdateChannelStates(currentTimeMs)

		statusReporter.CheckPeriodChange(r)

		audioBuf.Data = r.GenerateAudioChunk(samples)

		framesToWrite := chunkFrames
		if remain := totalFrames - framesWritten; remain < chunkFrames {
			framesToWrite = remain
			audioBuf.Data = audioBuf.Data[:remain*2] // stereo interleaved
		}

		if err := enc.Write(audioBuf); err != nil {
			enc.Close()
			return fmt.Errorf("write wav: %w", err)
		}

		framesWritten += framesToWrite

		if statusReporter.ShouldUpdateStatus() {
			statusReporter.DisplayStatus(r, currentTimeMs)
		}
	}

	return nil
}
