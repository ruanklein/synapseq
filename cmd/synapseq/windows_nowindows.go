//go:build !windows

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package main

import (
	"fmt"
)

// installWindowsFileAssociation is disabled for non-Windows builds
func installWindowsFileAssociation() error {
	return fmt.Errorf("this build does not support Windows file association installation")
}

// uninstallWindowsFileAssociation is disabled for non-Windows builds
func uninstallWindowsFileAssociation() error {
	return fmt.Errorf("this build does not support Windows file association removal")
}
