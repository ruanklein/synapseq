/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package sequence

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

func makeBigContent(lineLen, minBytes int) []string {
	line := strings.Repeat("X", lineLen)
	var (
		lines []string
		size  int
	)
	for size <= minBytes {
		lines = append(lines, line)
		size += lineLen + 1
	}
	return lines
}

func countConsumedBytes(readLines []string, fullLineLen int) (bytes, fullLines int, hadPartial bool) {
	for i, ln := range readLines {
		if len(ln) == fullLineLen {
			bytes += fullLineLen + 1
			fullLines++
		} else {
			if i != len(readLines)-1 {
				bytes += len(ln) + 1
			} else {
				bytes += len(ln)
				hadPartial = true
			}
		}
	}
	return
}

func TestLoadFile_Truncate_FromFilePath(ts *testing.T) {
	const (
		maxSize  = 32 * 1024
		lineLen  = 123 // 123+1=124, 32768%124 ≠ 0 --> partial line at the end
		minBytes = maxSize + 8192
	)
	lines := makeBigContent(lineLen, minBytes)
	path := writeTempFile(ts, "big-local.spsq", lines)

	sf, err := LoadFile(path)
	if err != nil {
		ts.Fatalf("LoadFile(%s) error: %v", path, err)
	}
	defer sf.Close()

	var got []string
	for sf.NextLine() {
		got = append(got, sf.CurrentLine)
	}
	if len(got) == len(lines) {
		ts.Fatalf("was expecting truncation (>32KB), but all lines were read")
	}

	consumed, _, _ := countConsumedBytes(got, lineLen)
	if consumed != maxSize {
		ts.Fatalf("unexpected consumed bytes (local file): got=%d want=%d", consumed, maxSize)
	}
}

func TestLoadFile_Truncate_FromStdin(ts *testing.T) {
	const (
		maxSize  = 32 * 1024
		lineLen  = 123
		minBytes = maxSize + 4096
	)
	lines := makeBigContent(lineLen, minBytes)
	path := writeTempFile(ts, "big-stdin.spsq", lines)

	orig := os.Stdin
	f, err := os.Open(path)
	if err != nil {
		ts.Fatalf("open temp file: %v", err)
	}
	defer func() {
		_ = f.Close()
		os.Stdin = orig
	}()
	os.Stdin = f

	sf, err := LoadFile("-")
	if err != nil {
		ts.Fatalf("LoadFile(-) error: %v", err)
	}
	defer sf.Close()

	var got []string
	for sf.NextLine() {
		got = append(got, sf.CurrentLine)
	}
	if len(got) == len(lines) {
		ts.Fatalf("was expecting truncation (>32KB via stdin), but all lines were read")
	}
	consumed, _, _ := countConsumedBytes(got, lineLen)
	if consumed != maxSize {
		ts.Fatalf("unexpected consumed bytes (stdin): got=%d want=%d", consumed, maxSize)
	}
}

func TestLoadFile_Truncate_FromHTTP_LocalServer(ts *testing.T) {
	const (
		maxSize  = 32 * 1024
		lineLen  = 123
		minBytes = maxSize + 2048
	)
	lines := makeBigContent(lineLen, minBytes)
	body := strings.Join(lines, "\n") + "\n"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	sf, err := LoadFile(srv.URL + "/big.spsq")
	if err != nil {
		ts.Fatalf("LoadFile(HTTP) error: %v", err)
	}
	defer sf.Close()

	var got []string
	for sf.NextLine() {
		got = append(got, sf.CurrentLine)
	}
	if len(got) == len(lines) {
		ts.Fatalf("was expecting truncation (>32KB via HTTP), but all lines were read")
	}
	consumed, _, _ := countConsumedBytes(got, lineLen)
	if consumed != maxSize {
		ts.Fatalf("unexpected consumed bytes (HTTP): got=%d want=%d", consumed, maxSize)
	}
}

func TestLoadFile_HTTP_InvalidContentType(ts *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	if _, err := LoadFile(srv.URL + "/data.spsq"); err == nil {
		ts.Fatalf("expected error for invalid content-type, got nil")
	} else if !strings.Contains(err.Error(), "invalid content-type") {
		ts.Fatalf("unexpected error message: %v", err)
	}
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

func TestLoadFile_FromStdin(ts *testing.T) {
	want := []string{"a", "b"}
	path := writeTempFile(ts, "stdin.spsq", want)

	orig := os.Stdin
	f, err := os.Open(path)
	if err != nil {
		ts.Fatalf("open temp file: %v", err)
	}
	defer func() {
		f.Close()
		os.Stdin = orig
	}()

	os.Stdin = f
	sf, err := LoadFile("-")
	if err != nil {
		ts.Fatalf("LoadFile(-) error: %v", err)
	}
	if sf.file != nil {
		ts.Errorf("for '-', expected sf.file=nil, got non-nil")
	}

	var got []string
	for sf.NextLine() {
		got = append(got, sf.CurrentLine)
	}
	if len(got) != len(want) {
		ts.Fatalf("expected %d lines from stdin, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			ts.Errorf("stdin line %d: expected %q, got %q", i+1, want[i], got[i])
		}
	}

	// No-op close when file is nil
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
