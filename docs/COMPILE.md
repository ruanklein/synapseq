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

For Windows users, we recommend using **Git Bash** instead of PowerShell or CMD, as the Makefile requires Unix-like shell commands.

1. Install [Git for Windows](https://git-scm.com/download/win) (includes Git Bash).
   After installation, you’ll have both **Git Bash** and **PowerShell** available.

2. Install [Scoop](https://scoop.sh/).
   Open PowerShell and run:

   ```powershell
   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
   Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
   ```

3. Install **Go and Make** using Scoop.
   In PowerShell, run:

   ```powershell
   scoop update
   scoop install go make
   ```

4. Open **Git Bash** and verify that everything is available:

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

**For Windows**

In Git Bash:

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

**Windows**

Using Git Bash (run as Administrator):

```bash
mkdir -p "/c/Program Files/SynapSeq"
cp bin/synapseq-windows-amd64.exe "/c/Program Files/SynapSeq/synapseq.exe"
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

After adding to PATH, **restart Git Bash or PowerShell** and verify:

```bash
synapseq -h
```

## Installing Documentation (macOS/Linux)

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
