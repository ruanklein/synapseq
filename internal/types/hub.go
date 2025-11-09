/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package types

const (
	// Hub Base URL for the Hub repository
	HubBaseURL = "https://ruanklein.github.io/synapseq-hub"
	// HubManifestURL is the URL to fetch the Hub manifest
	HubManifestURL = "https://ruanklein.github.io/synapseq-hub/manifest.json"
)

// HubDependencyType represents the type of a Hub dependency
type HubDependencyType string

const (
	// HubDependencyTypePresetList represents a preset list dependency
	HubDependencyTypePresetList HubDependencyType = "presetlist"
	// HubDependencyTypeBackground represents a background dependency
	HubDependencyTypeBackground HubDependencyType = "background"
)

// String returns the string representation of the HubDependencyType
func (dt HubDependencyType) String() string {
	return string(dt)
}

// HubDependency represents a dependency for a Hub entry
type HubDependency struct {
	Type        HubDependencyType `json:"type"`
	Name        string            `json:"name"`
	DownloadUrl string            `json:"download_url"`
}

// HubEntry represents an entry in the Hub index
type HubEntry struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Author       string          `json:"author"`
	Category     string          `json:"category"`
	Path         string          `json:"path"`
	DownloadUrl  string          `json:"download_url"`
	UpdatedAt    string          `json:"updated_at"`
	Dependencies []HubDependency `json:"dependencies,omitempty"`
}

// HubManifest represents the manifest of available Hub entries
type HubManifest struct {
	Version     string
	LastUpdated string
	Entries     []HubEntry
}
