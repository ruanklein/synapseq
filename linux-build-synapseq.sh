#!/bin/bash

# SynapSeq Linux build script
# Builds 32-bit and 64-bit binaries with MP3, OGG and ALSA support on x86_64
# Builds only ARM64 binary on aarch64 platforms

# Source common library
. ./lib.sh

section_header "Building SynapSeq for Linux with MP3, OGG and ALSA support..."

# Check for required tools
check_required_tools gcc

# Check distribution directory
create_dir_if_not_exists "dist"

# Detect host architecture
HOST_ARCH=$(uname -m)
info "Detected host architecture: $HOST_ARCH"

# Define paths for libraries
LIB_PATH_32="libs/linux-x86-libmad.a"
LIB_PATH_64="libs/linux-x86_64-libmad.a"
LIB_PATH_ARM64="libs/linux-arm64-libmad.a"
OGG_LIB_PATH_32="libs/linux-x86-libogg.a"
OGG_LIB_PATH_64="libs/linux-x86_64-libogg.a"
OGG_LIB_PATH_ARM64="libs/linux-arm64-libogg.a"
TREMOR_LIB_PATH_32="libs/linux-x86-libvorbisidec.a"
TREMOR_LIB_PATH_64="libs/linux-x86_64-libvorbisidec.a"
TREMOR_LIB_PATH_ARM64="libs/linux-arm64-libvorbisidec.a"

# Get the version number from the VERSION file
VERSION=$(cat VERSION)

# Skip 32-bit build on ARM64
SKIP_32BIT=0
if [ "$HOST_ARCH" = "aarch64" ]; then
    SKIP_32BIT=1
    warning "32-bit compilation is not supported on ARM64, skipping..."
fi

# Build 32-bit version
if [ $SKIP_32BIT = 0 ]; then
    section_header "Building 32-bit version..."

    # Set up compilation flags for 32-bit
    CFLAGS_32="-DT_LINUX -m32 -Wall -O3 -I."
    LIBS_32="-lm -lpthread -lasound"

    # Check for MP3 support (32-bit)
    if [ -f "$LIB_PATH_32" ]; then
        info "Including MP3 support for 32-bit using: $LIB_PATH_32"
        CFLAGS_32="$CFLAGS_32 -DMP3_DECODE"
        LIBS_32="$LIBS_32 $LIB_PATH_32"
    else
        warning "MP3 library not found at $LIB_PATH_32"
        warning "MP3 support will not be included in 32-bit build"
        warning "Run ./linux-build-libs.sh to build the required libraries"
    fi

    # Check for OGG support (32-bit)
    if [ -f "$OGG_LIB_PATH_32" ] && [ -f "$TREMOR_LIB_PATH_32" ]; then
        info "Including OGG support for 32-bit using: $OGG_LIB_PATH_32 and $TREMOR_LIB_PATH_32"
        CFLAGS_32="$CFLAGS_32 -DOGG_DECODE"
        # Order is important: first tremor, then ogg
        LIBS_32="$LIBS_32 $TREMOR_LIB_PATH_32 $OGG_LIB_PATH_32"
    else
        warning "OGG libraries not found at $OGG_LIB_PATH_32 or $TREMOR_LIB_PATH_32"
        warning "OGG support will not be included in 32-bit build"
        warning "Run ./linux-build-libs.sh to build the required libraries"
    fi

    # Compile 32-bit version
    info "Compiling 32-bit version with flags: $CFLAGS_32"
    info "Libraries: $LIBS_32"

    # Try to compile with 32-bit support
    gcc $CFLAGS_32 synapseq.c -o dist/synapseq-linux32 $LIBS_32

    if [ $? -eq 0 ]; then
        success "32-bit compilation successful! Binary created: synapseq-linux32"
    else
        error "32-bit compilation failed! You may need to install 32-bit development libraries."
    fi
else
    warning "Skipping 32-bit build..."
fi

# Build 64-bit version
section_header "Building 64-bit version..."

# Set up compilation flags for 64-bit
if [ "$HOST_ARCH" = "aarch64" ]; then
    # On ARM64, don't use -m64 flag as it's not supported
    CFLAGS_64="-DT_LINUX -Wall -O3 -I."
    info "Running on ARM64, using native gcc for 64-bit compilation"
else
    CFLAGS_64="-DT_LINUX -m64 -Wall -O3 -I."
fi
LIBS_64="-lm -lpthread -lasound"

# Check for MP3 support for 64-bit or ARM64
if [ "$HOST_ARCH" = "aarch64" ]; then
    if [ -f "$LIB_PATH_ARM64" ]; then
        info "Including MP3 support for ARM64 using: $LIB_PATH_ARM64"
        CFLAGS_64="$CFLAGS_64 -DMP3_DECODE"
        LIBS_64="$LIBS_64 $LIB_PATH_ARM64"
    else
        warning "MP3 library not found at $LIB_PATH_ARM64"
        warning "MP3 support will not be included in ARM64 build"
        warning "Run ./linux-build-libs.sh to build the required libraries"
    fi
else
    if [ -f "$LIB_PATH_64" ]; then
        info "Including MP3 support for 64-bit using: $LIB_PATH_64"
        CFLAGS_64="$CFLAGS_64 -DMP3_DECODE"
        LIBS_64="$LIBS_64 $LIB_PATH_64"
    else
        warning "MP3 library not found at $LIB_PATH_64"
        warning "MP3 support will not be included in 64-bit build"
        warning "Run ./linux-build-libs.sh to build the required libraries"
    fi
fi

# Check for OGG support for 64-bit or ARM64
if [ "$HOST_ARCH" = "aarch64" ]; then
    if [ -f "$OGG_LIB_PATH_ARM64" ] && [ -f "$TREMOR_LIB_PATH_ARM64" ]; then
        info "Including OGG support for ARM64 using: $OGG_LIB_PATH_ARM64 and $TREMOR_LIB_PATH_ARM64"
        CFLAGS_64="$CFLAGS_64 -DOGG_DECODE"
        # Order is important: first tremor, then ogg
        LIBS_64="$LIBS_64 $TREMOR_LIB_PATH_ARM64 $OGG_LIB_PATH_ARM64"
    else
        warning "OGG libraries not found at $OGG_LIB_PATH_ARM64 or $TREMOR_LIB_PATH_ARM64"
        warning "OGG support will not be included in ARM64 build"
        warning "Run ./linux-build-libs.sh to build the required libraries"
    fi
else
    if [ -f "$OGG_LIB_PATH_64" ] && [ -f "$TREMOR_LIB_PATH_64" ]; then
        info "Including OGG support for 64-bit using: $OGG_LIB_PATH_64 and $TREMOR_LIB_PATH_64"
        CFLAGS_64="$CFLAGS_64 -DOGG_DECODE"
        # Order is important: first tremor, then ogg
        LIBS_64="$LIBS_64 $TREMOR_LIB_PATH_64 $OGG_LIB_PATH_64"
    else
        warning "OGG libraries not found at $OGG_LIB_PATH_64 or $TREMOR_LIB_PATH_64"
        warning "OGG support will not be included in 64-bit build"
        warning "Run ./linux-build-libs.sh to build the required libraries"
    fi
fi

# Compile 64-bit version
info "Compiling 64-bit version with flags: $CFLAGS_64"
info "Libraries: $LIBS_64"

# Replace VERSION with the actual version number
sed "s/__VERSION__/\"$VERSION\"/" synapseq.c > synapseq.tmp.c

if [ "$HOST_ARCH" = "aarch64" ]; then
    gcc $CFLAGS_64 synapseq.tmp.c -o dist/synapseq-linux-arm64 $LIBS_64
else
    gcc $CFLAGS_64 synapseq.tmp.c -o dist/synapseq-linux64 $LIBS_64
fi

if [ $? -eq 0 ]; then
    if [ "$HOST_ARCH" = "aarch64" ]; then
        success "64-bit compilation successful! Created ARM64 binary: synapseq-linux-arm64"
    else
        success "64-bit compilation successful! Created 64-bit binary: synapseq-linux64"
    fi
else
    error "64-bit compilation failed!"
fi

# Remove the temporary file
rm -f synapseq.tmp.c

section_header "Build process completed!" 