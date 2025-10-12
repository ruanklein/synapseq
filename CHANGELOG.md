# Changelog

All notable changes to this project will be documented in this file.

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
