<h1 align="center">SynapSeq</h1>

<p align="center">
  <a href="#installation">Installation</a> |
  <a href="https://synapseq-hub.ruan.sh/">Examples</a> |
  <a href="docs/USAGE.md">Documentation</a> |
  <a href="docs/FAQ.md">FAQ</a>
</p>

<p align="center">
  <p align="center">
  <a href="https://github.com/ruanklein/synapseq/releases/latest"><img src="https://img.shields.io/github/v/release/ruanklein/synapseq?color=blue&logo=github" alt="Release"></a>
  <a href="COPYING.txt"><img src="https://img.shields.io/badge/license-GPL%20v2-blue.svg?logo=open-source-initiative&logoColor=white" alt="License"></a>
  <a href="#winget"><img src="https://img.shields.io/badge/Winget-Install-blue?logo=gnometerminal&logoColor=cyan" alt="Winget"></a>
  <a href="#homebrew"><img src="https://img.shields.io/badge/Homebrew-Install-brightgreen?logo=homebrew&logoColor=white" alt="Homebrew"></a>
  <a href="#scoop"><img src="https://img.shields.io/badge/Scoop-Install-blue?logo=gnometerminal&logoColor=00AEEF" alt="Scoop"></a>
  <a href="https://github.com/ruanklein/synapseq/commits"><img src="https://img.shields.io/github/commit-activity/m/ruanklein/synapseq?color=ff69b4&logo=git" alt="Commit Activity"></a>
</p>
</p>

<p align="center"><strong>Synapse-Sequenced Brainwave Generator</strong></p>

SynapSeq is a lightweight engine that sequences audio tones to guide brainwave states like relaxation, focus, and meditation using a simple text-based format.

## Installation

SynapSeq can be installed via package managers or by downloading precompiled binaries.

### Windows

#### Winget

If you use Winget, you can install SynapSeq directly from the official Microsoft repository:

```powershell
winget install synapseq
```

Check installation:

```powershell
synapseq -h
```

#### Scoop

If you use Scoop, you can install SynapSeq directly from the official bucket:

```powershell
scoop bucket add synapseq https://github.com/ruanklein/scoop-synapseq
scoop install synapseq
```

Check installation:

```powershell
synapseq -h
```

#### Windows Native Integration

On Windows, you can enable native file association and context menu integration for a better user experience:

```powershell
synapseq -install-file-association
```

To remove the integration:

```powershell
synapseq -uninstall-file-association
```

For more details, see the [-install-file-association](./docs/USAGE.md#-install-file-association-windows-only) command documentation.

### macOS/Linux

#### Homebrew

Just install via Homebrew:

```bash
brew tap ruanklein/synapseq
brew install synapseq
```

Check installation:

```bash
synapseq -h
```

### Precompiled Binaries

You can download the latest precompiled binaries for Windows, macOS, or Linux from the [Releases](https://github.com/ruanklein/synapseq/releases/latest) page.

## Compilation

For users who prefer to compile from source, please follow the [Compilation](docs/COMPILE.md) instructions.

## Documentation

For detailed information on all features and advanced usage, see the [USAGE.md](docs/USAGE.md) file.

## Go Library

SynapSeq can also be used as a Go library in your own projects. See the [Go Library Documentation](https://pkg.go.dev/github.com/ruanklein/synapseq/v3/core) for instructions on how to integrate SynapSeq into your Go applications.

## FAQ

For answers to common questions about SynapSeq and brainwave entrainment, see the [FAQ](docs/FAQ.md).

## Contributing

We welcome contributions!

Please read the [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines on how to contribute code, bug fixes, and documentation to the project.

### Contributing Sequence Files

If you'd like to share your own `.spsq` sequence files with the community, please contribute them to the [SynapSeq Hub Repository](https://github.com/ruanklein/synapseq-hub). All sequence files contributed to the Hub are licensed under [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).

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

- **[google/uuid](https://github.com/google/uuid)**  
  License: BSD 3-Clause  
  Copyright © 2009-2014 Google Inc.  
  Used for UUID generation and unique identifier handling.

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
