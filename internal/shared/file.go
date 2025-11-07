/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package shared

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"

	t "github.com/ruanklein/synapseq/v3/internal/types"
)

// readFile reads a file from the given reader up to maxSize bytes
func readFile(r io.Reader, maxSize int64) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxSize))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// getRemoteFile fetches a remote file and validates its content type and size
func getRemoteFile(url string, maxSize int64, typ t.FileFormat) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching remote file: %v", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	switch typ {
	case t.FormatText:
		if !strings.HasPrefix(contentType, "text/plain") {
			return nil, fmt.Errorf("invalid content-type for text file: %s", contentType)
		}
	case t.FormatJSON:
		if !slices.Contains([]string{"application/json", "application/x-json", "text/json"}, contentType) {
			return nil, fmt.Errorf("invalid content-type for json file: %s", contentType)
		}
	case t.FormatXML:
		if !slices.Contains([]string{"application/xml", "text/xml"}, contentType) {
			return nil, fmt.Errorf("invalid content-type for xml file: %s", contentType)
		}
	case t.FormatYAML:
		if !slices.Contains([]string{"application/x-yaml", "application/yaml", "text/yaml", "text/x-yaml"}, contentType) {
			return nil, fmt.Errorf("invalid content-type for yaml file: %s", contentType)
		}
	case t.FormatWAV:
		if !slices.Contains([]string{"audio/wav", "audio/x-wav", "audio/wave", "audio/vnd.wave"}, contentType) {
			return nil, fmt.Errorf("invalid content-type for wav file: %s", contentType)
		}
	}

	data, err := readFile(resp.Body, maxSize)
	if err != nil {
		return nil, fmt.Errorf("error reading remote file: %v", err)
	}

	return data, nil
}

// IsRemoteFile checks if the given file path is a remote URL
func IsRemoteFile(filePath string) bool {
	return strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://")
}

// GetFile retrieves a file from a local path or URL based on the specified type
func GetFile(filePath string, typ t.FileFormat) ([]byte, error) {
	maxSize := int64(0)
	switch typ {
	case t.FormatText:
		maxSize = t.MaxTextFileSize
	case t.FormatJSON, t.FormatXML, t.FormatYAML:
		maxSize = t.MaxStructuredFileSize
	case t.FormatWAV:
		maxSize = t.MaxBackgroundFileSize
	}

	if maxSize == 0 {
		return nil, fmt.Errorf("unsupported file type: %s", typ.String())
	}

	switch {
	case filePath == "-":
		data, err := readFile(os.Stdin, maxSize)
		if err != nil {
			return nil, fmt.Errorf("error reading from stdin: %v", err)
		}
		return data, nil

	case IsRemoteFile(filePath):
		return getRemoteFile(filePath, maxSize, typ)

	default:
		f, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %v", err)
		}
		defer f.Close()

		data, err := readFile(f, maxSize)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}
		return data, nil
	}
}
