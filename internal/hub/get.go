/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package hub

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	t "github.com/ruanklein/synapseq/v3/internal/types"
)

// HubGet retrieves a sequence by its ID from the Hub
func HubGet(sequenceID string) (*t.HubEntry, error) {
	manifest, err := GetManifest()
	if err != nil {
		return nil, err
	}

	var entry *t.HubEntry
	for _, e := range manifest.Entries {
		if e.ID == sequenceID {
			entry = &e
			break
		}
	}

	return entry, nil
}

// HubDownload downloads a sequence and its dependencies from the Hub
func HubDownload(entry *t.HubEntry) (string, error) {
	if entry == nil {
		return "", fmt.Errorf("hub entry is nil")
	}

	cache, err := GetCacheDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(cache, strings.TrimSuffix(entry.Path, ".spsq"))
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}

	sequencePath := filepath.Join(path, entry.Name+".spsq")
	if _, err := os.Stat(sequencePath); err == nil {
		return sequencePath, nil
	}

	for _, dep := range entry.Dependencies {
		var depPath string
		if dep.Type == t.HubDependencyTypeBackground {
			depPath = filepath.Join(path, dep.Name+".wav")
		} else {
			depPath = filepath.Join(path, dep.Name+".spsq")
		}

		resp, err := http.Get(dep.DownloadUrl)
		if err != nil {
			return "", fmt.Errorf("error downloading dependency %s: %v", dep.Name, err)
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error reading dependency %s: %v", dep.Name, err)
		}

		if err = os.WriteFile(depPath, data, 0644); err != nil {
			return "", fmt.Errorf("error saving dependency %s: %v", dep.Name, err)
		}
	}

	resp, err := http.Get(entry.DownloadUrl)
	if err != nil {
		return "", fmt.Errorf("error downloading sequence %s: %v", entry.Name, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading sequence %s: %v", entry.Name, err)
	}

	if err = os.WriteFile(sequencePath, data, 0644); err != nil {
		return "", fmt.Errorf("error saving sequence %s: %v", entry.Name, err)
	}

	return sequencePath, nil
}
