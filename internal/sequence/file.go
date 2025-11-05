/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package sequence

import (
	"bufio"

	s "github.com/ruanklein/synapseq/v3/internal/shared"
	t "github.com/ruanklein/synapseq/v3/internal/types"
)

// SequenceFile represents a sequence file
type SequenceFile struct {
	CurrentLine       string         // Current line in the file
	CurrentLineNumber int            // Current line number
	scanner           *bufio.Scanner // Scanner for reading the file
}

// LoadFile loads a sequence file
func LoadFile(fileName string) (*SequenceFile, error) {
	file, err := s.GetFile(fileName, t.FormatText)
	if err != nil {
		return nil, err
	}

	return &SequenceFile{
		scanner: bufio.NewScanner(file),
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
