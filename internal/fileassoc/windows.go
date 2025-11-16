//go:build windows

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package fileassoc

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

// Helper that recursively deletes a key tree
func deleteRegistryTree(base registry.Key, path string) error {
	k, err := registry.OpenKey(base, path, registry.ALL_ACCESS)
	if err != nil {
		return nil // If not exist, fine
	}
	defer k.Close()

	subKeys, err := k.ReadSubKeyNames(-1)
	if err == nil {
		for _, sub := range subKeys {
			_ = deleteRegistryTree(base, path+`\`+sub)
		}
	}

	err = registry.DeleteKey(base, path)
	if err != nil {
		return fmt.Errorf("failed to delete registry key %s: %w", path, err)
	}

	return nil
}

// CleanSynapSeqWindowsRegistry removes all SynapSeq-related registry keys.
// Safe to run even if some keys don't exist.
func CleanSynapSeqWindowsRegistry() error {
	extKeyPath := `Software\Classes\.spsq`

	extKey, err := registry.OpenKey(registry.CURRENT_USER, extKeyPath, registry.READ)
	if err == nil {
		defer extKey.Close()

		val, _, err := extKey.GetStringValue("")
		if err == nil && val == "SynapSeq.File" {
			registry.DeleteKey(registry.CURRENT_USER, extKeyPath)
		}
	}

	_ = deleteRegistryTree(registry.CURRENT_USER, `Software\Classes\SynapSeq.File`)
	_ = deleteRegistryTree(registry.CURRENT_USER,
		`Software\Classes\SystemFileAssociations\.wav\shell\SynapSeqExtract`)

	return nil
}

// InstallWindowsFileAssociation sets up the Windows registry to associate .spsq files with SynapSeq
func InstallWindowsFileAssociation() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exePath := filepath.Clean(exe)

	extKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Classes\.spsq`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer extKey.Close()

	if err := extKey.SetStringValue("", "SynapSeq.File"); err != nil {
		return err
	}

	progIDKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Classes\SynapSeq.File`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer progIDKey.Close()

	progIDKey.SetStringValue("", "SynapSeq Sequence File")

	iconKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Classes\SynapSeq.File\DefaultIcon`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer iconKey.Close()

	iconKey.SetStringValue("", exePath+",0")

	cmdKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Classes\SynapSeq.File\shell\open\command`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer cmdKey.Close()

	openCmd := fmt.Sprintf(`"%s" "%%1"`, exePath)
	cmdKey.SetStringValue("", openCmd)

	return nil
}

// InstallWindowsContextMenu adds SynapSeq options to the Windows context menu for .spsq files
func InstallWindowsContextMenu() error {
	base := `Software\Classes\SynapSeq.File\shell`

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get exe path: %w", err)
	}
	exePath := filepath.Clean(exe)

	// ===============================
	// Test Sequence
	// ===============================
	testKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		base+`\TestSequence`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create TestSequence menu: %w", err)
	}
	defer testKey.Close()

	testKey.SetStringValue("", "SynapSeq: Test sequence")
	testKey.SetStringValue("Icon", exePath+",0")

	testCmdKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		base+`\TestSequence\command`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create TestSequence command: %w", err)
	}
	defer testCmdKey.Close()

	testCmd := `cmd.exe /C synapseq -test "%1" & echo. & pause`
	testCmdKey.SetStringValue("", testCmd)

	// ===============================
	// Edit Sequence
	// ===============================
	editKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		base+`\EditSequence`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create EditSequence menu: %w", err)
	}
	defer editKey.Close()

	editKey.SetStringValue("", "SynapSeq: Edit sequence")
	editKey.SetStringValue("Icon", exePath+",0")

	editCmdKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		base+`\EditSequence\command`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create EditSequence command: %w", err)
	}
	defer editCmdKey.Close()

	editCmd := `notepad.exe "%1"`
	editCmdKey.SetStringValue("", editCmd)

	return nil
}

// InstallWindowsExtractMenu adds an "Extract sequence" option to the Windows context menu for .wav files
func InstallWindowsExtractMenu() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exePath := filepath.Clean(exe)

	base := `Software\Classes\SystemFileAssociations\.wav\shell\SynapSeqExtract`

	// Main key
	k, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		base,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create wav extract menu: %w", err)
	}
	defer k.Close()

	k.SetStringValue("", "SynapSeq: Extract sequence")
	k.SetStringValue("Icon", exePath+",0")

	// Command key
	cmdKey, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		base+`\command`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create extract command: %w", err)
	}
	defer cmdKey.Close()

	extractCmd := `cmd.exe /C synapseq -extract "%1" & echo. & pause`
	cmdKey.SetStringValue("", extractCmd)

	return nil
}
