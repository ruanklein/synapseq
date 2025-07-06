# SynapSeq

**Synapse-Sequenced Brainwave Generator**

SynapSeq is a lightweight and efficient engine for sequencing audio tones for brainwave entrainment, using a simple text-based format. It helps induce states such as relaxation, meditation, and focused awareness by guiding brainwave frequencies through sound.

---

**[Visit the official website](https://ruanklein.github.io/synapseq/)** for interactive examples, audio demonstrations, and complete documentation about the different types of brainwave entrainment.

---

## Table of Contents

- [Quick Start Example](#quick-start-example)
- [Installation](#installation)
- [Compilation](#compilation)
  - [macOS](#macos)
  - [Linux](#linux)
  - [Windows](#windows)
- [Documentation](#documentation)
- [License](#license)
- [Contact](#contact)
- [Credits](#credits)

## Quick Start Example

Save the following content as `relax.spsq`:

```
# Presets
alpha
  noise brown amplitude 40
  tone 250 binaural 10.0 amplitude 10
theta
  noise brown amplitude 40
  tone 250 binaural 5.0 amplitude 10

# Timeline sequence
00:00:00 silence
00:00:15 alpha
00:02:00 alpha
00:03:00 theta
00:04:00 theta
00:05:00 alpha
00:06:00 alpha
00:07:00 theta
00:08:00 theta
00:09:00 alpha
00:10:00 silence
```

Run SynapSeq to generate the audio file:

```bash
synapseq --output relax.wav relax.spsq    # Save the audio file
synapseq --output - relax.spsq | play -   # Or play directly (UNIX only)
```

The audio file will be created in the current directory.

### Phase Sequence

When processing this file, SynapSeq will execute the following sequence of phases:

```
Phases:
├─ 0:00-0:15: Fade-in from silence to alpha (start of sequence)
├─ 0:15-2:00: alpha (10Hz) - Brown noise + binaural tone
├─ 2:00-3:00: Transition: 10Hz → 5Hz (alpha → theta)
├─ 3:00-4:00: theta (5Hz) - Brown noise + binaural tone
├─ 4:00-5:00: Transition: 5Hz → 10Hz (theta → alpha)
├─ 5:00-6:00: alpha (10Hz)
├─ 6:00-7:00: Transition: 10Hz → 5Hz (alpha → theta)
├─ 7:00-8:00: theta (5Hz)
├─ 8:00-9:00: Transition: 5Hz → 10Hz (theta → alpha)
└─ 9:00-10:00: Fade-out from alpha to silence (end of sequence)
```

## Installation

SynapSeq is a command-line tool, not a traditional desktop application. It's designed to be compiled and used directly from the terminal as part of your audio workflow.

## Compilation

SynapSeq can be compiled for macOS, Linux, and Windows.

### macOS

Install the Xcode Command Line Tools on your system:

```bash
xcode-select --install
```

Install [Homebrew](https://brew.sh/) if you don't have it yet.

Install the required dependencies:

```bash
brew install pkg-config libvorbis libmad
```

Run the build script to create the binary:

```bash
./build/macos-build-synapseq.sh
```

The binary will be created in the `build/dist` folder.

To install the binary system-wide:

```bash
sudo cp build/dist/synapseq-macos-arm64 /usr/local/bin/synapseq
```

### Linux

On Ubuntu/Debian-based distributions, install the dependencies:

```bash
sudo apt-get install build-essential pkg-config libvorbis-dev libogg-dev libmad0-dev
```

Run the build script to create the binary:

```bash
./build/linux-build-synapseq.sh # Build a dynamic binary (recommended)
./build/linux-build-synapseq-static.sh # Or build a static binary (optional)
```

The binary will be created in the `build/dist` folder.

To install the binary system-wide:

```bash
sudo cp build/dist/synapseq-linux-x86_64 /usr/local/bin/synapseq # For x86_64
sudo cp build/dist/synapseq-linux-arm64 /usr/local/bin/synapseq # For arm64
```

### Windows

On Windows, the recommended way to build SynapSeq is using [Docker](https://www.docker.com/) with WSL2.

**Prerequisites:**

1. Install [WSL2](https://learn.microsoft.com/en-us/windows/wsl/install)
2. Install [Docker Desktop](https://www.docker.com/products/docker-desktop/)

Run this sequence of commands to build SynapSeq (x86-64 only):

```bash
docker compose -f build/compose.yml up build-windows-libs-x86-64
docker compose -f build/compose.yml up build-windows-synapseq-x86-64
```

The `.exe` file will be created in the `build/dist` folder.

## Documentation

For detailed information on all features and advanced usage, see the [USAGE.md](docs/USAGE.md) file.

## License

SynapSeq is distributed under the GPL license. See the [COPYING.txt](COPYING.txt) file for details.

## Contact

If you have any questions, please open a topic on the [discussions](https://github.com/ruanklein/synapseq/discussions) page.

## Credits

SynapSeq is based on SBaGen. See the [SBaGen project](https://uazu.net/sbagen/) for more information.
