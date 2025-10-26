/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package sequence

import (
	"fmt"

	t "github.com/ruanklein/synapseq/internal/types"
)

// ConvertToText converts a slice of Periods to a text-based sequence file.
func ConvertToText(sequence *t.Sequence) (string, error) {
	content := "# GENERATED FROM SYNAPSEQ STRUCTURED SEQUENCE FILE\n\n"
	for _, comments := range sequence.Comments {
		content += fmt.Sprintf("## %s\n", comments)
	}

	options := sequence.Options
	if options != nil {
		content += "\n# Options\n"
		content += fmt.Sprintf("%s%s %d", t.KeywordOption, t.KeywordOptionSampleRate, options.SampleRate)
		content += fmt.Sprintf("\n%s%s %d", t.KeywordOption, t.KeywordOptionVolume, options.Volume)

		if options.BackgroundPath != "" {
			content += fmt.Sprintf("\n%s%s %s", t.KeywordOption, t.KeywordOptionBackground, options.BackgroundPath)
			content += fmt.Sprintf("\n%s%s %s", t.KeywordOption, t.KeywordOptionGainLevel, options.GainLevel.String())
		}
		content += "\n"
	}

	content += "\n# Presets"

	var timeline []string
	for i, period := range sequence.Periods {
		presetID := fmt.Sprintf("tone-set-%03d", i+1)
		content += fmt.Sprintf("\n%s", presetID)

		for _, track := range period.TrackStart {
			if track.Type != t.TrackOff {
				content += fmt.Sprintf("\n  %s", track.String())
			}
		}

		timeline = append(timeline, fmt.Sprintf("%s %s %s", period.TimeString(), presetID, period.Transition.String()))
	}

	content += "\n\n# Timeline"

	for _, tline := range timeline {
		content += fmt.Sprintf("\n%s", tline)
	}

	return content, nil
}
