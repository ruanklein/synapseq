# SynapSeq

### Synapse-Sequenced Brainwave Generator

SynapSeq is a lightweight engine for sequencing audio tones for brainwave entrainment, using a simple text-based format. It helps induce states such as relaxation, meditation, and focused awareness by guiding brainwave frequencies through sound.

---

🌐 **[Visit the official website](https://ruanklein.github.io/synapseq/)** for interactive examples, audio demonstrations, and complete documentation about the different types of brainwave entrainment.

---

## Table of Contents

- [Example](#example)
- [Installation](#installation)
- [Compilation](#compilation)
  - [For macOS ](#for-macos)
  - [For Linux](#for-linux)
  - [For Windows](#for-windows)
- [Documentation](#documentation)
- [License](#license)
- [Contact](#contact)
- [Credits](#credits)

## Example

Save this file as `relax.spsq`:

```
# Presets
relax1
  noise brown amplitude 40
  tone 250 binaural 10.0 amplitude 10
relax2
  noise brown amplitude 40
  tone 250 binaural 5.0 amplitude 10

# Timeline sequence
00:00:00 silence
00:00:15 relax1
00:02:00 relax1
00:03:00 relax2
00:04:00 relax2
00:05:00 relax1
00:06:00 relax1
00:07:00 relax2
00:08:00 relax2
00:09:00 relax1
00:10:00 silence
```

Run SynapSeq to generate the audio file:

```bash
synapseq --output relax.wav relax.spsq # Save the audio file
synapseq --output - relax.spsq | play - # Or play the audio file directly
```

The audio file will be created in the current directory.

When processing this file, SynapSeq will execute the following sequence of phases:

```
Phases:
├─ 0:00-0:15: Fade-in of silence for relax1 (start of the sequence)
├─ 0:15-2:00: relax1 (10Hz) - Brown noise + binaural tone
├─ 2:00-3:00: Ramp: 10Hz -> 5Hz (relax1 -> relax2)
├─ 3:00-4:00: relax2 (5Hz) - Brown noise + binaural tone
├─ 4:00-5:00: Ramp: 5Hz -> 10Hz (relax2 -> relax1)
├─ 5:00-6:00: relax1 (10Hz)
├─ 6:00-7:00: Ramp: 10Hz -> 5Hz (relax1 -> relax2)
├─ 7:00-8:00: relax2 (5Hz)
├─ 8:00-9:00: Ramp: 5Hz -> 10Hz (relax2 -> relax1)
├─ 9:00-10:00: Fade-out of relax1 for silence (end of the sequence)
```

## Installation

SynapSeq is a command-line engine, not a traditional desktop application. It’s designed to be compiled and used directly from the terminal, as part of your audio workflow.

## Compilation

SynapSeq can be compiled for macOS, Linux and Windows.

### For macOS

Install the "Xcode Command Line Tools" in your system.

```bash
xcode-select --install
```

Install [homebrew](https://brew.sh/) if you don't have it yet.

Install dependencies:

```bash
brew install pkg-config libvorbis libmad
```

Run the build script to create the binary:

```bash
./build/macos-build-synapseq.sh
```

The binary will be created in the `build/dist` folder.

To install the binary, run the following command:

```bash
sudo cp build/dist/synapseq-macos-arm64 /usr/local/bin/synapseq
```

### For Linux

In Ubuntu/Debian based distributions, install dependencies:

```bash
sudo apt-get install build-essential pkg-config libvorbis-dev libogg-dev libmad0-dev
```

Run the build script to create the binary:

```bash
./build/linux-build-synapseq.sh
```

The binary will be created in the `build/dist` folder.

To install the binary, run the following command:

```bash
sudo cp build/dist/synapseq-linux-x86_64 /usr/local/bin/synapseq
```

### For Windows

In Windows, the best way to build SynapSeq is using [Docker](https://www.docker.com/) with WSL2.

1. Install [WSL2](https://learn.microsoft.com/en-us/windows/wsl/install)
2. Install [Docker Desktop](https://www.docker.com/products/docker-desktop/)

Run this sequence of commands to build SynapSeq (only x86-64):

```bash
docker compose -f build/compose.yml up build-windows-libs-x86-64
docker compose -f build/compose.yml up build-windows-synapseq-x86-64
```

The `.exe` will be created in the `build/dist` folder.

## Documentation

For detailed information on all features, see the [USAGE.md](docs/USAGE.md) file.

## License

SynapSeq is distributed under the GPL license. See the [COPYING.txt](COPYING.txt) file for details.

## Contact

If you have any questions, please open a topic on the [discussions](https://github.com/ruanklein/synapseq/discussions) page.

## Credits

SynapSeq is based on the SBaGen. See [SBaGen project](https://uazu.net/sbagen/).
