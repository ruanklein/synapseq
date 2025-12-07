/*
Package external provides integrations with external audio utilities
such as ffplay, ffmpeg, and ffprobe to extend the capabilities of the SynapSeq
engine without introducing additional internal complexity.

# Overview

The package uses SynapSeq's real-time PCM streaming (`AppContext.Stream`)
and sends it directly to the stdin of ffmpeg/ffplay. This avoids
temporary files, reduces memory usage, and maintains instant startup.

# External Utilities

The following utilities are supported:

  - ffplay – real-time playback of the SynapSeq-generated audio
  - ffmpeg – audio encoding (MP3, OGG, OPUS) from streamed PCM input
  - ffprobe – metadata extraction from encoded audio files

Custom paths may be provided when constructing FFplay, FFmpeg, or FFprobe.
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
	    if err := encoder.Convert(ctx, "mp3", nil); err != nil {
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
	    opts := &external.CodecOptions{
	        MP3Options: &external.MP3Options{Mode: external.MP3ModeCBR},
	    }
	    if err := encoder.Convert(ctx, "mp3", opts); err != nil {
	        log.Fatal(err)
	    }
	}

# Example: OGG/Vorbis Encoding

This example streams PCM audio to ffmpeg and saves it as an OGG file
using Vorbis codec at highest quality.

		package main

		import (
		    "log"

		    synapseq "github.com/ruanklein/synapseq/v3/core"
		    "github.com/ruanklein/synapseq/v3/external"
		)

		func main() {
		    // Create application context for OGG output
		    ctx, err := synapseq.NewAppContext("input.spsq", "output.ogg", "text")
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

	    // Encode OGG/Vorbis at highest quality (options can be nil for OGG)
	    if err := encoder.Convert(ctx, "ogg", nil); err != nil {
	        log.Fatal(err)
	    }
	}

# Example: OPUS EncodingThis example streams PCM audio to ffmpeg and saves it as an OPUS file.
Note that OPUS encoding requires a sample rate of 48000 Hz.

		package main

		import (
		    "log"

		    synapseq "github.com/ruanklein/synapseq/v3/core"
		    "github.com/ruanklein/synapseq/v3/external"
		)

		func main() {
		    // Create application context for OPUS output
		    ctx, err := synapseq.NewAppContext("input.spsq", "output.opus", "text")
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

	    // Encode OPUS at 96 kbps (options can be nil for OPUS)
	    if err := encoder.Convert(ctx, "opus", nil); err != nil {
	        log.Fatal(err)
	    }
	}

# Example: Metadata ExtractionThis example extracts the original SynapSeq text sequence from an encoded
audio file (MP3, OGG, or OPUS) using ffprobe.

	package main

	import (
	    "fmt"
	    "log"

	    "github.com/ruanklein/synapseq/v3/external"
	)

	func main() {
	    // Create ffprobe instance
	    probe, err := external.NewFFprobe("")
	    if err != nil {
	        log.Fatal(err)
	    }

	    // Extract text sequence from encoded file
	    content, err := probe.ExtractTextSequence("output.mp3")
	    if err != nil {
	        log.Fatal(err)
	    }

	    fmt.Println(content)
	}

# Example: Save Extracted Metadata

This example extracts and saves the original SynapSeq text sequence
from an encoded audio file to a new file.

	package main

	import (
	    "log"

	    "github.com/ruanklein/synapseq/v3/external"
	)

	func main() {
	    // Create ffprobe instance
	    probe, err := external.NewFFprobe("")
	    if err != nil {
	        log.Fatal(err)
	    }

	    // Extract and save text sequence
	    if err := probe.SaveExtractedTextSequence("output.mp3", "extracted.spsq"); err != nil {
	        log.Fatal(err)
	    }
	}

# Audio Format Support

The Convert method supports the following audio formats:

  - MP3: Uses libmp3lame encoder with two modes (mode selection is required):

  - VBR (Variable Bit Rate): LAME V0 quality preset, best quality-to-size ratio

  - CBR (Constant Bit Rate): Fixed 320 kbps bitrate, useful for compatibility

  - OGG: Uses libvorbis encoder at highest quality (q:a 10)

  - OPUS: Uses libopus encoder at 96 kbps (requires 48000 Hz sample rate)

# Metadata Embedding and Extraction

When encoding to MP3, OGG, or OPUS formats, SynapSeq automatically embeds
metadata into the output file, including:

  - synapseq_id: Unique identifier for the sequence
  - synapseq_generated: Generation timestamp
  - synapseq_version: SynapSeq version used
  - synapseq_platform: Platform information
  - synapseq_content: Base64-encoded original sequence content

This metadata can be extracted later using FFprobe, allowing full recovery
of the original sequence definition from encoded audio files.

To disable metadata embedding in output files, you can create a context copy
using the WithUnsafeNoMetadata() method:

	// Create context and disable metadata embedding
	ctx, err := synapseq.NewAppContext("input.spsq", "output.mp3", "text")
	if err != nil {
	    log.Fatal(err)
	}

	// Disable metadata embedding (text format only)
	ctx, err = ctx.WithUnsafeNoMetadata()
	if err != nil {
	    log.Fatal(err)
	}

	// Load and encode without metadata
	if err := ctx.LoadSequence(); err != nil {
	    log.Fatal(err)
	}

	encoder, err := external.NewFFmpeg("")
	if err != nil {
	    log.Fatal(err)
	}

	opts := &external.CodecOptions{
	    MP3Options: &external.MP3Options{Mode: external.MP3ModeVBR},
	}
	if err := encoder.Convert(ctx, "mp3", opts); err != nil {
	    log.Fatal(err)
	}

# Error Handling

If an external tool does not exist or is not executable, constructors
(NewFFPlay, NewFFmpeg, NewFFprobe) return an error. If the tool exits with
a non-zero status code, the returned error contains both streaming and
command errors.

# Platform Notes

  - On Linux/macOS, executable permission bits are checked.
  - On Windows, lookups rely on PATH and associated .exe resolution.
  - Streaming uses stdin pipes and does not rely on temporary files.

# More Information

Full documentation and examples are available at:
https://github.com/ruanklein/synapseq
*/
package external
