#!/bin/bash

# SynapSeq Windows build script
# Builds 64-bit Windows binaries with MP3 and OGG support using MinGW

# Build directory
BUILD_DIR="$PWD/build"

# Source common library
. $BUILD_DIR/lib.sh

# Source directory
SRC_DIR="$PWD/src"

# Binary name
OUTPUT_BINARY=synapseq-windows-win64.exe

section_header "Building SynapSeq for Windows (64-bit)..."

# Check for MinGW cross-compilers
if ! command -v x86_64-w64-mingw32-gcc &> /dev/null; then
    error "MinGW cross-compilers not found. Please install them."
    info "On Debian/Ubuntu: sudo apt-get install mingw-w64"
    info "On Fedora: sudo dnf install mingw64-gcc"
    info "On Arch: sudo pacman -S mingw-w64-gcc"
    exit 1
fi

# Check distribution directory
create_dir_if_not_exists "$BUILD_DIR/dist"

# Define paths for libraries. Change it to the correct path for your system.
LIBMAD_PATH_64="$BUILD_DIR/libs/libmad-win64.a"
LIBOGG_PATH_64="$BUILD_DIR/libs/libogg-win64.a"
LIBVORBIS_PATH_64="$BUILD_DIR/libs/libvorbis-win64.a"
LIBVORBISFILE_PATH_64="$BUILD_DIR/libs/libvorbisfile-win64.a"

# Build 64-bit version
section_header "Building 64-bit version..."

# Set up compilation flags for 64-bit
CFLAGS_64="-DT_WIN32 -Wall -O3"
LIBS_64="-lwinmm"

# Check for MP3 support (64-bit)
if [ -f "$LIBMAD_PATH_64" ]; then
    info "Including MP3 support for 64-bit using: $LIBMAD_PATH_64"
    CFLAGS_64="$CFLAGS_64 -DMP3_DECODE"
    LIBS_64="$LIBS_64 $LIBMAD_PATH_64"
else
    warning "MP3 library not found at $LIBMAD_PATH_64"
    warning "MP3 support will not be included in 64-bit build"
    warning "Run ./windows-build-libs.sh to build the required libraries"
fi

# Check for OGG support (64-bit)
if [ -f "$LIBOGG_PATH_64" ] && [ -f "$LIBVORBIS_PATH_64" ] && [ -f "$LIBVORBISFILE_PATH_64" ]; then
    info "Including OGG support for 64-bit using: $LIBOGG_PATH_64 and $LIBVORBIS_PATH_64 and $LIBVORBISFILE_PATH_64"
    CFLAGS_64="$CFLAGS_64 -DOGG_DECODE -I$BUILD_DIR/libs/include"
    # Order is important: first tremor, then ogg
    LIBS_64="$LIBS_64 $LIBVORBISFILE_PATH_64 $LIBVORBIS_PATH_64 $LIBOGG_PATH_64"
else
    warning "OGG libraries not found at $LIBOGG_PATH_64 or $LIBVORBIS_PATH_64 or $LIBVORBISFILE_PATH_64"
    warning "OGG support will not be included in 64-bit build"
    warning "Run ./windows-build-libs.sh to build the required libraries"
fi

# Compile 64-bit version
info "Compiling 64-bit version with flags: $CFLAGS_64"
info "Libraries: $LIBS_64"

x86_64-w64-mingw32-gcc $CFLAGS_64 $SRC_DIR/synapseq.c -o $BUILD_DIR/dist/$OUTPUT_BINARY $LIBS_64

if [ $? -eq 0 ]; then
    success "64-bit compilation successful! Created 64-bit binary: $BUILD_DIR/dist/$OUTPUT_BINARY"
else
    error "64-bit compilation failed!"
fi

section_header "Build process completed!" 