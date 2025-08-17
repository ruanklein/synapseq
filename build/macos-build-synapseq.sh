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

# Get native macOS version
MACOS_MIN_VERSION=$(sw_vers -productVersion | cut -d '.' -f 1-2)

# Define base compilation flags
CFLAGS="-DT_POSIX -arch $ARCH -mmacosx-version-min=$MACOS_MIN_VERSION -I."
LIBS=""

section_header "Building SynapSeq native binary ($ARCH) for macOS $MACOS_MIN_VERSION..."

if [ ! "$(which pkg-config)" 2> /dev/null ]; then
    error "pkg-config is not installed"
    error "Please install pkg-config using 'brew install pkg-config'"
    exit 1
fi

if [[ -n $HOMEBREW_PREFIX ]]; then
    LIB_PATH=$HOMEBREW_PREFIX/lib
elif command -v port >/dev/null 2>&1; then
    MACPORTS_PREFIX=$(dirname $(dirname $(command -v port)))
    LIB_PATH=$MACPORTS_PREFIX/lib
else
    error "Missing required dependencies"
    info "Follow the steps in README.md under Compilation -> macOS to install either Homebrew or MacPorts and required dependencies"
    exit 1
fi

# Check distribution directory
create_dir_if_not_exists "$BUILD_DIR/dist"

# Check for MP3 support - FORCE STATIC LINKING
if pkg-config --exists mad; then
    info "Including MP3 support (STATIC LINKING)"
    CFLAGS="$CFLAGS -DMP3_DECODE $(pkg-config --cflags mad)"
    
    # Force static linking by using .a files directly
    MAD_LIB="$LIB_PATH/libmad.a"
    if [ -f "$MAD_LIB" ]; then
        LIBS="$LIBS $MAD_LIB -lm"
        info "Using static library: $MAD_LIB"
    else
        LIBS="$LIBS $(pkg-config --libs mad)"
        warning "Static libmad.a not found, using dynamic${MACPORTS_PREFIX:+ due to MacPorts not providing a static variant of libmad}"
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
    VORBISFILE_LIB="$LIB_PATH/libvorbisfile.a"
    VORBIS_LIB="$LIB_PATH/libvorbis.a"
    OGG_LIB="$LIB_PATH/libogg.a"
    
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

gcc $CFLAGS $SRC_DIR/synapseq.c -o $BUILD_DIR/dist/synapseq-macos-$ARCH $LIBS

if [ $? -eq 0 ]; then
    success "Compilation successful! Binary created: dist/synapseq-macos-$ARCH"
    # Strip the binary
    strip $BUILD_DIR/dist/synapseq-macos-$ARCH
else
    error "Compilation failed!"
fi

section_header "Build process completed!" 