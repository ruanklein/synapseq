# Compilation

You need to install Go (v1.25 or later) and make on your system before compiling SynapSeq.

## Table of Contents

- [Installing Go](#installing-go)
- [Compiling SynapSeq](#compiling-synapseq)
- [Installing the Binary](#installing-the-binary)
- [Installing Documentation (Optional)](#installing-documentation-optional)
- [Additional Make Commands](#additional-make-commands)

## Installing Go

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

**Option 1: Git Bash**

1. Download and install [Git for Windows](https://git-scm.com/download/win) (includes Git Bash)
2. Install [Scoop](https://scoop.sh/) (package manager for Windows).

3. Install Go and make using Scoop:

```powershell
scoop update
scoop install go make
```

4. Close and reopen your terminal, then verify installation in Git Bash:

```bash
go version
make --version
```

**Option 2: WSL2**

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

## Compiling SynapSeq

First, clone the repository:

```bash
git clone https://github.com/ruanklein/synapseq.git
cd synapseq
```

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

## Installing the Binary

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
synapseq -h
```

## Installing Documentation (Optional)

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

## Additional Make Commands

- `make build` - Build for your current platform
- `make build-*` - Build for a specific platform (e.g., `make build-windows`, `make build-macos`, `make build-linux`)
- `make clean` - Remove all compiled binaries and generated documentation
- `make all` - Same as `make build`
- `make man` - Generate man page documentation (requires pandoc)
- `make install-man` - Install man page system-wide (requires pandoc and sudo)
