/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package sequence

import (
	"encoding/json"
	"fmt"
	"os"

	t "github.com/ruanklein/synapseq/internal/types"
)

func LoadJSONSequence(filename string) (*LoadResult, error) {
	var result LoadResult

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	var jsonInput t.SynapSeqInput
	if err := json.Unmarshal(data, &jsonInput); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	if len(jsonInput.Sequence) == 0 {
		return nil, fmt.Errorf("no sequence data found in JSON")
	}

	gainLevel := t.GainLevelMedium
	if jsonInput.Options.GainLevel != "" {
		switch jsonInput.Options.GainLevel {
		case t.KeywordOptionGainLevelVeryLow:
			gainLevel = t.GainLevelVeryLow
		case t.KeywordOptionGainLevelLow:
			gainLevel = t.GainLevelLow
		case t.KeywordOptionGainLevelMedium:
			gainLevel = t.GainLevelMedium
		case t.KeywordOptionGainLevelHigh:
			gainLevel = t.GainLevelHigh
		case t.KeywordOptionGainLevelVeryHigh:
			gainLevel = t.GainLevelVeryHigh
		default:
			return nil, fmt.Errorf("invalid gain level: %s", jsonInput.Options.GainLevel)
		}
	}

	// Initialize audio options
	options := &t.Option{
		SampleRate:     jsonInput.Options.Samplerate,
		Volume:         jsonInput.Options.Volume,
		BackgroundPath: jsonInput.Options.BackgroundPath,
		PresetPath:     jsonInput.Options.PresetPath,
		GainLevel:      gainLevel,
	}

	if err := options.Validate(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	var periods []t.Period
	for idx, seq := range jsonInput.Sequence {
		if idx == 0 && seq.Time != 0 {
			return nil, fmt.Errorf("first timeline must start at 0ms (00:00:00)")
		}
		if idx >= 1 && seq.Time <= jsonInput.Sequence[idx-1].Time {
			return nil, fmt.Errorf("timeline %d time must be greater than previous timeline time", idx+1)
		}

	}

	return &result, nil
}
