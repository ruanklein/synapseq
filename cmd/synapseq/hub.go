/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	synapseq "github.com/ruanklein/synapseq/v3/core"
	"github.com/ruanklein/synapseq/v3/internal/hub"
	t "github.com/ruanklein/synapseq/v3/internal/types"
)

// hubRunUpdate updates the local Hub manifest
func hubRunUpdate() error {
	if err := hub.HubUpdate(); err != nil {
		return fmt.Errorf("failed to update hub. Error\n  %v", err)
	}
	manifest, err := hub.GetManifest()
	if err != nil {
		return fmt.Errorf("failed to get hub manifest. Error\n  %v", err)
	}
	fmt.Printf("Fetched %d entries from the Hub. Last update: %s\n", len(manifest.Entries), manifest.LastUpdated)
	return nil
}

// hubRunClean cleans the local Hub cache
func hubRunClean() error {
	if err := hub.HubClean(); err != nil {
		return fmt.Errorf("failed to clean hub cache. Error\n  %v", err)
	}
	fmt.Println("Hub cache cleaned successfully.")
	return nil
}

// hubRunGet retrieves and processes a sequence from the Hub
func hubRunGet(sequenceId, outputFile string, quiet bool) error {
	entry, err := hub.HubGet(sequenceId)
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}
	if entry == nil {
		return fmt.Errorf("sequence not found in hub: %s", sequenceId)
	}

	inputFile, err := hub.HubDownload(entry)
	if err != nil {
		return fmt.Errorf("failed to download sequence from hub. Error\n  %v", err)
	}

	appCtx, err := synapseq.NewAppContext(inputFile, outputFile, "text")
	if err != nil {
		return fmt.Errorf("failed to create application context. Error\n  %v", err)
	}

	if !quiet && outputFile != "-" {
		appCtx = appCtx.WithVerbose(os.Stdout)
	}

	if err := appCtx.LoadSequence(); err != nil {
		return fmt.Errorf("failed to load sequence. Error\n  %v", err)
	}

	if outputFile == "-" {
		if err := appCtx.Stream(os.Stdout); err != nil {
			return fmt.Errorf("failed to stream sequence. Error\n  %v", err)
		}
		return nil
	}

	if !quiet {
		for _, c := range appCtx.Comments() {
			fmt.Printf("> %s\n", c)
		}
	}

	if err := appCtx.WAV(); err != nil {
		return fmt.Errorf("failed to generate WAV file. Error\n  %v", err)
	}
	return nil
}

// / hubRunList prints all available sequences from the Hub manifest in a tabular format
func hubRunList() error {
	manifest, err := hub.GetManifest()
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}

	fmt.Printf("SynapSeq Hub — %d available sequences  (Last updated: %s)\n\n",
		len(manifest.Entries), manifest.LastUpdated)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tAUTHOR\tCATEGORY\tUPDATED")

	for _, e := range manifest.Entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			e.ID,
			e.Author,
			e.Category,
			e.UpdatedAt[:10],
		)
	}

	w.Flush()
	return nil
}

// hubRunSearch searches for sequences in the Hub by keyword (case-insensitive)
func hubRunSearch(query string) error {
	if strings.TrimSpace(query) == "" {
		return fmt.Errorf("missing search term")
	}

	manifest, err := hub.GetManifest()
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}

	query = strings.ToLower(query)
	var results []t.HubEntry

	for _, e := range manifest.Entries {
		if strings.Contains(strings.ToLower(e.ID), query) ||
			strings.Contains(strings.ToLower(e.Name), query) ||
			strings.Contains(strings.ToLower(e.Author), query) ||
			strings.Contains(strings.ToLower(e.Category), query) {
			results = append(results, e)
		}
	}

	if len(results) == 0 {
		fmt.Printf("No matches found for %q\n", query)
		return nil
	}

	fmt.Printf("SynapSeq Hub - %d matching results for %q\n\n", len(results), query)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tAUTHOR\tCATEGORY\tUPDATED")

	for _, e := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			e.ID,
			e.Author,
			e.Category,
			e.UpdatedAt[:10],
		)
	}

	w.Flush()
	return nil
}

// hubRunDownload downloads a sequence and all its dependencies into a given folder
func hubRunDownload(sequenceID, targetDir string) error {
	if strings.TrimSpace(sequenceID) == "" {
		return fmt.Errorf("missing sequence ID")
	}

	outDir := targetDir
	if outDir == "" {
		outDir = "."
	}

	manifest, err := hub.GetManifest()
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}

	var entry *t.HubEntry
	for _, e := range manifest.Entries {
		if e.ID == sequenceID {
			entry = &e
			break
		}
	}
	if entry == nil {
		return fmt.Errorf("sequence not found: %s", sequenceID)
	}

	seqDir := filepath.Join(outDir, entry.Name)
	if err := os.MkdirAll(seqDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	fmt.Printf("Preparing download package: %s\n", entry.Name)
	fmt.Printf("Destination: %s\n", seqDir)
	fmt.Println("Dependencies:")
	if len(entry.Dependencies) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, d := range entry.Dependencies {
			fmt.Printf("  - %s (%s)\n", d.Name, d.Type)
		}
	}
	fmt.Println()

	// inline helper for progress bar
	writeProgress := func(name string, total int64, reader io.Reader, dstPath string) error {
		outFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		bar := func(percent float64) string {
			width := 25
			filled := int(percent / 4)
			if filled > width {
				filled = width
			}
			return fmt.Sprintf("[%s%s]", strings.Repeat("=", filled), strings.Repeat(" ", width-filled))
		}

		pr := io.TeeReader(reader, outFile)
		buf := make([]byte, 32*1024)
		var written int64

		for {
			n, err := pr.Read(buf)
			if n > 0 {
				written += int64(n)
				percent := float64(written) / float64(total) * 100
				if total <= 0 {
					percent = 100
				}
				fmt.Printf("\r  ↳ %-18s %s %3.0f%%", name, bar(percent), percent)
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
		}

		fmt.Printf("\r  ✓ %-18s (%d KB)\n", name, written/1024)
		return nil
	}

	for _, dep := range entry.Dependencies {
		var depPath string
		if dep.Type == t.HubDependencyTypeBackground {
			depPath = filepath.Join(seqDir, dep.Name+".wav")
		} else {
			depPath = filepath.Join(seqDir, dep.Name+".spsq")
		}

		if _, err := os.Stat(depPath); err == nil {
			fmt.Printf("Skipping existing file: %s\n", filepath.Base(depPath))
			continue
		}

		resp, err := http.Get(dep.DownloadUrl)
		if err != nil {
			return fmt.Errorf("failed to download %s: %v", dep.Name, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to download %s: received status code %d", dep.Name, resp.StatusCode)
		}

		if err = writeProgress(dep.Name, resp.ContentLength, resp.Body, depPath); err != nil {
			return fmt.Errorf("error saving %s: %v", dep.Name, err)
		}
	}

	resp, err := http.Get(entry.DownloadUrl)
	if err != nil {
		return fmt.Errorf("failed to download sequence %s: %v", entry.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download sequence %s: received status code %d", entry.Name, resp.StatusCode)
	}

	sequencePath := filepath.Join(seqDir, entry.Name+".spsq")
	if err = writeProgress(entry.Name, resp.ContentLength, resp.Body, sequencePath); err != nil {
		return fmt.Errorf("error saving sequence %s: %v", entry.Name, err)
	}

	// Track the download event
	hub.TrackDownload(entry.ID)

	fmt.Printf("\nAll files successfully saved to: %s\n", seqDir)
	return nil
}
