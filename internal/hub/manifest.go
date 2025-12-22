//go:build !nohub

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package hub

import (
	"encoding/json"
	"os"

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

	return hubManifest, nil
}

// ManifestExists checks if the Hub manifest file exists in the cache
func ManifestExists() bool {
	cache, err := GetCacheDir()
	if err != nil {
		return false
	}

	manifestPath := cache + "/manifest.json"
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return false
	}

	return true
}
