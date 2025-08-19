package sequence

import (
	"bufio"
	"fmt"
	"os"
)

// SequenceFile represents a sequence file
type SequenceFile struct {
	CurrentLine       string
	CurrentLineNumber int
	scanner           *bufio.Scanner
	file              *os.File
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

	var err error
	file, err = os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", fileName, err)
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
