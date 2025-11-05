/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package audio

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gopxl/beep/v2"
	bwav "github.com/gopxl/beep/v2/wav"

	s "github.com/ruanklein/synapseq/v3/internal/shared"
	t "github.com/ruanklein/synapseq/v3/internal/types"
)

// Max background file size: 10MB
const maxBackgroundFileSize = 10 * 1024 * 1024

// BackgroundAudio handles background WAV file playback with looping
type BackgroundAudio struct {
	filePath      string
	decoder       beep.StreamSeekCloser
	sampleRate    int
	channels      int
	bitDepth      int
	isEnabled     bool
	hasReachedEOF bool

	// Cache of the loaded file data
	cachedData []byte

	// Buffer for reading samples
	buffer     []int
	bufferSize int
}

// NewBackgroundAudio creates a new background audio processor
func NewBackgroundAudio(filePath string) (*BackgroundAudio, error) {
	if filePath == "" {
		return &BackgroundAudio{isEnabled: false}, nil
	}

	bg := &BackgroundAudio{
		filePath:   filePath,
		bufferSize: t.BufferSize * audioChannels, // Stereo
		isEnabled:  true,
	}

	if err := bg.loadAndCache(); err != nil {
		return nil, err
	}

	if err := bg.openFromCache(); err != nil {
		return nil, fmt.Errorf("failed to open background file: %w", err)
	}

	bg.buffer = make([]int, bg.bufferSize)
	return bg, nil
}

// allowedWavContentType checks Content-Type for WAV
func allowedWavContentType(ct string) bool {
	ct = strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
	switch ct {
	case "audio/wav", "audio/x-wav", "audio/wave":
		return true
	default:
		return false
	}
}

// loadAndCache loads the file (local or remote) into memory cache
func (bg *BackgroundAudio) loadAndCache() error {
	if bg.cachedData != nil {
		return nil
	}

	data, err := s.GetFile(bg.filePath, t.FormatWAV)
	if err != nil {
		return fmt.Errorf("failed to load background file: %w", err)
	}

	// Read all data into cache
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, data)
	if err != nil {
		return fmt.Errorf("failed to read background file data: %w", err)
	}

	bg.cachedData = buf.Bytes()
	return nil
}

// openFromCache opens a decoder from the cached data
func (bg *BackgroundAudio) openFromCache() error {
	if bg.cachedData == nil {
		return fmt.Errorf("no cached data available")
	}

	reader := bytes.NewReader(bg.cachedData)
	s, f, err := bwav.Decode(reader)
	if err != nil {
		return fmt.Errorf("invalid WAV file: %s: %w", bg.filePath, err)
	}

	bg.decoder = s
	bg.sampleRate = int(f.SampleRate)
	bg.channels = f.NumChannels
	bg.bitDepth = f.Precision * 8
	bg.hasReachedEOF = false

	return nil
}

// restart reopens the decoder from cache for looping
func (bg *BackgroundAudio) restart() error {
	if !bg.isEnabled {
		return nil
	}

	// Close current decoder
	if bg.decoder != nil {
		_ = bg.decoder.Close()
		bg.decoder = nil
	}

	// Reopen from cache
	if err := bg.openFromCache(); err != nil {
		return fmt.Errorf("failed to restart background file: %w", err)
	}

	return nil
}

// ReadSamples reads background audio samples with automatic looping
func (bg *BackgroundAudio) ReadSamples(samples []int, numSamples int) (int, error) {
	if !bg.isEnabled || bg.decoder == nil {
		// Fill with silence if no background
		for i := 0; i < numSamples; i++ {
			samples[i] = 0
		}
		return numSamples, nil
	}

	samplesRead := 0

	for samplesRead < numSamples {
		remaining := numSamples - samplesRead
		bufferOffset := samplesRead

		// Try to read from current position. If remaining is smaller than a frame,
		// read at least one full frame and then copy the needed tail.
		var n int
		var err error
		if remaining < bg.channels {
			tmp := make([]int, bg.channels)
			n, err = bg.readFromDecoder(tmp, bg.channels)
			if n > 0 {
				// copy only what's requested
				copy(samples[bufferOffset:bufferOffset+remaining], tmp[:remaining])
				n = remaining
			}
		} else {
			n, err = bg.readFromDecoder(samples[bufferOffset:bufferOffset+remaining], remaining)
		}
		samplesRead += n

		if err == io.EOF {
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
	if bg.decoder == nil {
		return 0, io.EOF
	}

	// Calculate how many frames to read (a frame is a set of samples for all channels)
	framesToRead := maxSamples / bg.channels
	if framesToRead <= 0 {
		// Need at least one frame to progress
		framesToRead = 1
	}

	buf := make([][2]float64, framesToRead)
	nFrames, ok := bg.decoder.Stream(buf)
	if !ok || nFrames == 0 {
		if err := bg.decoder.Err(); err != nil {
			return 0, err
		}
		return 0, io.EOF
	}

	const scale = 8388608.0 // 2^23
	outN := nFrames * bg.channels
	// Limit to maxSamples to avoid overrun when we read a full frame but caller
	// requested fewer samples than a full frame.
	if outN > maxSamples {
		outN = maxSamples
	}
	framesOut := outN / bg.channels
	for i := 0; i < framesOut; i++ {
		l := int(buf[i][0] * scale)
		r := int(buf[i][1] * scale)

		// clip to valid range
		if l > audioMaxValue {
			l = audioMaxValue
		}
		if l < audioMinValue {
			l = audioMinValue
		}
		if r > audioMaxValue {
			r = audioMaxValue
		}
		if r < audioMinValue {
			r = audioMinValue
		}

		samples[2*i] = l
		if bg.channels >= 2 {
			// only write second sample if we still have space
			if 2*i+1 < outN {
				samples[2*i+1] = r
			}
		} else {
			if 2*i+1 < outN {
				samples[2*i+1] = l
			}
		}
	}

	return outN, nil
}

// Close closes the background audio decoder
func (bg *BackgroundAudio) Close() error {
	bg.isEnabled = false
	if bg.decoder != nil {
		return bg.decoder.Close()
	}
	return nil
}

// IsEnabled returns whether background audio is enabled
func (bg *BackgroundAudio) IsEnabled() bool {
	return bg.isEnabled
}
