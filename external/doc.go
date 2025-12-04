/*
Package external provides integrations with external audio utilities
such as ffplay and ffmpeg to extend the capabilities of the SynapSeq
engine without introducing additional internal complexity.

# Overview

The package uses SynapSeq's real-time PCM streaming (`AppContext.Stream`)
and sends it directly to the stdin of ffmpeg/ffplay. This avoids
temporary files, reduces memory usage, and maintains instant startup.

# External Utilities

The following utilities are supported:

  - ffplay – real-time playback of the SynapSeq-generated audio
  - ffmpeg – MP3 encoding from streamed PCM input

Custom paths may be provided when constructing FFplay or FFmpeg.
If no path is given, the package attempts to locate the utilities
using the system PATH.

# Example: Real-Time Playback

This example shows how to play a SynapSeq sequence directly through
ffplay using streaming PCM audio.

	package main

	import (
	    "log"
	    "os"

	    synapseq "github.com/ruanklein/synapseq/v3/core"
	    "github.com/ruanklein/synapseq/v3/external"
	)

	func main() {
	    // Create application context
	    ctx, err := synapseq.NewAppContext("input.spsq", "", "text")
	    if err != nil {
	        log.Fatal(err)
	    }

	    // Load sequence
	    if err := ctx.LoadSequence(); err != nil {
	        log.Fatal(err)
	    }

	    // Create ffplay utility (uses PATH by default)
	    player, err := external.NewFFPlay("")
	    if err != nil {
	        log.Fatal(err)
	    }

	    // Play audio in real time
	    if err := player.Play(ctx); err != nil {
	        log.Fatal(err)
	    }
	}

# Example: MP3 Encoding (VBR Mode)

This example streams PCM audio to ffmpeg and saves it as an MP3 file
using Variable Bit Rate (VBR) encoding at highest quality (V0).

	package main

	import (
	    "log"

	    synapseq "github.com/ruanklein/synapseq/v3/core"
	    "github.com/ruanklein/synapseq/v3/external"
	)

	func main() {
	    // Create application context for MP3 output
	    ctx, err := synapseq.NewAppContext("input.spsq", "output.mp3", "text")
	    if err != nil {
	        log.Fatal(err)
	    }

	    // Load sequence
	    if err := ctx.LoadSequence(); err != nil {
	        log.Fatal(err)
	    }

	    // Create ffmpeg converter
	    encoder, err := external.NewFFmpeg("")
	    if err != nil {
	        log.Fatal(err)
	    }

	    // Encode MP3 using VBR (default - highest quality)
	    if err := encoder.MP3(ctx, nil); err != nil {
	        log.Fatal(err)
	    }
	}

# Example: MP3 Encoding (CBR Mode)

This example encodes MP3 using Constant Bit Rate (CBR) at 320 kbps
for applications requiring fixed bitrate output.

	package main

	import (
	    "log"

	    synapseq "github.com/ruanklein/synapseq/v3/core"
	    "github.com/ruanklein/synapseq/v3/external"
	)

	func main() {
	    // Create application context for MP3 output
	    ctx, err := synapseq.NewAppContext("input.spsq", "output.mp3", "text")
	    if err != nil {
	        log.Fatal(err)
	    }

	    // Load sequence
	    if err := ctx.LoadSequence(); err != nil {
	        log.Fatal(err)
	    }

	    // Create ffmpeg converter
	    encoder, err := external.NewFFmpeg("")
	    if err != nil {
	        log.Fatal(err)
	    }

	    // Encode MP3 using CBR at 320 kbps
	    if err := encoder.MP3(ctx, &external.MP3Options{Mode: external.MP3ModeCBR}); err != nil {
	        log.Fatal(err)
	    }
	}

# MP3 Encoding Modes

The MP3 method supports two encoding modes:

  - VBR (Variable Bit Rate): Uses LAME V0 quality preset, providing the best
    quality-to-size ratio. This is the default mode when options is nil.

  - CBR (Constant Bit Rate): Uses fixed 320 kbps bitrate, useful for compatibility
    with older players or when consistent file size is required.

# Error Handling

If an external tool does not exist or is not executable, constructors
(NewFFPlay, NewFFmpeg) return an error. If the tool exits with a non-zero
status code, the returned error contains both streaming and command errors.

# Platform Notes

  - On Linux/macOS, executable permission bits are checked.
  - On Windows, lookups rely on PATH and associated .exe resolution.
  - Streaming uses stdin pipes and does not rely on temporary files.

# More Information

Full documentation and examples are available at:
https://github.com/ruanklein/synapseq
*/
package external
