/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package hub

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetCacheDir(ts *testing.T) {
	tmp := ts.TempDir()
	origHome := os.Getenv("HOME")
	origLocal := os.Getenv("LOCALAPPDATA")
	origXDG := os.Getenv("XDG_CACHE_HOME")

	ts.Cleanup(func() {
		os.Setenv("HOME", origHome)
		os.Setenv("LOCALAPPDATA", origLocal)
		os.Setenv("XDG_CACHE_HOME", origXDG)
	})

	tests := []struct {
		name         string
		goos         string
		setup        func()
		wantContains string
	}{
		{
			name:         "Darwin uses Library/Caches/org.synapseq",
			goos:         "darwin",
			setup:        func() { os.Setenv("HOME", tmp) },
			wantContains: "Library/Caches/org.synapseq",
		},
		{
			name: "Windows uses LOCALAPPDATA",
			goos: "windows",
			setup: func() {
				os.Setenv("LOCALAPPDATA", filepath.Join(tmp, "AppData", "Local"))
			},
			wantContains: "SynapSeq/Cache",
		},
		{
			name:         "Linux with XDG_CACHE_HOME",
			goos:         "linux",
			setup:        func() { os.Setenv("XDG_CACHE_HOME", filepath.Join(tmp, "xdg")) },
			wantContains: "xdg/synapseq",
		},
		{
			name: "Linux fallback to HOME/.cache/synapseq",
			goos: "linux",
			setup: func() {
				os.Unsetenv("XDG_CACHE_HOME")
				os.Setenv("HOME", tmp)
			},
			wantContains: ".cache/synapseq",
		},
	}

	for _, tt := range tests {
		ts.Run(tt.name, func(t *testing.T) {
			tt.setup()

			gotGOOS := runtime.GOOS
			switch tt.goos {
			case "darwin":
				if gotGOOS != "darwin" {
					t.Skip("skipping: test for darwin")
				}
			case "windows":
				if gotGOOS != "windows" {
					t.Skip("skipping: test for windows")
				}
			default:
				if gotGOOS != "linux" && gotGOOS != "freebsd" {
					t.Skip("skipping: test for linux/freebsd")
				}
			}

			dir, err := GetCacheDir()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if _, err := os.Stat(dir); err != nil {
				t.Errorf("expected directory to exist, but got error: %v", err)
			}

			if !filepath.IsAbs(dir) {
				t.Errorf("expected absolute path, got %s", dir)
			}

			if !strings.Contains(filepath.ToSlash(dir), tt.wantContains) {
				t.Errorf("expected path to contain %s, got %s", tt.wantContains, dir)
			}
		})
	}
}
