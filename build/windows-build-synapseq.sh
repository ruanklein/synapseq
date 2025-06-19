#!/bin/bash

# SynapSeq Windows build script
# Builds 64-bit Windows binaries with MP3 and OGG support using MinGW

# Build directory
BUILD_DIR="$PWD/build"

# Source common library
. $BUILD_DIR/lib.sh

# Source directory
SRC_DIR="$PWD/src"

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

# Get version from VERSION file
VERSION=$(cat $BUILD_DIR/VERSION)

# Extract numeric version and build number for RC file
NUMERIC_VERSION=$(echo $VERSION | sed 's/-.*$//')
BUILD_DATE=$(echo $VERSION | sed -n 's/.*-dev\.\([0-9]\{8\}\)\..*$/\1/p')
BUILD_NUMBER="0"

if [ ! -z "$BUILD_DATE" ]; then
    # Use last 4 digits of date as build number (avoid overflow)
    BUILD_NUMBER=$(echo $BUILD_DATE | tail -c 5)  # Gets "0606"
fi

# Create version for RC file (format: major,minor,patch,build)
VERSION_RC=$(echo $NUMERIC_VERSION | sed 's/\./,/g'),$BUILD_NUMBER

# Create resource file with version information
cat > /tmp/synapseq.rc << EOF
#include <windows.h>

// Include icon
1 ICON "$BUILD_DIR/assets/synapseq.ico"

VS_VERSION_INFO VERSIONINFO
FILEVERSION     $VERSION_RC
PRODUCTVERSION  $VERSION_RC
FILEFLAGSMASK   VS_FFI_FILEFLAGSMASK
FILEFLAGS       0
FILEOS          VOS__WINDOWS32
FILETYPE        VFT_APP
FILESUBTYPE     0
BEGIN
    BLOCK "StringFileInfo"
    BEGIN
        BLOCK "040904E4"
        BEGIN
            VALUE "CompanyName",      "SynapSeq"
            VALUE "FileDescription",  "Synapse-Sequenced Brainwave Generator"
            VALUE "FileVersion",      "$VERSION"
            VALUE "InternalName",     "synapseq"
            VALUE "LegalCopyright",   "GPLv2"
            VALUE "OriginalFilename", "synapseq.exe"
            VALUE "ProductName",      "SynapSeq"
            VALUE "ProductVersion",   "$VERSION"
        END
    END
    BLOCK "VarFileInfo"
    BEGIN
        VALUE "Translation", 0x409, 1252
    END
END
EOF

# Compile resource file for both architectures
x86_64-w64-mingw32-windres /tmp/synapseq.rc -O coff -o /tmp/synapseq64.res

# Define paths for libraries. Change it to the correct path for your system.
LIBMAD_PATH_64="$BUILD_DIR/libs/libmad-win64.a"
LIBOGG_PATH_64="$BUILD_DIR/libs/libogg-win64.a"
LIBVORBIS_PATH_64="$BUILD_DIR/libs/libvorbis-win64.a"
LIBVORBISFILE_PATH_64="$BUILD_DIR/libs/libvorbisfile-win64.a"

# Build 64-bit version
section_header "Building 64-bit version..."

# Set up compilation flags for 64-bit
CFLAGS_64="-DT_MINGW -Wall -O3"
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

# Replace VERSION with the actual version number
sed "s/__VERSION__/\"$VERSION\"/" $SRC_DIR/synapseq.c > $SRC_DIR/synapseq.tmp.c

x86_64-w64-mingw32-gcc $CFLAGS_64 $SRC_DIR/synapseq.tmp.c /tmp/synapseq64.res -o $BUILD_DIR/dist/synapseq-win64.exe $LIBS_64

if [ $? -eq 0 ]; then
    success "64-bit compilation successful! Created 64-bit binary: $BUILD_DIR/dist/synapseq-win64.exe"
else
    error "64-bit compilation failed!"
fi

# Clean up temporary files
rm -f /tmp/synapseq.rc /tmp/synapseq64.res $SRC_DIR/synapseq.tmp.c

section_header "Build process completed!" 