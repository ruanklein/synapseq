/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package sequence

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// maxFileSize defines the maximum file size (32KB)
const maxFileSize = 32 * 1024

// SequenceFile represents a sequence file
type SequenceFile struct {
	CurrentLine       string         // Current line in the file
	CurrentLineNumber int            // Current line number
	scanner           *bufio.Scanner // Scanner for reading the file
	file              *os.File       // File handle
}

// loadRemoteFile loads a remote file from a given URL (max 32KB and text/plain)
func loadRemoteFile(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching remote file: %v", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/plain") {
		return nil, fmt.Errorf("invalid content-type: %s", contentType)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxFileSize))
	if err != nil {
		return nil, fmt.Errorf("error reading remote file: %v", err)
	}

	return strings.NewReader(string(data)), nil
}

// LoadFile loads a sequence file
func LoadFile(fileName string) (*SequenceFile, error) {
	if fileName == "-" {
		reader := io.LimitReader(os.Stdin, maxFileSize)
		return &SequenceFile{
			scanner: bufio.NewScanner(reader),
			file:    nil,
		}, nil
	}

	if strings.HasPrefix(fileName, "http://") || strings.HasPrefix(fileName, "https://") {
		reader, err := loadRemoteFile(fileName)
		if err != nil {
			return nil, err
		}
		return &SequenceFile{
			scanner: bufio.NewScanner(reader),
			file:    nil,
		}, nil
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	reader := io.LimitReader(file, maxFileSize)

	return &SequenceFile{
		scanner: bufio.NewScanner(reader),
		file:    file,
	}, nil
}

// NextLine advances to the next line in the sequence file
func (sf *SequenceFile) NextLine() bool {
	if sf.scanner == nil {
		return false
	}

	if sf.scanner.Scan() {
		sf.CurrentLine = sf.scanner.Text()
		sf.CurrentLineNumber++
		return true
	}
	return false
}

// Close closes the sequence file
func (sf *SequenceFile) Close() {
	if sf.file != nil {
		sf.file.Close()
	}
}
