//go:build !nohub

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package hub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/ruanklein/synapseq/v3/internal/info"
	t "github.com/ruanklein/synapseq/v3/internal/types"
)

// TrackDownload sends an anonymous download event to the SynapSeq Hub analytics endpoint.
// It only sends technical metadata, no personal or identifying information.
func TrackDownload(sequenceID string, action t.HubActionTracking) error {
	if sequenceID == "" {
		return fmt.Errorf("empty sequence ID")
	}

	payload := map[string]string{
		"id": sequenceID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", t.HubTrackEndpoint, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-SYNAPSEQ-SOURCE", "CLI")
	req.Header.Set("X-SYNAPSEQ-VERSION", info.VERSION)
	req.Header.Set("X-SYNAPSEQ-PLATFORM", runtime.GOOS)
	req.Header.Set("X-SYNAPSEQ-ARCH", runtime.GOARCH)
	req.Header.Set("X-SYNAPSEQ-ACTION", strings.ToUpper(action.String()))

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// Fail silently. Tracking must never break CLI functionality
		return nil
	}
	defer resp.Body.Close()

	return nil
}
