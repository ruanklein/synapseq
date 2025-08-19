package sequence

import (
	"fmt"
)

const (
	KeywordComment    = "#"          // Represents a comment
	KeywordOption     = "@"          // Represents an option
	KeywordTone       = "tone"       // Represents a tone
	KeywordBinaural   = "binaural"   // Represents a binaural tone
	KeywordMonaural   = "monaural"   // Represents a monaural tone
	KeywordIsochronic = "isochronic" // Represents an isochronic tone
	KeywordAmplitude  = "amplitude"  // Represents an amplitude
	KeywordNoise      = "noise"      // Represents a noise
	KeywordWhite      = "white"      // Represents a white noise
	KeywordPink       = "pink"       // Represents a pink noise
	KeywordBrown      = "brown"      // Represents a brown noise
	KeywordSpin       = "spin"       // Represents a spin
	KeywordWidth      = "width"      // Represents a width
	KeywordRate       = "rate"       // Represents a rate
	KeywordIntensity  = "intensity"  // Represents an intensity

	// Regex for validating preset names
	// RegexPreset = `^[a-zA-Z][a-zA-Z0-9_-]*$`
)

// LoadSequence loads a sequence from a file
// and parses its contents into a preset
// and create channels as needed
func LoadSequence(fileName string) error {
	file, err := LoadFile(fileName)
	if err != nil {
		return fmt.Errorf("error loading sequence file: %v", err)
	}

	for file.NextLine() {
		// Debugging output
		fmt.Printf("Processing line %d: %s\n", file.CurrentLineNumber, file.CurrentLine)
	}

	return nil
}
