package types

const (
	MaxPresets     = 32        // Maximum number of presets
	BuiltinSilence = "silence" // Represents silence built-in preset
)

// Preset represents a named preset
type Preset struct {
	Name  string                  // Name of preset
	Voice [NumberOfChannels]Voice // Voice-set for it
}
