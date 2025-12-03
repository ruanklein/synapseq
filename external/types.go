/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package external

import (
	"fmt"
	"os"
	"os/exec"
)

// externalTool holds paths to external Unix tools like ffmpeg and ffplay
type externalTool struct {
	utilityPath string
}

// newUtility creates a new externalTool instance with given utility path
func newUtility(utilPath string) (*externalTool, error) {
	if utilPath == "" {
		return nil, fmt.Errorf("utility path cannot be empty")
	}

	filePath, err := exec.LookPath(utilPath)
	if err == nil {
		return &externalTool{
			utilityPath: filePath,
		}, nil
	}

	fileInfo, err := os.Stat(utilPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("executable not found at custom path: %s", utilPath)
		}
		return nil, fmt.Errorf("error checking path: %s, error: %v", utilPath, err)
	}

	if fileInfo.Mode().IsRegular() && (fileInfo.Mode().Perm()&0111 != 0) {
		return &externalTool{
			utilityPath: utilPath,
		}, nil
	}

	return nil, fmt.Errorf("file at path is not executable: %s", utilPath)
}
