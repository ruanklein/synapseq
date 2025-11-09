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
	"strings"

	t "github.com/ruanklein/synapseq/v3/internal/types"
)

func HubUpdate() error {
	cache, err := GetCacheDir()
	if err != nil {
		return err
	}

	resp, err := http.Get(t.HubManifestURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return fmt.Errorf("invalid content-type for manifest file: %s", contentType)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	manifestPath := cache + "/manifest.json"
	if err = os.WriteFile(manifestPath, data, 0644); err != nil {
		return err
	}

	return nil
}
