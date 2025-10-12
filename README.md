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
- [FAQ](#faq)
- [Contributing](#contributing)
  - [Code of Conduct](#code-of-conduct)
- [License](#license)
  - [Third-Party License](#third-party-licenses)
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

You can find additional example scripts in the `samples/` folder of this repository. See the [samples/README.md](samples/README.md) for detailed information about each example.

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

For Windows users, we recommend using **Git Bash** or **WSL2** (Windows Subsystem for Linux) instead of PowerShell or CMD, as the Makefile requires Unix-like shell commands.

**Option 1: Git Bash (Recommended for simplicity)**

1. Download and install [Git for Windows](https://git-scm.com/download/win) (includes Git Bash)
2. Install [Chocolatey](https://chocolatey.org/install) (package manager for Windows)
3. Open PowerShell as Administrator and install Go and make:

```powershell
choco install golang make -y
```

4. Close and reopen your terminal, then verify installation in Git Bash:

```bash
go version
make --version
```

**Option 2: WSL2 (Recommended for full Linux experience)**

1. Install WSL2 following [Microsoft's guide](https://learn.microsoft.com/en-us/windows/wsl/install)
2. Install Ubuntu from Microsoft Store
3. Open Ubuntu terminal and run:

```bash
sudo apt update
sudo apt install golang-go make
```

**Verify installation:**

```bash
go version
make --version
```

## Compilation

SynapSeq can be compiled using the provided Makefile.

**For macOS and Linux:**

Simply run:

```bash
make
```

This will automatically compile SynapSeq for your current operating system and architecture, creating a binary in the `bin/` directory.

**For Windows (using Git Bash or WSL2):**

Open Git Bash or your WSL2 terminal and run:

```bash
make build-windows
```

This will generate Windows executables (`.exe`) in the `bin/` directory.

### Installing the Binary

After compilation, you can install the binary system-wide:

**macOS/Linux:**

```bash
sudo make install
```

This will install the SynapSeq binary to `/usr/local/bin/synapseq`.

**Windows (Git Bash or WSL2):**

Using Git Bash (run as Administrator):

```bash
mkdir -p "/c/Program Files/SynapSeq"
cp bin/synapseq-windows-amd64.exe "/c/Program Files/SynapSeq/synapseq.exe"
```

Or using WSL2, you can copy to a Windows directory:

```bash
mkdir -p "/mnt/c/Program Files/SynapSeq"
cp bin/synapseq-windows-amd64.exe "/mnt/c/Program Files/SynapSeq/synapseq.exe"
```

**Adding to PATH:**

After copying the executable, add `C:\Program Files\SynapSeq` to your PATH environment variable.

_Option 1: Using PowerShell (run as Administrator):_

```powershell
# Add to PATH permanently for current user
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\SynapSeq", "User")

# Or add for all users (requires Administrator)
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\SynapSeq", "Machine")
```

_Option 2: Using Windows Settings (GUI):_

1. Open **Start Menu** and search for "Environment Variables"
2. Click **"Edit the system environment variables"**
3. Click **"Environment Variables..."** button
4. Under **"User variables"** or **"System variables"**, find and select **"Path"**
5. Click **"Edit..."**
6. Click **"New"**
7. Add: `C:\Program Files\SynapSeq`
8. Click **"OK"** on all dialogs

After adding to PATH, **restart your terminal** and verify:

```bash
synapseq --version
```

### Installing Documentation (Optional)

**macOS/Linux:**

You can generate and install a man page for offline documentation:

```bash
# Generate the man page (requires pandoc)
make man

# Install the man page system-wide
sudo make install-man
```

After installation, you can access the documentation with:

```bash
man synapseq
```

**Note:** The `man` target requires [pandoc](https://pandoc.org/) to be installed on your system.

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
- `make clean` - Remove all compiled binaries and generated documentation
- `make all` - Same as `make build`
- `make man` - Generate man page documentation (requires pandoc)
- `make install-man` - Install man page system-wide (requires pandoc and sudo)

## Documentation

For detailed information on all features and advanced usage, see the [USAGE.md](docs/USAGE.md) file.

## FAQ

For answers to common questions about SynapSeq and brainwave entrainment, see the [FAQ](docs/FAQ.md).

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

SynapSeq is distributed under the GPL v2 license. See the [COPYING.txt](COPYING.txt) file for details.

### Third-Party Licenses

SynapSeq makes use of third-party libraries, which remain under their own licenses.  
All original code in SynapSeq is licensed under the GNU GPL v2, but the following components are included and redistributed under their respective terms:

- **[beep](https://github.com/gopxl/beep)**  
  License: MIT  
  Used for audio encoding/decoding.

- **[go-yaml](https://github.com/goccy/go-yaml)**  
  License: MIT  
  Used for YAML parsing and processing.

- **[pkg/errors](https://github.com/pkg/errors)**  
  License: BSD 2-Clause  
  Used indirectly via `beep` for error wrapping and stack trace utilities.

All third-party copyright notices and licenses are preserved in this repository in compliance with their original terms.

## Contact

We'd love to hear from you! Here's how to get in touch:

### Issues (Bug Reports & Feature Requests)

Use [GitHub Issues](https://github.com/ruanklein/synapseq/issues) for:

- Bug reports and technical problems
- Feature requests and enhancement suggestions
- Documentation improvements

### Discussions (Questions & Community)

Use [GitHub Discussions](https://github.com/ruanklein/synapseq/discussions) for:

- General questions and support (e.g., "How do I use `@presetlist`?")
- Help with your sequences (e.g., "My sequence isn't working, can you help?")
- Sharing your own sequences and presets with the community
- Discussing ideas and best practices
- Showcasing creative use cases

### Quick Guidelines

- **Found a bug?** → Open an Issue
- **Want a new feature?** → Open an Issue
- **Need help or have questions?** → Start a Discussion
- **Want to share your sequences?** → Post in Discussions
- **General feedback or ideas?** → Start a Discussion
