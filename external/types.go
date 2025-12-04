/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package external

import "os/exec"

// baseUtility represents a base external utility
type baseUtility struct{ path string }

// Path returns the path of the external utility
func (bu *baseUtility) Path() string {
	return bu.path
}

// Command creates an exec.Cmd for the utility with given arguments
func (bu *baseUtility) Command(args ...string) *exec.Cmd {
	return exec.Command(bu.path, args...)
}
