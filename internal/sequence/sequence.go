package sequence

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ruanklein/synapseq/internal/audio"
)

const (
	keywordComment       = "#"          // Represents a comment
	keywordDoubleComment = "##"         // Represents a double comment
	keywordOption        = "@"          // Represents an option
	keywordWaveform      = "waveform"   // Represents a waveform
	keywordSine          = "sine"       // Represents a sine wave
	keywordSquare        = "square"     // Represents a square wave
	keywordTriangle      = "triangle"   // Represents a triangle wave
	keywordSawtooth      = "sawtooth"   // Represents a sawtooth wave
	keywordTone          = "tone"       // Represents a tone
	keywordBinaural      = "binaural"   // Represents a binaural tone
	keywordMonaural      = "monaural"   // Represents a monaural tone
	keywordIsochronic    = "isochronic" // Represents an isochronic tone
	keywordAmplitude     = "amplitude"  // Represents an amplitude
	keywordNoise         = "noise"      // Represents a noise
	keywordWhite         = "white"      // Represents a white noise
	keywordPink          = "pink"       // Represents a pink noise
	keywordBrown         = "brown"      // Represents a brown noise
	keywordSpin          = "spin"       // Represents a spin
	keywordWidth         = "width"      // Represents a width
	keywordRate          = "rate"       // Represents a rate
	keywordEffect        = "effect"     // Represents an effect
	keywordBackground    = "background" // Represents a background
	keywordPulse         = "pulse"      // Represents a pulse
	keywordIntensity     = "intensity"  // Represents an intensity
)

// LoadSequence loads a sequence from a file
func LoadSequence(fileName string) error {
	file, err := LoadFile(fileName)
	if err != nil {
		return fmt.Errorf("error loading sequence file: %v", err)
	}

	for file.NextLine() {
		line := file.CurrentLine
		lineNumber := file.CurrentLineNumber

		// Split the line into fields
		fields := strings.Fields(line)

		// Skip empty lines and comments
		if len(fields) == 0 || fields[0] == keywordComment {
			continue
		}

		// Print if double comment
		if fields[0] == keywordDoubleComment {
			fmt.Printf("> %s\n", strings.Join(fields[1:], " "))
			continue
		}

		// Preset definition
		if line[0] != ' ' && IsPreset(fields[0]) {
			presetName := strings.ToLower(fields[0])
			if presetName == BuiltinSilence {
				return fmt.Errorf("cannot load built-in preset '%s' at line %d", BuiltinSilence, lineNumber)
			}

			var preset Preset
			if len(PresetList) > 0 {
				for i := 0; i < len(PresetList); i++ {
					if PresetList[i].Name == presetName {
						return fmt.Errorf("preset '%s' already exists at line %d", presetName, lineNumber)
					}
				}

				PresetList[len(PresetList)-1].Next = &preset
			}

			// Create a new preset
			preset.Name = presetName
			preset.InitVoices()
			PresetList = append(PresetList, preset)

			continue
		}

		// Voice line
		if len(line) > 3 && line[0] == ' ' && line[1] == ' ' && line[2] != ' ' {
			if len(PresetList) == 0 {
				return fmt.Errorf("voice defined without a preset list at line %d: %s", lineNumber, line)
			}

			preset := &PresetList[len(PresetList)-1]
			var voice audio.Voice

			switch fields[0] {
			case keywordWaveform: // waveform is valid with tone, spin, and effect
				if len(fields) < 2 {
					return fmt.Errorf("invalid voice definition at line %d: %s", lineNumber, line)
				}

				waveform, err := ParseWaveformType(fields[1])
				if err != nil {
					return fmt.Errorf("invalid waveform type at line %d: %s", lineNumber, line)
				}

				// Tone line
				if len(fields) == 8 {
					if fields[2] != keywordTone {
						return fmt.Errorf("expected tone keyword at line %d: %s", lineNumber, line)
					}

					carrier, err := strconv.ParseFloat(fields[3], 64)
					if err != nil {
						return fmt.Errorf("invalid carrier frequency at line %d: %s", lineNumber, line)
					}
				}
			}
			continue
		}

		// Error if no valid syntax is found
		return fmt.Errorf("invalid syntax at line %d: %s", lineNumber, line)
	}

	// Debug presets
	for _, p := range PresetList {
		if p.HasNext() {
			fmt.Printf("Preset: %s | Next: %s\n", p.Name, p.Next.Name)
		} else {
			fmt.Printf("Preset: %s\n", p.Name)
		}
	}

	return nil
}

// Regex for validating preset names
func IsPreset(s string) bool {
	regexPreset := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	return regexPreset.MatchString(s)
}

// ParseWaveformType parses a string into a WaveformType
// Returns an error if the string does not match any known waveform type
func ParseWaveformType(s string) (audio.WaveformType, error) {
	switch s {
	case keywordSine:
		return audio.WaveformSine, nil
	case keywordSquare:
		return audio.WaveformSquare, nil
	case keywordTriangle:
		return audio.WaveformTriangle, nil
	case keywordSawtooth:
		return audio.WaveformSawtooth, nil
	default:
		return -1, fmt.Errorf("invalid waveform type: %s", s)
	}
}
