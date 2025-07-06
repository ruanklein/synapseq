#!/bin/bash

# SynapSeq Linux static build script
# Builds native binary with MP3 and OGG (static linking)

# Build directory
BUILD_DIR="$PWD/build"

# Source common library
. $BUILD_DIR/lib.sh

# Source directory
SRC_DIR="$PWD/src"

# Get Architecture
ARCH=$(uname -m)

section_header "Building SynapSeq native binary (static linking - $ARCH)..."

if [ ! "$(which pkg-config)" 2> /dev/null ]; then
    error "pkg-config is not installed"
    error "Please install pkg-config using your package manager:"
    error "  Ubuntu/Debian: sudo apt install pkg-config"
    error "  CentOS/RHEL:   sudo yum install pkgconfig"
    error "  Fedora:        sudo dnf install pkgconfig"
    error "  Arch:          sudo pacman -S pkgconf"
    exit 1
fi

# Check distribution directory
create_dir_if_not_exists "$BUILD_DIR/dist"

# Define base compilation flags for static linking
CFLAGS="-DT_POSIX -Wall -O3 -I. -static"
LIBS="-lm -lpthread"

# Check for MP3 support using pkg-config (static)
if pkg-config --exists mad; then
    info "Including MP3 support using pkg-config (libmad - static)"
    CFLAGS="$CFLAGS -DMP3_DECODE $(pkg-config --cflags mad)"
    LIBS="$LIBS $(pkg-config --libs --static mad)"
else
    warning "libmad not found via pkg-config"
    warning "MP3 support will not be included"
    warning "Install libmad development packages using your package manager:"
    warning "  Ubuntu/Debian: sudo apt install libmad0-dev"
    warning "  CentOS/RHEL:   sudo yum install libmad-devel"
    warning "  Fedora:        sudo dnf install libmad-devel"
    warning "  Arch:          sudo pacman -S libmad"
fi

# Check for OGG support using pkg-config (static)
if pkg-config --exists vorbis vorbisfile ogg; then
    info "Including OGG support using pkg-config (libvorbis + libogg - static)"
    CFLAGS="$CFLAGS -DOGG_DECODE $(pkg-config --cflags vorbis vorbisfile ogg)"
    LIBS="$LIBS $(pkg-config --libs --static vorbis vorbisfile ogg)"
else
    warning "libvorbis/libogg not found via pkg-config"
    warning "OGG support will not be included"
    warning "Install vorbis development packages using your package manager:"
    warning "  Ubuntu/Debian: sudo apt install libvorbis-dev libogg-dev"
    warning "  CentOS/RHEL:   sudo yum install libvorbis-devel libogg-devel"
    warning "  Fedora:        sudo dnf install libvorbis-devel libogg-devel"
    warning "  Arch:          sudo pacman -S libvorbis libogg"
fi

# Add additional libraries that might be needed for static linking
LIBS="$LIBS -ldl"

# Compile synapseq (static)
section_header "Starting synapseq static compilation..."
info "Compilation flags: $CFLAGS"
info "Libraries: $LIBS"

# Determine output binary name based on architecture
OUTPUT_BINARY=synapseq-linux-static
if [ "$ARCH" = "aarch64" ]; then
    OUTPUT_BINARY="$BUILD_DIR/dist/${OUTPUT_BINARY}-arm64"
elif [ "$ARCH" = "x86_64" ]; then
    OUTPUT_BINARY="$BUILD_DIR/dist/${OUTPUT_BINARY}-x86_64"
else
    OUTPUT_BINARY="$BUILD_DIR/dist/${OUTPUT_BINARY}-$ARCH"
    warning "Unknown architecture $ARCH, using generic binary name"
fi

info "Output binary: $(basename $OUTPUT_BINARY)"

# Attempt static compilation
gcc $CFLAGS $SRC_DIR/synapseq.c -o $OUTPUT_BINARY $LIBS

if [ $? -eq 0 ]; then
    success "Static compilation successful! Binary created: $(basename $OUTPUT_BINARY)"
    info "Architecture: $ARCH"
    
    # Check if binary is actually statically linked
    if ldd $OUTPUT_BINARY 2>&1 | grep -q "not a dynamic executable"; then
        success "Binary is statically linked"
    else
        warning "Binary may not be fully statically linked:"
        ldd $OUTPUT_BINARY
    fi
    
    # Strip the binary
    strip $OUTPUT_BINARY
else
    error "Static compilation failed!"
    warning ""
    warning "Static compilation troubleshooting:"
    warning "1. Make sure static libraries are installed:"
    warning "   Ubuntu/Debian: sudo apt install libc6-dev libmad0-dev libvorbis-dev libogg-dev"
    warning "2. Some distributions may not provide static libraries by default"
    warning "3. Try compiling with dynamic linking using linux-build-synapseq.sh instead"
fi

section_header "Static build process completed!" 