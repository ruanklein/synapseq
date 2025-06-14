#!/bin/bash

# Source common library
. ./lib.sh

# Store the original directory
ORIGINAL_DIR="$PWD"

# Define variables for temporary paths
TEMP_DIR="/tmp/cross_compile"
INSTALL_DIR_WIN64="$TEMP_DIR/win64"
SRC_DIR="$PWD/build"
OUTPUT_DIR="$TEMP_DIR/windows"
LIBOGG_VERSION="1.3.5"
LIBVORBIS_VERSION="1.3.7"
LIBMAD_VERSION="0.15.1b"

# Create temporary directories
section_header "Creating temporary directories..."
create_dir_if_not_exists "$INSTALL_DIR_WIN64"
create_dir_if_not_exists "$SRC_DIR"
create_dir_if_not_exists "$OUTPUT_DIR"

# Create libs directory if it doesn't exist
create_dir_if_not_exists "libs"

# Download libraries if not present
cd "$SRC_DIR"
section_header "Downloading libogg..."
if [ ! -f "libogg-$LIBOGG_VERSION.tar.gz" ]; then
    curl -L -O -s "https://downloads.xiph.org/releases/ogg/libogg-$LIBOGG_VERSION.tar.gz" > /dev/null
    check_error "Failed to download libogg"
    info "Extracting libogg..."
    tar -xzf "libogg-$LIBOGG_VERSION.tar.gz" > /dev/null
    check_error "Failed to extract libogg"
fi
section_header "Downloading libvorbis..."
if [ ! -f "libvorbis-$LIBVORBIS_VERSION.tar.gz" ]; then
    curl -L -O -s "https://downloads.xiph.org/releases/vorbis/libvorbis-$LIBVORBIS_VERSION.tar.gz" > /dev/null
    check_error "Failed to download libvorbis"
    info "Extracting libvorbis..."
    tar -xzf "libvorbis-$LIBVORBIS_VERSION.tar.gz" > /dev/null
    check_error "Failed to extract libvorbis"
fi
section_header "Downloading libmad..."
if [ ! -f "libmad-$LIBMAD_VERSION.tar.gz" ]; then
    curl -L -O -s "https://downloads.sourceforge.net/project/mad/libmad/$LIBMAD_VERSION/libmad-$LIBMAD_VERSION.tar.gz" > /dev/null
    check_error "Failed to download libmad"
    info "Extracting libmad..."
    tar -xzf "libmad-$LIBMAD_VERSION.tar.gz" > /dev/null
    check_error "Failed to extract libmad"
fi

# Compile libogg for Win64
cd "$SRC_DIR/libogg-$LIBOGG_VERSION"
info "Configuring libogg for Win64..."
./configure --host=x86_64-w64-mingw32 \
    --prefix="$INSTALL_DIR_WIN64" \
    --enable-static \
    --disable-shared \
    > "$TEMP_DIR/libogg-win64.log" 2>&1
check_error "libogg configuration for Win64 failed" "$TEMP_DIR/libogg-win64.log"

info "Compiling libogg for Win64..."
make -j$(nproc) >> "$TEMP_DIR/libogg-win64.log" 2>&1 && make install >> "$TEMP_DIR/libogg-win64.log" 2>&1
check_error "libogg compilation for Win64 failed" "$TEMP_DIR/libogg-win64.log"

# Compile libvorbis for Win64
cd "$SRC_DIR/libvorbis-$LIBVORBIS_VERSION"
info "Configuring libvorbis for Win64..."
CPPFLAGS="-I${INSTALL_DIR_WIN64}/include" \
LDFLAGS="-L${INSTALL_DIR_WIN64}/lib" \
./configure --host=x86_64-w64-mingw32 \
    --prefix="$INSTALL_DIR_WIN64" \
    --with-ogg-includes="${INSTALL_DIR_WIN64}/include" \
    --with-ogg-libraries="${INSTALL_DIR_WIN64}/lib" \
    --enable-static \
    --disable-shared \
    > "$TEMP_DIR/libvorbis-win64.log" 2>&1
check_error "libvorbis configuration for Win64 failed" "$TEMP_DIR/libvorbis-win64.log"

info "Compiling libvorbis for Win64..."
make -j$(nproc) >> "$TEMP_DIR/libvorbis-win64.log" 2>&1 && make install >> "$TEMP_DIR/libvorbis-win64.log" 2>&1
check_error "libvorbis compilation for Win64 failed" "$TEMP_DIR/libvorbis-win64.log"

# Compile libmad for Win64
cd "$SRC_DIR/libmad-$LIBMAD_VERSION"

# Patch libmad configure to remove unsupported gcc flags
section_header "Patching libmad configure script..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS version of sed
    sed -i '' 's/-fforce-mem//g' configure
    sed -i '' 's/-fthread-jumps//g' configure
    sed -i '' 's/-fcse-follow-jumps//g' configure
    sed -i '' 's/-fcse-skip-blocks//g' configure
    sed -i '' 's/-fregmove//g' configure
    sed -i '' 's/-fexpensive-optimizations//g' configure
    sed -i '' 's/-fschedule-insns2//g' configure
else
    # Linux version of sed
    sed -i 's/-fforce-mem//g' configure
    sed -i 's/-fthread-jumps//g' configure
    sed -i 's/-fcse-follow-jumps//g' configure
    sed -i 's/-fcse-skip-blocks//g' configure
    sed -i 's/-fregmove//g' configure
    sed -i 's/-fexpensive-optimizations//g' configure
    sed -i 's/-fschedule-insns2//g' configure
fi
check_error "Failed to patch libmad configure script"

info "Configuring libmad for Win64..."
./configure --host=x86_64-w64-mingw32 \
    --prefix="$INSTALL_DIR_WIN64" \
    --enable-static \
    --disable-shared \
    --disable-debugging \
    > "$TEMP_DIR/libmad-win64.log" 2>&1
check_error "libmad configuration for Win64 failed" "$TEMP_DIR/libmad-win64.log"

info "Compiling libmad for Win64..."
make -j$(nproc) >> "$TEMP_DIR/libmad-win64.log" 2>&1 && make install >> "$TEMP_DIR/libmad-win64.log" 2>&1
check_error "libmad compilation for Win64 failed" "$TEMP_DIR/libmad-win64.log"

# Copy libraries to output directory
section_header "Copying libraries and headers to output directory..."
cp "$INSTALL_DIR_WIN64/lib/libogg.a" "$OUTPUT_DIR/libogg-win64.a"
cp "$INSTALL_DIR_WIN64/lib/libvorbis.a" "$OUTPUT_DIR/libvorbis-win64.a"
cp "$INSTALL_DIR_WIN64/lib/libvorbisfile.a" "$OUTPUT_DIR/libvorbisfile-win64.a"
cp "$INSTALL_DIR_WIN64/lib/libvorbisenc.a" "$OUTPUT_DIR/libvorbisenc-win64.a"
cp "$INSTALL_DIR_WIN64/lib/libmad.a" "$OUTPUT_DIR/libmad-win64.a"
cp -R "$INSTALL_DIR_WIN64/include" "$OUTPUT_DIR/include"
check_error "Failed to copy libraries and headers to output directory"

# Copy libraries to the project's libs directory
section_header "Copying libraries and headers to project's libs directory..."
cd "$ORIGINAL_DIR"
cp -R $OUTPUT_DIR/* $ORIGINAL_DIR/libs

# Cleanup: Remove temporary directories and logs, keeping only the libraries
section_header "Cleaning up temporary directories and logs..."
rm -rf "$INSTALL_DIR_WIN64" "$SRC_DIR" "$TEMP_DIR"/*.log
success "Compilation completed! Libraries have been copied to the libs folder." 