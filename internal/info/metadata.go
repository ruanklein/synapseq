/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package info

import (
	"encoding/base64"
	"runtime"
	"time"

	"github.com/google/uuid"
)

// Metadata holds the embedded metadata information
type Metadata struct {
	// ID is the unique identifier for the metadata
	id string
	// Generated is the timestamp of when the file was generated
	generated string
	// Version of the application that generated the file
	version string
	// Platform is the target platform for the generated file
	platform string
	// Content is the actual embedded content (e.g., sequence data)
	content string
}

// NewMetadata creates a new Metadata instance with current information
func NewMetadata(content []byte) (*Metadata, error) {
	return &Metadata{
		id:        uuid.New().String(),
		generated: time.Now().UTC().Format(time.RFC3339),
		version:   VERSION,
		platform:  runtime.GOOS + "/" + runtime.GOARCH,
		content:   base64.StdEncoding.EncodeToString(content),
	}, nil
}

// ID returns the unique identifier
func (m *Metadata) ID() string {
	return m.id
}

// Generated returns the generation timestamp
func (m *Metadata) Generated() string {
	return m.generated
}

// Version returns the application version
func (m *Metadata) Version() string {
	return m.version
}

// Platform returns the target platform
func (m *Metadata) Platform() string {
	return m.platform
}

// Content returns the embedded content
func (m *Metadata) Content() string {
	return m.content
}
