//go:build !nohub

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/tabwriter"

	synapseq "github.com/ruanklein/synapseq/v3/core"
	"github.com/ruanklein/synapseq/v3/internal/hub"
	s "github.com/ruanklein/synapseq/v3/internal/shared"
	t "github.com/ruanklein/synapseq/v3/internal/types"
)

// hubRunUpdate updates the local Hub manifest
func hubRunUpdate(quiet bool) error {
	if err := hub.HubUpdate(); err != nil {
		return fmt.Errorf("failed to update hub. Error\n  %v", err)
	}
	manifest, err := hub.GetManifest()
	if err != nil {
		return fmt.Errorf("failed to get hub manifest. Error\n  %v", err)
	}
	if !quiet {
		fmt.Printf("Fetched %d entries from the Hub. Last update: %s\n", len(manifest.Entries), manifest.LastUpdated)
	}
	return nil
}

// hubRunClean cleans the local Hub cache
func hubRunClean(quiet bool) error {
	if err := hub.HubClean(); err != nil {
		return fmt.Errorf("failed to clean hub cache. Error\n  %v", err)
	}
	if !quiet {
		fmt.Println("Hub cache cleaned successfully.")
	}
	return nil
}

// hubRunGet retrieves and processes a sequence from the Hub
func hubRunGet(sequenceId, outputFile string, quiet bool) error {
	var wg sync.WaitGroup

	entry, err := hub.HubGet(sequenceId)
	if err != nil {
		return fmt.Errorf("failed to load hub manifest. Error\n  %v", err)
	}
	if entry == nil {
		return fmt.Errorf("sequence not found in hub: %s", sequenceId)
	}

	inputFile, err := hub.HubDownload(entry, t.HubActionTrackingGet, &wg)
	if err != nil {
		return fmt.Errorf("failed to download sequence from hub. Error\n  %v", err)
	}

	if outputFile == "" {
		outputFile = entry.Name + ".wav"
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

	wg.Wait()
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
func hubRunDownload(sequenceID, targetDir string, quiet bool) error {
	var wg sync.WaitGroup

	if strings.TrimSpace(sequenceID) == "" {
		return fmt.Errorf("missing sequence ID")
	}

	if targetDir == "" || targetDir == "." {
		var err error
		targetDir, err = os.Getwd()
		if err != nil {
			return err
		}
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

	seqFile, err := hub.HubDownload(entry, t.HubActionTrackingDownload, &wg)
	if err != nil {
		return fmt.Errorf("failed to download sequence from hub. Error\n  %v", err)
	}

	if err := s.CopyDir(filepath.Dir(seqFile), filepath.Join(targetDir, entry.Name)); err != nil {
		return fmt.Errorf("failed to copy files to target directory. Error\n  %v", err)
	}

	wg.Wait()
	if !quiet {
		fmt.Printf("Sequence %q and its dependencies have been downloaded to %s\n", entry.Name, targetDir)
	}

	return nil
}

// hubRunInfo shows information about a sequence from the Hub
func hubRunInfo(sequenceID string) error {
	var wg sync.WaitGroup

	if strings.TrimSpace(sequenceID) == "" {
		return fmt.Errorf("missing sequence ID")
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

	seqFile, err := hub.HubDownload(entry, t.HubActionTrackingInfo, &wg)
	if err != nil {
		return fmt.Errorf("failed to download sequence from hub. Error\n  %v", err)
	}

	appCtx, err := synapseq.NewAppContext(seqFile, "", "text")
	if err != nil {
		return fmt.Errorf("failed to create application context. Error\n  %v", err)
	}

	if err := appCtx.LoadSequence(); err != nil {
		return fmt.Errorf("failed to load sequence. Error\n  %v", err)
	}

	fmt.Printf("Name:        %s\n", entry.Name)
	fmt.Printf("Author:      %s\n", entry.Author)
	fmt.Printf("Category:    %s\n", entry.Category)
	fmt.Printf("Updated At:  %s\n", entry.UpdatedAt[:10])

	dependencies := "\nDependencies: None\n"
	if len(entry.Dependencies) > 0 {
		dependencies = "\nDependencies:\n"
		for _, dep := range entry.Dependencies {
			dependencies += fmt.Sprintf("  - %s (%s)\n", dep.Name, dep.Type.String())
		}
	}
	fmt.Printf("%s", dependencies)

	description := "\nDescription: No description available.\n"
	comments := appCtx.Comments()
	if len(comments) > 0 {
		description = "\nDescription:\n"
		for _, comment := range comments {
			description += fmt.Sprintf("  %s\n", comment)
		}
	}
	fmt.Printf("%s", description)

	wg.Wait()
	return nil
}
