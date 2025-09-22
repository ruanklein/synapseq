/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package sequence

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(ts *testing.T, name string, lines []string) string {
	ts.Helper()
	dir := ts.TempDir()
	p := filepath.Join(dir, name)
	f, err := os.Create(p)
	if err != nil {
		ts.Fatalf("create temp file: %v", err)
	}
	for _, ln := range lines {
		if _, err := f.WriteString(ln + "\n"); err != nil {
			ts.Fatalf("write temp file: %v", err)
		}
	}
	if err := f.Close(); err != nil {
		ts.Fatalf("close temp file: %v", err)
	}
	return p
}

func TestLoadFile_FromFilePath(ts *testing.T) {
	want := []string{"line-1", "line-2", "line-3"}
	path := writeTempFile(ts, "seq.spsq", want)

	sf, err := LoadFile(path)
	if err != nil {
		ts.Fatalf("LoadFile(%s) error: %v", path, err)
	}
	defer os.Remove(path)

	var got []string
	for sf.NextLine() {
		got = append(got, sf.CurrentLine)
	}
	if len(got) != len(want) {
		ts.Fatalf("expected %d lines, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			ts.Errorf("line %d: expected %q, got %q", i+1, want[i], got[i])
		}
	}
	if sf.CurrentLineNumber != len(want) {
		ts.Errorf("expected CurrentLineNumber=%d, got %d", len(want), sf.CurrentLineNumber)
	}

	// Must close underlying file descriptor
	if sf.file == nil {
		ts.Fatalf("expected non-nil file handle")
	}
	sf.Close()
	if _, err := sf.file.Stat(); err == nil {
		ts.Errorf("expected error when stat on closed file, got nil")
	}

	// Idempotency
	sf.Close()
}

func TestLoadFile_NotFound(ts *testing.T) {
	if _, err := LoadFile(filepath.Join(ts.TempDir(), "missing.spsq")); err == nil {
		ts.Errorf("expected error for missing file, got nil")
	}
}

func TestNextLine_NoScanner(ts *testing.T) {
	sf := &SequenceFile{}
	if sf.NextLine() {
		ts.Errorf("expected NextLine()=false with nil scanner")
	}
}
