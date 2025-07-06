#!/bin/bash

# SynapSeq Linux build script
# Builds native binary with MP3, OGG and ALSA support

# Build directory
BUILD_DIR="$PWD/build"

# Source common library
. $BUILD_DIR/lib.sh

# Source directory
SRC_DIR="$PWD/src"

# Get Architecture
ARCH=$(uname -m)

section_header "Building SynapSeq native binary ($ARCH)..."

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

# Define base compilation flags
CFLAGS="-DT_POSIX -Wall -O3 -I."
LIBS="-lm -lpthread"

# Check for MP3 support using pkg-config
if pkg-config --exists mad; then
    info "Including MP3 support using pkg-config (libmad)"
    CFLAGS="$CFLAGS -DMP3_DECODE $(pkg-config --cflags mad)"
    LIBS="$LIBS $(pkg-config --libs mad)"
else
    warning "libmad not found via pkg-config"
    warning "MP3 support will not be included"
    warning "Install libmad using your package manager:"
    warning "  Ubuntu/Debian: sudo apt install libmad0-dev"
    warning "  CentOS/RHEL:   sudo yum install libmad-devel"
    warning "  Fedora:        sudo dnf install libmad-devel"
    warning "  Arch:          sudo pacman -S libmad"
fi

# Check for OGG support using pkg-config
if pkg-config --exists vorbis vorbisfile ogg; then
    info "Including OGG support using pkg-config (libvorbis + libogg)"
    CFLAGS="$CFLAGS -DOGG_DECODE $(pkg-config --cflags vorbis vorbisfile ogg)"
    LIBS="$LIBS $(pkg-config --libs vorbis vorbisfile ogg)"
else
    warning "libvorbis/libogg not found via pkg-config"
    warning "OGG support will not be included"
    warning "Install vorbis libraries using your package manager:"
    warning "  Ubuntu/Debian: sudo apt install libvorbis-dev libogg-dev"
    warning "  CentOS/RHEL:   sudo yum install libvorbis-devel libogg-devel"
    warning "  Fedora:        sudo dnf install libvorbis-devel libogg-devel"
    warning "  Arch:          sudo pacman -S libvorbis libogg"
fi

# Compile synapseq
section_header "Starting synapseq compilation..."
info "Compilation flags: $CFLAGS"
info "Libraries: $LIBS"

# Determine output binary name based on architecture
OUTPUT_BINARY=synapseq-linux
if [ "$ARCH" = "aarch64" ]; then
    OUTPUT_BINARY="$BUILD_DIR/dist/${OUTPUT_BINARY}-arm64"
elif [ "$ARCH" = "x86_64" ]; then
    OUTPUT_BINARY="$BUILD_DIR/dist/${OUTPUT_BINARY}-x86_64"
else
    OUTPUT_BINARY="$BUILD_DIR/dist/${OUTPUT_BINARY}-$ARCH"
    warning "Unknown architecture $ARCH, using generic binary name"
fi

gcc $CFLAGS $SRC_DIR/synapseq.c -o $OUTPUT_BINARY $LIBS

if [ $? -eq 0 ]; then
    success "Compilation successful! Binary created: $(basename $OUTPUT_BINARY)"
    info "Architecture: $ARCH"
    # Strip the binary
    strip $OUTPUT_BINARY
else
    error "Compilation failed!"
fi

section_header "Build process completed!" 