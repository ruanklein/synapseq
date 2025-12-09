//go:build !wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package external

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type FFprobe struct{ baseUtility }

// NewFFprobe creates a new FFprobe instance with given ffprobe path
func NewFFprobe(ffprobePath string) (*FFprobe, error) {
	if ffprobePath == "" {
		ffprobePath = "ffprobe"
	}

	util, err := newUtility(ffprobePath)
	if err != nil {
		return nil, err
	}

	return &FFprobe{baseUtility: *util}, nil
}

// NewFFprobeUnsafe creates an FFprobe instance without validating the path.
// Useful for documentation examples and testing environments.
func NewFFprobeUnsafe(path string) *FFprobe {
	if path == "" {
		path = "ffprobe"
	}
	return &FFprobe{baseUtility: baseUtility{path: path}}
}

// extractMetadata extracts metadata tags from inputFile for mp3/ogg/opus files.
func (fp *FFprobe) extractMetadata(inputFile string) (map[string]string, error) {
	if inputFile == "" {
		return nil, fmt.Errorf("input file cannot be empty")
	}

	cmd := fp.Command(
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		inputFile,
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffprobe failed: %v: %s", err, stderr.String())
	}

	var probe struct {
		Format struct {
			Tags map[string]string `json:"tags"`
		} `json:"format"`
		Streams []struct {
			Tags map[string]string `json:"tags"`
		} `json:"streams"`
	}

	if err := json.Unmarshal(out.Bytes(), &probe); err != nil {
		return nil, fmt.Errorf("failed to parse ffprobe output: %v", err)
	}

	result := make(map[string]string)

	if probe.Format.Tags != nil {
		for k, v := range probe.Format.Tags {
			result[k] = v
		}
	}

	for _, s := range probe.Streams {
		if s.Tags == nil {
			continue
		}
		for k, v := range s.Tags {
			result[k] = v
		}
	}

	if b64, ok := result["synapseq_content"]; ok && b64 != "" {
		decoded, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			// Try URL-safe base64 as fallback (some producers use it)
			decoded, err = base64.URLEncoding.DecodeString(b64)
			if err != nil {
				return nil, fmt.Errorf("failed to base64 decode synapseq_content: %v", err)
			}
		}
		result["synapseq_content"] = string(decoded)
	}

	return result, nil
}

// ExtractTextSequence reconstructs the original text sequence metadata
// from MP3/OGG/OPUS files using ffprobe.
func (fp *FFprobe) ExtractTextSequence(inputFile string) (string, error) {
	if inputFile == "" {
		return "", fmt.Errorf("input file cannot be empty")
	}

	meta, err := fp.extractMetadata(inputFile)
	if err != nil {
		return "", err
	}

	id := meta["synapseq_id"]
	gen := meta["synapseq_generated"]
	ver := meta["synapseq_version"]
	plat := meta["synapseq_platform"]
	content := meta["synapseq_content"]

	if id == "" || gen == "" || ver == "" || plat == "" || content == "" {
		return "", fmt.Errorf("missing required synapseq_* metadata fields in file")
	}

	var out strings.Builder

	out.WriteString("# ================================================\n")
	out.WriteString("#  This sequence was exported from SynapSeq\n")
	out.WriteString(fmt.Sprintf("#  ID       : %s\n", id))
	out.WriteString(fmt.Sprintf("#  Date     : %s\n", gen))
	out.WriteString(fmt.Sprintf("#  Version  : %s\n", ver))
	out.WriteString(fmt.Sprintf("#  Platform : %s\n", plat))
	out.WriteString("# ================================================\n\n\n")

	out.WriteString(content)

	return out.String(), nil
}

// SaveExtractedTextSequence extracts and saves the text sequence
// from inputFile into outputFile.
func (fp *FFprobe) SaveExtractedTextSequence(inputFile, outputFile string) error {
	content, err := fp.ExtractTextSequence(inputFile)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write extracted text sequence to file: %v", err)
	}

	return nil
}
