package sequence

import (
	"fmt"
	"strings"
)

const (
	keywordComment    = "#"          // Represents a comment
	keywordOption     = "@"          // Represents an option
	keywordWaveform   = "waveform"   // Represents a waveform
	keywordSine       = "sine"       // Represents a sine wave
	keywordSquare     = "square"     // Represents a square wave
	keywordTriangle   = "triangle"   // Represents a triangle wave
	keywordSawtooth   = "sawtooth"   // Represents a sawtooth wave
	keywordTone       = "tone"       // Represents a tone
	keywordBinaural   = "binaural"   // Represents a binaural tone
	keywordMonaural   = "monaural"   // Represents a monaural tone
	keywordIsochronic = "isochronic" // Represents an isochronic tone
	keywordAmplitude  = "amplitude"  // Represents an amplitude
	keywordNoise      = "noise"      // Represents a noise
	keywordWhite      = "white"      // Represents a white noise
	keywordPink       = "pink"       // Represents a pink noise
	keywordBrown      = "brown"      // Represents a brown noise
	keywordSpin       = "spin"       // Represents a spin
	keywordWidth      = "width"      // Represents a width
	keywordRate       = "rate"       // Represents a rate
	keywordIntensity  = "intensity"  // Represents an intensity
)

// isDoubleComment checks if a string is a double comment
func isDoubleComment(s string) bool {
	return s == strings.Repeat(keywordComment, 2)
}

// LoadSequence loads a sequence from a file
func LoadSequence(fileName string) error {
	file, err := LoadFile(fileName)
	if err != nil {
		return fmt.Errorf("error loading sequence file: %v", err)
	}

	for file.NextLine() {
		fields := strings.Fields(file.CurrentLine)

		// Skip empty lines and comments
		if len(fields) == 0 || fields[0] == keywordComment {
			continue
		}

		// Printable comment
		if isDoubleComment(fields[0]) {
			fmt.Printf("> %s\n", strings.Join(fields[1:], " "))
			continue
		}

		// Preset definition
		if !strings.HasPrefix(file.CurrentLine, " ") && IsPreset(fields[0]) {
			presetName := strings.ToLower(fields[0])
			if presetName == BuiltinSilence {
				return fmt.Errorf("cannot load built-in preset '%s' at line %d", BuiltinSilence, file.CurrentLineNumber)
			}

			// Check for a comment after the preset name
			if len(fields) > 1 {
				if isDoubleComment(fields[1]) {
					fmt.Printf("> %s\n", strings.Join(fields[2:], " "))
				} else if fields[1] == keywordComment {
					// Skip
				} else {
					return fmt.Errorf("invalid preset definition at line %d: %s", file.CurrentLineNumber, file.CurrentLine)
				}
			}

			var preset Preset
			if len(PresetList) > 0 {
				for i := 0; i < len(PresetList); i++ {
					if PresetList[i].Name == presetName {
						return fmt.Errorf("preset '%s' already exists at line %d", presetName, file.CurrentLineNumber)
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
