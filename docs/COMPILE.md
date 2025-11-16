# Compilation

You need to install Go (v1.25 or later) and make on your system before compiling SynapSeq.

## Table of Contents

- [Installing Go](#installing-go)
- [Compiling SynapSeq](#compiling-synapseq)
- [Installing the Binary](#installing-the-binary)
- [Compiling without Hub Support](#compiling-without-hub-support)
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

**For POSIX systems (macOS/Linux):**

Simply run:

```bash
make
```

This will automatically compile SynapSeq for your current operating system and architecture, creating a binary in the `bin/` directory.

**For Windows:**

Windows users should use the platform-specific targets to ensure proper `.exe` extension, application icon, and Windows-specific command-line options:

```bash
make build-windows-amd64    # Windows 64-bit (Intel/AMD) - Recommended
make build-windows-arm64    # Windows 64-bit (ARM)
```

**Note:** Do not use `make build` on Windows, as it will create a binary without the `.exe` extension, missing the application icon and Windows-specific features.

**For cross-compilation to other platforms:**

You can build for different platforms and architectures:

```bash
# Linux
make build-linux-amd64      # Linux 64-bit (Intel/AMD)
make build-linux-arm64      # Linux 64-bit (ARM)

# macOS
make build-macos            # macOS ARM64 (Apple Silicon)
```

**Note for Windows builds:** The Makefile automatically generates Windows resource files (including application icon and metadata) when building for Windows. This requires the `goversioninfo` tool, which will be automatically downloaded during the build process.

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

## Compiling without Hub Support

If you prefer to use SynapSeq without any Hub features (including analytics tracking), you can compile a Hub-disabled version using the `nohub` build tag. This completely removes all Hub-related code from the binary, resulting in a smaller executable with no network connections to the SynapSeq Hub infrastructure.

**To compile without Hub support:**

```bash
# Build for POSIX system (macOS/Linux)
make build-nohub

# Or build for specific platforms and architectures
make build-windows-nohub-amd64    # Windows 64-bit (Intel/AMD)
make build-windows-nohub-arm64    # Windows 64-bit (ARM)
make build-linux-nohub-amd64      # Linux 64-bit (Intel/AMD)
make build-linux-nohub-arm64      # Linux 64-bit (ARM)
make build-macos-nohub            # macOS ARM64 (Apple Silicon)
```

**What's different in the Hub-disabled build:**

- All `-hub-*` commands will return an error message
- No network connections to the Hub infrastructure
- No analytics or tracking of any kind
- Smaller binary size (Hub code is excluded)
- All other SynapSeq features work normally

**Installing the Hub-disabled binary:**

After compilation, you can install it system-wide:

```bash
# macOS/Linux
sudo make install-nohub
# Windows (in Git Bash as Administrator)
mkdir -p "/c/Program Files/SynapSeq"
cp bin/synapseq-windows-amd64-nohub.exe "/c/Program Files/SynapSeq/synapseq-nohub.exe"
```

## Additional Make Commands

Other available make targets:

**Testing:**

- `make test` - Run all tests

**Cleanup:**

- `make clean` - Remove all compiled binaries and generated files

**Utilities:**

- `make all` - Same as `make build`
- `make prepare` - Create `bin/` directory
- `make windows-res-amd64` - Generate Windows resource file (icon and metadata) for x64 build
- `make windows-res-arm64` - Generate Windows resource file (icon and metadata) for ARM64 build
