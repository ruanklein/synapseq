# SynapSeq v3.0.0 Roadmap

This document details the steps required to port the SynapSeq project from C to Go for v3, based on a deep analysis of the original `synapseq.c` source code. The goal is to maintain all core features, improve maintainability, and provide a modern, idiomatic Go codebase.

---

### Project Guidelines

**Standard Library First:**  
The goal is to use Go's standard library as much as possible. External or third-party packages may be included only if they provide significant value to the project and will be vendored (included in the repository) to ensure long-term stability and reproducibility. All parsing, audio processing, file I/O, and CLI functionality should preferably be implemented using native Go packages.

**Explicit Naming:**  
All function, method, variable, struct, and interface names must be explicit and descriptive. Avoid abbreviations. For example, do not port a function as `corrVal`; instead, use a full, meaningful name such as `CorrectParameterValue`. Every identifier should clearly express its purpose and usage.

**Do Not Migrate Dead Code:**  
Do not migrate code that is obsolete or irrelevant for SynapSeq, such as remnants of real-time playback or any functionality not aligned with the current project goals. Only port code that is meaningful and necessary for the offline sequence-to-WAV workflow.

**Code Formatting:**  
All Go code must be formatted using [`gofmt`](https://pkg.go.dev/cmd/gofmt). Please ensure your code is properly formatted before submitting a pull request.

---

### Stage Descriptions

1. **Project Structure & Planning:**  
   Define the Go module and repository structure, modularize the codebase, and set up version control and documentation.

2. **Data Structures:**  
   Translate C structs and constants into Go idioms, using slices, pointers, and appropriate types.

3. **Parsing & Sequence File Handling:**  
   Implement parsing of `.spsq` files, handling comments, options, presets, timeline entries, and validation.

4. **Timeline Construction & Validation:**  
   Build and validate the sequence of periods, handle time parsing, insert transitional periods, and normalize amplitudes.

5. **Audio Synthesis Engine:**  
   Port all audio generation logic, including tones, noise, effects, and parameter interpolation.

6. **Buffering and Background Audio:**  
   Replace C's threaded buffer logic with Go goroutines and channels, and handle background WAV file reading and mixing.

7. **Main Processing Loop:**  
   Implement the main loop for sample generation, parameter interpolation, mixing, and buffer management.

8. **WAV File Output:**  
   Write the WAV header and PCM data using Go's binary utilities, supporting output to file or stdout.

9. **CLI and Library Interface:**  
   Design a CLI for command-line usage and expose the core as a Go package for library use.

10. **Testing & Validation:**  
    Write unit and integration tests, validate output, and test edge cases.

11. **Documentation:**  
    Document all modules, provide usage examples, and write a migration guide.

12. **Packaging & Release:**  
    Prepare build scripts, ensure cross-platform compatibility, and publish the project.

---

## 1. Project Structure & Planning

- Define the Go module and repository structure.
- Plan for modular packages: audio synthesis, sequence parsing, timeline management, WAV output, CLI, and utilities.
- Set up version control, initial documentation, and continuous integration.

---

## 2. Data Structures

- Translate C structs (`Voice`, `Channel`, `Period`, `NameDef`) into Go structs.
- Use slices and pointers for dynamic lists (e.g., periods, voices).
- Replace C macros and constants with Go `const` and `var`.

---

## 3. Parsing & Sequence File Handling

- Implement a parser for `.spsq` files using Go's `bufio.Scanner` and string utilities.
- Port logic for:
  - Skipping comments and blank lines.
  - Handling options (`@background`, `@gainlevel`, etc.).
  - Parsing name definitions (presets) and timeline entries.
  - Validating names and syntax.
- Ensure error handling is idiomatic (using Go errors, not `exit()`).

---

## 4. Timeline Construction & Validation

- Recreate the logic of `readTimeLine`, `readNameDef`, and `correctPeriods`:
  - Build a doubly-linked or slice-based list of `Period` structs.
  - Implement time parsing and validation (start/end times, chronological order).
  - Handle automatic insertion of transitional periods.
  - Implement amplitude normalization and validation as in `normalizeAmplitude`.
  - Remove redundant or zero-length periods after validation.

---

## 5. Audio Synthesis Engine

- Port all audio generation logic:
  - Tone generation (binaural, monaural, isochronic).
  - Noise generation (`noise2` for pink, `white_noise`, `brown_noise`).
  - Spin and effect logic (`create_noise_spin_effect`).
  - Interpolation (ramp/slide) between period values (as in `corrVal`).
- Use Go's math and random packages for calculations.
- Implement waveform tables (sine, square, triangle, sawtooth) as slices.

---

## 6. Buffering and Background Audio

- Replace C's threaded buffer logic (`inbuf_*`, `volatile` variables) with Go goroutines and channels.
- Implement a producer-consumer pattern for background audio mixing.
- Handle WAV file reading for background audio using Go's `os` and `encoding/binary` packages.

---

## 7. Main Processing Loop

- Port the main loop (`loop` and `outChunk`):
  - For each output buffer, interpolate parameters, generate samples, and mix channels.
  - Apply volume, dithering, and normalization.
  - Handle background audio mixing and looping.
  - Ensure correct buffer management and output chunking.

---

## 8. WAV File Output

- Implement WAV header and PCM data writing using Go's binary writing utilities.
- Ensure correct handling of sample rate, bit depth, and stereo channels.
- Support output to file or stdout as in the original.

---

## 9. CLI and Library Interface

- Design a CLI using Go's `flag` for command-line options.
- Expose core functionality as a Go package for use in other projects (not just CLI).
- Ensure clear separation between CLI and core logic.

---

## 10. Testing & Validation

- Write unit and integration tests for all modules.
- **Validate output WAV files against those generated by the original C version.**
- Test edge cases: invalid sequences, overlapping periods, extreme parameter values.

---

## 11. Documentation

- Document all exported functions, structs, and packages.
- Provide usage examples for both CLI and library usage.
- Write a migration guide for users familiar with the C version.

---

## 12. Packaging & Release

- Prepare build scripts and release instructions.
- Ensure cross-platform compatibility (Linux, macOS, Windows).
- Publish the project and documentation.

**This roadmap ensures a faithful, maintainable, and idiomatic Go port of SynapSeq, covering all technical and architectural aspects identified in the original C code.**
