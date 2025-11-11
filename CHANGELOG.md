# Changelog

All notable changes to this project will be documented in this file.

## [3.3.0]

### New Features

- **SynapSeq Hub Integration**: Added built-in support for the SynapSeq Hub, allowing users to easily browse, download, and generate community-contributed sequences directly from the CLI. New flags include:
  - `-hub-update`: Update local index of available sequences from the Hub.
  - `-hub-list`: List all available sequences in the Hub.
  - `-hub-get <sequence> <output>`: Download and generate a specified sequence from the Hub to the given output file.
  - `-hub-search <query>`: Search for sequences in the Hub matching the query term.
  - `-hub-download <sequence> <output>`: Download a sequence file from the Hub without generating audio.
  - `-hub-info <sequence>`: Display detailed information about a specific sequence in the Hub.
  - `-hub-clean`: Remove all locally cached Hub sequences.

The Hub integration includes anonymous usage tracking for downloads and generation actions to help understand community engagement. See [Usage Documentation](docs/USAGE.md) for details.

### Improvements

- Updated documentation to include SynapSeq Hub usage instructions and examples.
- Fixed an issue where `@presetlist` and `@background` references only worked when the CLI was executed from the same directory as the .spsq file. Now, these references are resolved relative to the location of the .spsq file, allowing more flexible organization of sequence files and assets.
- Fixed metadata embedding for sequences loaded from stdin or HTTP/HTTPS, ensuring full provenance tracking.
- Disabled metadata embedding for sequences using presetlist inheritance, since embedding would make the resulting WAV non-reproducible (breaking determinism).
- Simplified help output, focusing on core options and the new Hub commands.
- Removed legacy man references (now redundant with online documentation).
- Distribution archives are now slimmer, only the binary and checksum are included per platform.

### Maintenance

- Removed the legacy `contrib/` directory; all community works must now be submitted via **SynapSeq Hub**.
- Retained `samples/` for demonstration and testing purposes, providing a minimal set of local examples.
- Added internal Hub version matching (`HUB_VERSION`) to ensure compatibility between CLI and manifest schema.

## [3.2.2]

### Bug Fixes

- Removed "Press Enter to exit..." prompt on Windows to ensure smooth execution in CI/CD and non-interactive environments.

### Improvements

- Updated CLI help output to exclude manpage suggestion on Windows, as manpages are not typically available on that platform.
- stderr is now used for error messages, while stdout is reserved for regular output. This ensures proper redirection and logging of errors.
- -version flag includes build date and commit hash for better traceability of builds.

## [3.2.1]

### Fixes and Maintenance

- **Go Module Path Update:**  
   The module path has been updated from:

  ```
  github.com/ruanklein/synapseq
  ```

  to:

  ```
  github.com/ruanklein/synapseq/v3
  ```

  This change aligns with Go's [Semantic Import Versioning](https://go.dev/doc/modules/version-numbers) requirements for modules at version `v2` and higher.  
  Without this suffix, Go tools (such as `go get` and `pkg.go.dev`) would incorrectly fetch the legacy `v2.1.1+incompatible` release instead of the latest version.

- **Fixed Import Paths:**  
  All internal and public imports have been updated to reference the correct v3 module path.

## [3.2.0]

### New Features

- **Preset Inheritance**: Presets can now inherit from other presets, allowing more reusable and modular sequence definitions. This enables fine-grained control over variations while reducing repetition in `.spsq` files.
- **Embedded Metadata in WAV Files**: SynapSeq now embeds the original sequence and metadata into the generated WAV files. You can extract this information at any time using the new `-extract` CLI flag
- **Reverse Conversion**: Added CLI support to convert structured formats **(JSON, XML, YAML)** back into `.spsq` format using the `-convert` flag. This enables round-trip conversion and easier editing of machine-generated sequences.
- **Syntax Test Mode**: Added the `-test` CLI flag to validate the syntax of `.spsq` and structured formats (JSON, XML, YAML) without generating output. This replaces the previous `-debug` flag, which has been removed.
- **Core Library Extraction**: The internal sequencing engine has been modularized into a standalone Go library, enabling integration with third-party tools, GUIs, and services without relying on CLI execution.
- **Precompiled Binaries and CI/CD Build Automation**: SynapSeq is now distributed with precompiled binaries for all major platforms:
  - **Windows** (`.zip` for x86_64 and ARM64)
  - **macOS** (`.tar.gz` for ARM64)
  - **Linux** (`.tar.gz` for x86_64 and ARM64)

These binaries are published on GitHub Releases and allow anyone to run SynapSeq without needing to install Go or Make.

### Improvements

- Updated CLI documentation to reflect new `-convert` usage.
- The `-debug` CLI flag was removed and replaced by `-test`, which validates the syntax of sequence and structured files without generating output.
- Improved internal code organization by separating core logic and CLI wrapper.

### Notes

- Signed binary support is not included. On Windows/macOS, you may need to confirm execution due to unsigned status.

## [3.1.0]

### New Features

- Added transition types: `steady`, `ease-in`, `ease-out`, and `smooth` for more natural brainwave entrainment progressions.
- Added support for structured formats: JSON, XML, and YAML via `-json`, `-xml`, and `-yaml` CLI flags.
- Added `@presetlist` global option, allowing preset reuse across different sequences.
- Added support for HTTP/HTTPS URLs as input:
  - CLI: load sequences directly from web URLs
  - `@background`: load background audio from web URLs
  - `@presetlist`: load preset files from web URLs
- Added new sample files demonstrating all new features:
  - Transition examples:
    - [sample-transitions.spsq](samples/sample-transitions.spsq)
  - Preset reuse examples using `@presetlist`:
    - [presets-focus.spsq](samples/presets-focus.spsq) - Shared focus presets
    - [presets-relax.spsq](samples/presets-relax.spsq) - Shared relaxation presets
    - [sample-focus-one.spsq](samples/sample-focus-one.spsq) - Focus session part one
    - [sample-focus-two.spsq](samples/sample-focus-two.spsq) - Focus session part two
    - [sample-relax-one.spsq](samples/sample-relax-one.spsq) - Relaxation session part one
    - [sample-relax-two.spsq](samples/sample-relax-two.spsq) - Relaxation session part two
    - [sample-genesis.spsq](samples/sample-genesis.spsq) - Complete session using both preset files
  - Structured format samples:
    - [samples/structured/](samples/structured/) - JSON, XML, and YAML examples

### Bug Fixes

- Fixed audio clipping issue in `@background` option for better audio quality.
- Fixed `-quiet` option to properly suppress non-error output.

### Improvements

- Improved status output during audio generation for better tracking of sequence progress.
- Added file size limits for all file types to prevent memory overflow:
  - Text format (`.spsq`): 32 KB max
  - Structured formats (JSON/XML/YAML): 128 KB max
  - Background audio (`.wav`): 10 MB max
- Added `man` target to [Makefile](Makefile) for automatic man page generation from documentation, facilitating offline documentation access.
- Added Content-Type validation for HTTP/HTTPS file loading.

## [3.0.1]

- Replaced audio dependency: migrated from go-audio to gopxl/beep v2 for WAV encoding/decoding and streaming.
  - Updated RenderWav to use beep/wav encoder.
  - Updated background WAV decoder to beep/wav and adapted internal read loop.
  - Reworked Render to provide interleaved int24 samples directly (no go-audio types).
  - Updated RenderRaw and all audio tests to drop go-audio completely.
- Licensing note: go-audio’s Apache-2.0 license is not compatible with SynapSeq’s GPLv2-only licensing. The new dependency (beep) uses a permissive license compatible with GPLv2.
- The 3.0.0 release/tag was removed from the repository to prevent distributing artifacts built with the previous, incompatible dependency set. Version 3.0.1 supersedes 3.0.0.
- Added "pure" (no beats) tones.

## [3.0.0]

- Fully rewritten from scratch in Go (previous versions were based on C code).
- No longer based on or forked from SBaGen; now only inspired by its workflow.
- Unified syntax for backgrounds and effects (spin and pulse).
- Embedded spin noise removed; now handled via unified effects.
- Simplified build process using Makefile; removed shell scripts and C dependencies.
- Fade-in and fade-out rules are now explicit; there is no longer automatic fading between different tracks.

## [2.1.1]

- Fixed build scripts.
- Added `--version` option.

## [2.1.0]

- Removed playback real-time.
- Removed verbose mode.
- Removed support for ALSA, CoreAudio and WIN32 audio API.
- Removed `--wav` option.
- Removed `--output-raw-file` option.
- Removed `--output-wav-file` option.
- Removed `--verbose` option.
- Removed `--buffer-size` option (macOS only).
- Removed `--device` option (Linux only).
- Removed setup installer scripts.
- Removed assets, icons and background image.
- Removed `@verbose` global option.
- Removed `@test` global option.
- Removed `@quiet` global option.
- Removed `@waveform` global option.
- Added `--output` option.
- WAV output is now the default.

## [2.0.0]

- New, streamlined syntax. Simpler and more user-oriented.
- Removed SBaGen old features.
- `libvorbisidec` replaced by `libvorbis`.
- Command line options are now more consistent with the syntax.
- Background sounds with looping support.
- Gain level control for background sounds.
- New format for sequence files: `.spsq`.
- New icon and new background image for dmg.
- Added Monaural Beats
