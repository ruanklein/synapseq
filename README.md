# SynapSeq

**Synapse-Sequenced Brainwave Generator**

SynapSeq is a lightweight and efficient engine for sequencing audio tones for brainwave entrainment, using a simple text-based format. It helps induce states such as relaxation, meditation, and focused awareness by guiding brainwave frequencies through sound.

## Table of Contents

- [Quick Start Example](#quick-start-example)
- [Installation](#installation)
- [Compilation](#compilation)
  - [macOS](#macos)
  - [Linux](#linux)
  - [Windows](#windows)
- [Documentation](#documentation)
- [Contributing](#contributing)
  - [Code of Conduct](#code-of-conduct)
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
synapseq relax.spsq relax.wav
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

### More Examples

You can find additional example scripts in the `samples/` folder of this repository. These include various types of brainwave entrainment sequences that you can download and test:

- `sample-binaural.spsq` - Binaural beats example
- `sample-isochronic.spsq` - Isochronic tones example
- `sample-monaural.spsq` - Monaural beats example
- `sample-noise.spsq` - Noise-based entrainment
- `sample-background-spin.spsq` - "Spin" effect
- `sample-background-pulse.spsq` - "Pulse" effect
- `sample-waveform.spsq` - Custom waveform example

## Installation

SynapSeq is a command-line tool that needs to be compiled from source. Follow the instructions below for your operating system.

### Prerequisites

You need to install Go (v1.25 or later) and make on your system before compiling SynapSeq.

#### Installing Go

**macOS:**

```bash
# Using Homebrew
brew install go

# Using MacPorts
sudo port install go
```

**Linux (Ubuntu/Debian):**

```bash
# Update package list
sudo apt update

# Install Go
sudo apt install golang-go make

# Or install a newer version using snap
sudo snap install go --classic
```

**Linux (CentOS/RHEL/Fedora):**

```bash
# For Fedora
sudo dnf install golang make

# For CentOS/RHEL
sudo yum install golang make
```

**Windows:**

```powershell
# Using Chocolatey (install Chocolatey first from https://chocolatey.org/)
choco install golang make

# Using Scoop (install Scoop first from https://scoop.sh/)
scoop install go make

# Using winget (Windows 10/11)
winget install GoLang.Go
winget install GnuWin32.Make
```

**Verify installation:**

```bash
go version
make --version
```

## Compilation

SynapSeq can be compiled using the provided Makefile. For most users, simply run:

```bash
make
```

This will automatically compile SynapSeq for your current operating system and architecture, creating a binary in the `bin/` directory.

### Installing the Binary

After compilation, install the binary system-wide:

**macOS/Linux:**

```bash
sudo cp bin/synapseq /usr/local/bin/synapseq
```

**Windows:**

```cmd
# Run Command Prompt as Administrator
mkdir "C:\Program Files\SynapSeq"
copy "bin\synapseq.exe" "C:\Program Files\SynapSeq\synapseq.exe"
```

Then add `C:\Program Files\SynapSeq` to your PATH environment variable.

### Cross-Platform Compilation

If you need to build for a different platform, use these specific commands:

#### macOS

```bash
make build-macos
```

Creates: `bin/synapseq-macos-arm64`

#### Linux

```bash
make build-linux
```

Creates:

- `bin/synapseq-linux-amd64`
- `bin/synapseq-linux-arm64`

#### Windows

```bash
make build-windows
```

Creates:

- `bin/synapseq-windows-amd64.exe`
- `bin/synapseq-windows-arm64.exe`

### Additional Make Commands

- `make build` - Build for your current platform
- `make clean` - Remove all compiled binaries
- `make all` - Same as `make build`

## Documentation

For detailed information on all features and advanced usage, see the [USAGE.md](docs/USAGE.md) file.

## Contributing

We welcome contributions!

Please read the [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines on how to contribute code, bug fixes, documentation, or new sequence files to the project.

### Code of Conduct

Please note that all contributors are expected to follow our [Code of Conduct](CODE_OF_CONDUCT.md).

- Be respectful and considerate in all interactions.
- Harassment or abusive behavior will not be tolerated.
- Help us maintain a friendly and inclusive community.

If you experience or witness unacceptable behavior, please report it as described in the [Code of Conduct](CODE_OF_CONDUCT.md).

## License

SynapSeq is distributed under the GPL license. See the [COPYING.txt](COPYING.txt) file for details.

This repository includes vendored third-party libraries under the `vendor/` directory.  
Each vendored library retains its own license, which applies independently of the SynapSeq license.

## Contact

If you have any questions, please open a topic on the [discussions](https://github.com/ruanklein/synapseq/discussions) page.

## Credits

- **SBaGen** — SynapSeq was inspired by the [SBaGen project](https://uazu.net/sbagen/) (written in C) and follows a similar workflow.  
  SynapSeq has been completely rewritten from scratch in Go, but the conceptual foundation comes from SBaGen’s pioneering work in brainwave entrainment.

- **go-audio** — This project uses parts of the [go-audio](https://github.com/go-audio) libraries for audio encoding and decoding support, which provided a solid foundation for handling WAV data in Go.

We gratefully acknowledge these projects as the basis and inspiration for SynapSeq’s development.
