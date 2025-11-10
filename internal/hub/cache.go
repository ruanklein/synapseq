/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package hub

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetCacheDir returns the path to the cache directory for storing Hub data
func GetCacheDir() (string, error) {
	var base string

	switch runtime.GOOS {
	case "darwin":
		base = filepath.Join(os.Getenv("HOME"), "Library", "Caches", "org.synapseq")
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return "", os.ErrNotExist
		}
		base = filepath.Join(localAppData, "SynapSeq", "Cache")
	default: // Linux, BSD, etc.
		if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
			base = filepath.Join(xdg, "synapseq")
		} else {
			base = filepath.Join(os.Getenv("HOME"), ".cache", "synapseq")
		}
	}

	if err := os.MkdirAll(base, 0755); err != nil {
		return "", err
	}

	return base, nil
}

// HubClean removes the cache directory and all its contents
func HubClean() error {
	cache, err := GetCacheDir()
	if err != nil {
		return err
	}

	if err = os.RemoveAll(cache); err != nil {
		return err
	}

	return nil
}
