#!/bin/bash

# SynapSeq macOS build script
# Builds a native binary with MP3 and OGG support

# Build directory
BUILD_DIR="$PWD/build"

# Source common library
. $BUILD_DIR/lib.sh

# Source directory
SRC_DIR="$PWD/src"

# Get Architecture
ARCH=$(uname -m)

# Default macOS version
MACOS_MIN_VERSION=15.0

# Define base compilation flags
CFLAGS="-DT_MACOSX -arch $ARCH -mmacosx-version-min=$MACOS_MIN_VERSION -I."
LIBS="-framework CoreAudio"

# Get the version number from the VERSION file
VERSION=$(cat $BUILD_DIR/VERSION)

section_header "Building SynapSeq native binary ($ARCH) for macOS $MACOS_MIN_VERSION..."

if [ ! "$(which pkg-config)" 2> /dev/null ]; then
    error "pkg-config is not installed"
    error "Please install pkg-config using 'brew install pkg-config'"
    exit 1
fi

if [ -z $HOMEBREW_PREFIX ]; then
    error "HOMEBREW_PREFIX is not set"
    error "Please install Homebrew using 'https://brew.sh/'"
    exit 1
fi

# Check distribution directory
create_dir_if_not_exists "$BUILD_DIR/dist"

# Check for MP3 support - FORCE STATIC LINKING
if pkg-config --exists mad; then
    info "Including MP3 support (STATIC LINKING)"
    CFLAGS="$CFLAGS -DMP3_DECODE $(pkg-config --cflags mad)"
    
    # Force static linking by using .a files directly
    MAD_LIB="$HOMEBREW_PREFIX/lib/libmad.a"
    if [ -f "$MAD_LIB" ]; then
        LIBS="$LIBS $MAD_LIB -lm"
        info "Using static library: $MAD_LIB"
    else
        LIBS="$LIBS $(pkg-config --libs mad)"
        warning "Static libmad.a not found, using dynamic"
    fi
else
    warning "libmad not found via pkg-config"
    warning "MP3 support will not be included"
    warning "Install libmad using 'brew install mad'"
fi

# Check for OGG support - FORCE STATIC LINKING
if pkg-config --exists vorbis vorbisfile ogg; then
    info "Including OGG support (STATIC LINKING)"
    CFLAGS="$CFLAGS -DOGG_DECODE $(pkg-config --cflags vorbis vorbisfile ogg)"
    
    # Force static linking by using .a files directly
    VORBISFILE_LIB="$HOMEBREW_PREFIX/lib/libvorbisfile.a"
    VORBIS_LIB="$HOMEBREW_PREFIX/lib/libvorbis.a"
    OGG_LIB="$HOMEBREW_PREFIX/lib/libogg.a"
    
    if [ -f "$VORBISFILE_LIB" ] && [ -f "$VORBIS_LIB" ] && [ -f "$OGG_LIB" ]; then
        LIBS="$LIBS $VORBISFILE_LIB $VORBIS_LIB $OGG_LIB -lm"
        info "Using static libraries: libvorbisfile.a, libvorbis.a, libogg.a"
    else
        LIBS="$LIBS $(pkg-config --libs vorbis vorbisfile ogg)"
        warning "Static vorbis/ogg libraries not found, using dynamic"
    fi
else
    warning "libvorbis/libogg not found via pkg-config"
    warning "OGG support will not be included"
    warning "Install vorbis libraries using 'brew install libogg libvorbis'"
fi

# Compile synapseq
section_header "Starting synapseq compilation..."
info "Compilation flags: $CFLAGS"
info "Libraries: $LIBS"

# Replace VERSION with the actual version number
sed "s/__VERSION__/\"$VERSION\"/" $SRC_DIR/synapseq.c > $SRC_DIR/synapseq.tmp.c

gcc $CFLAGS $SRC_DIR/synapseq.tmp.c -o $BUILD_DIR/dist/synapseq-macos-universal $LIBS

if [ $? -eq 0 ]; then
    success "Compilation successful! Universal binary created: dist/synapseq-macos-universal"
    # Strip the binary
    strip $BUILD_DIR/dist/synapseq-macos-universal
else
    error "Compilation failed!"
fi

# Remove the temporary file
rm -f $SRC_DIR/synapseq.tmp.c



section_header "Build process completed!" 