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
	"os"
)

// SequenceFile represents a sequence file
type SequenceFile struct {
	CurrentLine       string         // Current line in the file
	CurrentLineNumber int            // Current line number
	scanner           *bufio.Scanner // Scanner for reading the file
	file              *os.File       // File handle
}

// LoadFile loads a sequence file
func LoadFile(fileName string) (*SequenceFile, error) {
	var file *os.File

	if fileName == "-" {
		return &SequenceFile{
			scanner: bufio.NewScanner(os.Stdin),
			file:    file,
		}, nil
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &SequenceFile{
		scanner: bufio.NewScanner(file),
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
