# <img style="border-radius: 15%;" src="build/assets/synapseq.png" alt="SynapSeq Logo" width="32" height="32"> SynapSeq - Synapse-Sequenced Brainwave Generator

SynapSeq is a text-based tool for generating tones to stimulate the brainwave to help with meditation, relaxation, and altering states of consciousness.

## 📑 Table of Contents

- [💡 Example](#-example)
- [📥 Installation](#-installation)
  - [🪟 Installing on Windows](#🪟Windows-Installation)
  - [🍎 Installing on macOS](#-installing-on-macos)
- [🚀 Basic Usage](#-basic-usage)
- [📚 Documentation](#-documentation)
- [🔍 Research](#-research)
- [🛠️ Compilation](#️-compilation)
  - [📁 Build Scripts Structure](#-build-scripts-structure)
  - [🐳 Building with Docker](#-option-1-using-docker-compose-simplest-method)
  - [💻 Building Natively](#-option-2-building-natively)
- [⚖️ License](#️-license)
- [👏 Credits](#-credits)

## 💡 Example

To get started with SynapSeq, create a new text file called `Relax.spsq` with the following content and double click on the file to open it with SynapSeq (Windows/macOS) or run `synapseq Relax.spsq` on Terminal (all platforms).

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

## 📥 Installation

You can download pre-built binaries on Windows (only 64-bit) and macOS (only Apple Silicon for macOS 15+) from the [releases page](https://github.com/ruanklein/synapseq/releases).

### 🪟 Windows Installation

1. Download the installer:

   - [synapseq-windows-setup.exe](https://github.com/ruanklein/synapseq/releases/download/v2.0.0/synapseq-windows-setup.exe)

2. Verify the SHA256 checksum of the installer with the checksum on the releases page.

3. Run the installer and follow the instructions.

⚠️ **Warning about antivirus on Windows**

Some versions of Windows Defender or other antivirus software may falsely detect `SynapSeq` as a threat.

This happens because the executable is **not digitally signed**, and as a command-line program, it may be flagged as suspicious by default.

`SynapSeq` is an open-source project, and the source code is publicly available in this repository for inspection.

✅ **Temporary solution:** if you trust the source of the executable, add an exception in your antivirus for the file or the folder where `SynapSeq` is installed.

### 🍎 Installing on macOS

1. Download the macOS Installer:

   - [SynapSeq Installer.dmg](https://github.com/ruanklein/synapseq/releases/download/v2.0.0/SynapSeq-Installer.dmg)

2. Verify the SHA256 checksum of the installer with the checksum on the releases page.

3. Open the DMG file and drag the `SynapSeq` application to the Applications folder.

4. Run the `SynapSeq` application from the Applications folder, accept the license agreement and click the `View Examples` button to view examples of spsq files.

5. Click in the .spsq file to play, edit or convert it. Also, you can drop spsq files on the `SynapSeq` application icon to open them.

**Important:** The `SynapSeq` application is not digitally signed, so you may need to add an exception on the `System Settings -> Security & Privacy -> General tab`.

### 🐧 Linux Installation

In Linux, you need build SynapSeq from source. See [Compilation](#-compilation) section for more details.

## 🛠️ Compilation

SynapSeq can be compiled for macOS, Linux and Windows.

### 🍎 For macOS

Install the "Xcode Command Line Tools" in your system.

```bash
xcode-select --install
```

Install [homebrew](https://brew.sh/) if you don't have it yet.

Install dependencies:

```bash
brew install pkg-config libvorbis libmad create-dmg pandoc
```

Run the build script to create the binary:

```bash
./build/macos-build-synapseq.sh
```

The binary will be created in the `build/dist` folder.

If you want to create a installer DMG file, run the following script to create the installer DMG file:

```bash
./build/macos-create-installer.sh
```

The installer DMG file will be created in the `build/dist` folder.

### 🐧 For GNU/Linux

In Ubuntu/Debian based distributions, install dependencies:

```bash
sudo apt-get install build-essential pkg-config libasound2-dev libvorbis-dev libogg-dev libmad0-dev
```

Run the build script to create the binary:

```bash
./build/linux-build-synapseq.sh
```

The binary will be created in the `build/dist` folder.

### 🪟 For Windows

In Windows, the best way to build SynapSeq is using [Docker](https://www.docker.com/) with WSL2.

1. Install [WSL2](https://learn.microsoft.com/en-us/windows/wsl/install)
2. Install [Docker Desktop](https://www.docker.com/products/docker-desktop/)

Run this sequence of commands to build SynapSeq (only x86-64):

```bash
docker compose -f build/compose.yml up build-windows-libs-x86-64
docker compose -f build/compose.yml up build-windows-synapseq-x86-64
docker compose -f build/compose.yml up build-windows-installer-x86-64 # Optional, if you want to create a installer
```

The `.exe` will be created in the `build/dist` folder.

## 🚀 Basic Usage

See [USAGE.md](USAGE.md) for more information on how to use SynapSeq.

## 📚 Documentation

For detailed information on all features, see the [SYNAPSEQ.txt](docs/SYNAPSEQ.txt) file.

## 🔍 Research

For the scientific background behind SynapSeq, check out [RESEARCH.md](RESEARCH.md).

## ⚖️ License

SynapSeq is distributed under the GPL license. See the [COPYING.txt](COPYING.txt) file for details.

## 👏 Credits

SynapSeq is based on the SBaGen by Jim Peters. See [SBaGen project](https://uazu.net/sbagen/).

ALSA support is based on this [patch](https://github.com/jave/sbagen-alsa/blob/master/sbagen.c).
