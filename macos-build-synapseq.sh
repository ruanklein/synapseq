#!/bin/bash

# SynapSeq macOS build script
# Builds a native binary with MP3 and OGG support

# Source common library
. ./lib.sh

# Get Architecture
ARCH=$(uname -m)

section_header "Building SynapSeq native binary ($ARCH) with MP3 and OGG support..."

if [ ! "$(which pkg-config)" 2> /dev/null ]; then
    error "pkg-config is not installed"
    error "Please install pkg-config using 'brew install pkg-config'"
    exit 1
fi

# Check distribution directory
create_dir_if_not_exists "dist"

# Define base compilation flags
CFLAGS="-DT_MACOSX -arch $ARCH -mmacosx-version-min=12.0 -I."
LIBS="-framework CoreAudio"

# Get the version number from the VERSION file
VERSION=$(cat VERSION)

# Check for MP3 support using pkg-config
if pkg-config --exists mad; then
    info "Including MP3 support using pkg-config (libmad)"
    CFLAGS="$CFLAGS -DMP3_DECODE $(pkg-config --cflags mad)"
    LIBS="$LIBS $(pkg-config --libs mad)"
else
    warning "libmad not found via pkg-config"
    warning "MP3 support will not be included"
    warning "Install libmad using 'brew install mad'"
fi

# Check for OGG support using pkg-config
if pkg-config --exists vorbis vorbisfile ogg; then
    info "Including OGG support using pkg-config (libvorbis + libogg)"
    CFLAGS="$CFLAGS -DOGG_DECODE $(pkg-config --cflags vorbis vorbisfile ogg)"
    LIBS="$LIBS $(pkg-config --libs vorbis vorbisfile ogg)"
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
sed "s/__VERSION__/\"$VERSION\"/" synapseq.c > synapseq.tmp.c

gcc $CFLAGS synapseq.tmp.c -o dist/synapseq-macos-universal $LIBS

if [ $? -eq 0 ]; then
    success "Compilation successful! Universal binary created: dist/synapseq-macos-universal"
    #info "Supported architectures:"
    #lipo -info dist/synapseq-macos-universal
else
    error "Compilation failed!"
fi

# Remove the temporary file
rm -f synapseq.tmp.c

section_header "Build process completed!" 