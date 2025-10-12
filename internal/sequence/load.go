/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package sequence

import t "github.com/ruanklein/synapseq/internal/types"

// LoadResult holds the result of loading a sequence
type LoadResult struct {
	Periods  []t.Period
	Options  *t.Option
	Comments []string
}
