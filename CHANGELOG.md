# Changelog

All notable changes to this project will be documented in this file.

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
