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
	"path/filepath"
	"testing"

	"github.com/ruanklein/synapseq/v3/internal/info"
	t "github.com/ruanklein/synapseq/v3/internal/types"
)

func TestGetManifest(ts *testing.T) {
	tmp := ts.TempDir()
	os.Setenv("HOME", tmp)

	cacheDir := filepath.Join(tmp, "Library", "Caches", "org.synapseq")
	os.MkdirAll(cacheDir, 0755)

	valid := t.HubManifest{
		Version:     info.HUB_VERSION,
		LastUpdated: "2025-11-09T00:00:00Z",
		Entries: []t.HubEntry{
			{Name: "test", Author: "ruan"},
		},
	}
	data, _ := json.Marshal(valid)
	os.WriteFile(filepath.Join(cacheDir, "manifest.json"), data, 0644)

	got, err := GetManifest()
	if err != nil {
		ts.Fatalf("unexpected error: %v", err)
	}
	if got.Version != info.HUB_VERSION {
		ts.Errorf("expected version %s, got %s", info.HUB_VERSION, got.Version)
	}

	invalid := valid
	invalid.Version = "999.0.0"
	data, _ = json.Marshal(invalid)
	os.WriteFile(filepath.Join(cacheDir, "manifest.json"), data, 0644)

	if _, err := GetManifest(); err == nil {
		ts.Error("expected version mismatch error, got nil")
	}
}
