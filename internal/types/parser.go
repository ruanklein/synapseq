package types

const (
	KeywordComment                 = "#"          // Represents a comment
	KeywordOption                  = "@"          // Represents an option
	KeywordOptionSampleRate        = "samplerate" // Represents a sample rate option
	KeywordOptionVolume            = "volume"     // Represents a volume option
	KeywordOptionBackground        = "background" // Represents a background option
	KeywordOptionGainLevel         = "gainlevel"  // Represents a gain level option
	KeywordOptionGainLevelVeryLow  = "verylow"    // Represents a very low gain level option
	KeywordOptionGainLevelLow      = "low"        // Represents a low gain level option
	KeywordOptionGainLevelMedium   = "medium"     // Represents a medium gain level option
	KeywordOptionGainLevelHigh     = "high"       // Represents a high gain level option
	KeywordOptionGainLevelVeryHigh = "veryhigh"   // Represents a very high gain level option
	KeywordWaveform                = "waveform"   // Represents a waveform
	KeywordSine                    = "sine"       // Represents a sine wave
	KeywordSquare                  = "square"     // Represents a square wave
	KeywordTriangle                = "triangle"   // Represents a triangle wave
	KeywordSawtooth                = "sawtooth"   // Represents a sawtooth wave
	KeywordTone                    = "tone"       // Represents a tone
	KeywordBinaural                = "binaural"   // Represents a binaural tone
	KeywordMonaural                = "monaural"   // Represents a monaural tone
	KeywordIsochronic              = "isochronic" // Represents an isochronic tone
	KeywordAmplitude               = "amplitude"  // Represents an amplitude
	KeywordNoise                   = "noise"      // Represents a noise
	KeywordWhite                   = "white"      // Represents a white noise
	KeywordPink                    = "pink"       // Represents a pink noise
	KeywordBrown                   = "brown"      // Represents a brown noise
	KeywordSpin                    = "spin"       // Represents a spin
	KeywordWidth                   = "width"      // Represents a width
	KeywordRate                    = "rate"       // Represents a rate
	KeywordEffect                  = "effect"     // Represents an effect
	KeywordBackground              = "background" // Represents a background
	KeywordPulse                   = "pulse"      // Represents a pulse
	KeywordIntensity               = "intensity"  // Represents an intensity
)

// Parser defines the interface for parsing different line types
type Parser interface {
	// ParseComment parses a comment line
	ParseComment() (string, error)
	// ParseOption parses an option line
	ParseOption(*Option) error
	// ParsePreset parses a preset line
	ParsePreset() (*Preset, error)
	// ParseVoice parses a voice line
	ParseVoice() (*Voice, error)
}
