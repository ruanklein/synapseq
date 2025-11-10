/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package hub

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ruanklein/synapseq/v3/internal/info"
	t "github.com/ruanklein/synapseq/v3/internal/types"
)

// GetManifest retrieves and parses the Hub manifest file from the cache
func GetManifest() (*t.HubManifest, error) {
	cache, err := GetCacheDir()
	if err != nil {
		return nil, err
	}

	manifestPath := cache + "/manifest.json"
	manifest, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	var hubManifest *t.HubManifest
	if err := json.Unmarshal(manifest, &hubManifest); err != nil {
		return nil, err
	}

	if hubManifest.Version != info.HUB_VERSION {
		return nil, fmt.Errorf(
			"hub manifest version mismatch: expected %s, got %s\n"+
				"please update SynapSeq to the latest version to ensure compatibility",
			info.HUB_VERSION, hubManifest.Version,
		)
	}

	return hubManifest, nil
}
