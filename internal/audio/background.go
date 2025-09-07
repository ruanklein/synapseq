package audio

import (
	"fmt"
	"io"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	t "github.com/ruanklein/synapseq/internal/types"
)

// BackgroundAudio handles background WAV file playback with looping
type BackgroundAudio struct {
	filePath      string
	file          *os.File
	decoder       *wav.Decoder
	sampleRate    int
	channels      int
	bitDepth      int
	isEnabled     bool
	hasReachedEOF bool
	// Buffer for reading samples
	buffer     []int
	bufferSize int
}

// NewBackgroundAudio creates a new background audio processor
func NewBackgroundAudio(filePath string, sampleRate int) (*BackgroundAudio, error) {
	if filePath == "" {
		return &BackgroundAudio{isEnabled: false}, nil
	}

	bg := &BackgroundAudio{
		filePath:   filePath,
		sampleRate: sampleRate,
		bufferSize: t.BufferSize * 2, // Stereo
		isEnabled:  true,
	}

	if err := bg.openFile(); err != nil {
		return nil, fmt.Errorf("failed to open background file: %w", err)
	}

	bg.buffer = make([]int, bg.bufferSize)
	return bg, nil
}

// openFile opens and initializes the WAV file
func (bg *BackgroundAudio) openFile() error {
	var err error
	bg.file, err = os.Open(bg.filePath)
	if err != nil {
		return fmt.Errorf("cannot open file %s: %w", bg.filePath, err)
	}

	bg.decoder = wav.NewDecoder(bg.file)
	if !bg.decoder.IsValidFile() {
		bg.file.Close()
		return fmt.Errorf("invalid WAV file: %s", bg.filePath)
	}

	bg.sampleRate = int(bg.decoder.SampleRate)
	bg.channels = int(bg.decoder.NumChans)
	bg.bitDepth = int(bg.decoder.BitDepth)
	bg.hasReachedEOF = false

	// TODO: Support more formats
	if bg.channels != 2 {
		return fmt.Errorf("unsupported channel count: %d", bg.channels)
	}
	if bg.bitDepth != 24 {
		return fmt.Errorf("unsupported bit depth: %d", bg.bitDepth)
	}

	return nil
}

// restart reopens the file for looping
func (bg *BackgroundAudio) restart() error {
	if !bg.isEnabled {
		return nil
	}

	// Close current file
	if bg.file != nil {
		bg.file.Close()
	}

	// Reopen file
	if err := bg.openFile(); err != nil {
		return fmt.Errorf("failed to restart background file: %w", err)
	}

	return nil
}

// ReadSamples reads background audio samples with automatic looping
func (bg *BackgroundAudio) ReadSamples(samples []int, numSamples int) (int, error) {
	if !bg.isEnabled || bg.decoder == nil {
		// Fill with silence if no background
		for i := range numSamples {
			samples[i] = 0
		}
		return numSamples, nil
	}

	samplesRead := 0

	for samplesRead < numSamples {
		remaining := numSamples - samplesRead
		bufferOffset := samplesRead

		// Try to read from current position
		n, err := bg.readFromDecoder(samples[bufferOffset:bufferOffset+remaining], remaining)
		samplesRead += n

		if err == io.EOF || n == 0 {
			// End of file reached, restart for looping
			if restartErr := bg.restart(); restartErr != nil {
				// If restart fails, fill remaining with silence
				for i := samplesRead; i < numSamples; i++ {
					samples[i] = 0
				}
				return numSamples, nil
			}
			// Continue reading after restart
			continue
		} else if err != nil {
			return samplesRead, fmt.Errorf("error reading background audio: %w", err)
		}

		// If we read less than requested but no error, we're at EOF
		if n < remaining {
			if restartErr := bg.restart(); restartErr != nil {
				// Fill remaining with silence if restart fails
				for i := samplesRead; i < numSamples; i++ {
					samples[i] = 0
				}
				return numSamples, nil
			}
		}
	}

	return samplesRead, nil
}

// readFromDecoder reads raw samples from the WAV decoder
func (bg *BackgroundAudio) readFromDecoder(samples []int, maxSamples int) (int, error) {
	// Create audio buffer for go-audio/audio
	audioBuf := &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: bg.channels,
			SampleRate:  bg.sampleRate,
		},
		Data:           bg.buffer[:maxSamples],
		SourceBitDepth: bg.bitDepth,
	}

	n, err := bg.decoder.PCMBuffer(audioBuf)
	if err != nil {
		return 0, err
	}

	// Copy the read samples
	copy(samples, audioBuf.Data[:n])

	return n, nil
}

// Close closes the background audio file
func (bg *BackgroundAudio) Close() error {
	bg.isEnabled = false
	if bg.file != nil {
		return bg.file.Close()
	}
	return nil
}

// IsEnabled returns whether background audio is enabled
func (bg *BackgroundAudio) IsEnabled() bool {
	return bg.isEnabled
}

// GetSampleRate returns the sample rate of the background audio
func (bg *BackgroundAudio) GetSampleRate() int {
	return bg.sampleRate
}

// GetChannels returns the number of channels in the background audio
func (bg *BackgroundAudio) GetChannels() int {
	return bg.channels
}

// CalculateBackgroundGain calculates the gain factor based on GainLevel
func calculateBackgroundGain(gainLevel t.GainLevel) float64 {
	// Convert dB reduction to linear gain factor
	// gain = 10^(-dB/20)
	switch gainLevel {
	case t.GainLevelVeryLow: // -20dB
		return 0.1
	case t.GainLevelLow: // -16dB
		return 0.158
	case t.GainLevelMedium: // -12dB
		return 0.25
	case t.GainLevelHigh: // -6dB
		return 0.5
	case t.GainLevelVeryHigh: // 0dB
		return 1.0
	default:
		return 0.25 // Default to medium
	}
}
