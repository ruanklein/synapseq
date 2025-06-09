#!/bin/bash

# Source common library
. ./lib.sh

# Store the original directory
ORIGINAL_DIR="$PWD"

# Define variables for temporary paths
TEMP_DIR="/tmp/cross_compile"
INSTALL_DIR_WIN32="$TEMP_DIR/win32"
INSTALL_DIR_WIN64="$TEMP_DIR/win64"
SRC_DIR="$PWD/build"
OUTPUT_DIR="$TEMP_DIR/windows"
LIBOGG_VERSION="1.3.5"
LIBVORBISIDEC_VERSION="1.0.2+svn16259"
LIBMAD_VERSION="0.15.1b"

# Check for required tools
check_required_tools curl make automake autoconf libtool
check_required_tools i686-w64-mingw32-gcc x86_64-w64-mingw32-gcc

# Create temporary directories
section_header "Creating temporary directories..."
create_dir_if_not_exists "$INSTALL_DIR_WIN32"
create_dir_if_not_exists "$INSTALL_DIR_WIN64"
create_dir_if_not_exists "$SRC_DIR"
create_dir_if_not_exists "$OUTPUT_DIR"

# Toolchain settings for cross-compilation
WIN32_CC="i686-w64-mingw32-gcc"
WIN64_CC="x86_64-w64-mingw32-gcc"

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
section_header "Downloading libvorbisidec (Tremor)..."
if [ ! -f "libvorbisidec_$LIBVORBISIDEC_VERSION.orig.tar.gz" ]; then
    curl -L -o "libvorbisidec_$LIBVORBISIDEC_VERSION.orig.tar.gz" -s "https://launchpadlibrarian.net/35151187/libvorbisidec_$LIBVORBISIDEC_VERSION.orig.tar.gz" > /dev/null
    check_error "Failed to download libvorbisidec"
    info "Extracting libvorbisidec..."
    tar -xzf "libvorbisidec_$LIBVORBISIDEC_VERSION.orig.tar.gz" > /dev/null
    check_error "Failed to extract libvorbisidec"
    mv "libvorbisidec-$LIBVORBISIDEC_VERSION" "libvorbisidec_$LIBVORBISIDEC_VERSION"
fi
section_header "Downloading libmad..."
if [ ! -f "libmad-$LIBMAD_VERSION.tar.gz" ]; then
    curl -L -O -s "https://downloads.sourceforge.net/project/mad/libmad/$LIBMAD_VERSION/libmad-$LIBMAD_VERSION.tar.gz" > /dev/null
    check_error "Failed to download libmad"
    info "Extracting libmad..."
    tar -xzf "libmad-$LIBMAD_VERSION.tar.gz" > /dev/null
    check_error "Failed to extract libmad"
fi

# Update config.sub and config.guess for libmad to support MinGW
cd "$SRC_DIR/libmad-$LIBMAD_VERSION"
section_header "Updating config.sub and config.guess for libmad..."
# Copy config files from libs directory instead of downloading
cp "$PWD/../../libs/config.sub" ./config.sub
cp "$PWD/../../libs/config.guess" ./config.guess
chmod +x config.sub config.guess
check_error "Failed to update config.sub and config.guess for libmad"

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

# Compile libogg for Win32
cd "$SRC_DIR/libogg-$LIBOGG_VERSION"
info "Configuring libogg for Win32..."
./configure --host=i686-w64-mingw32 --prefix="$INSTALL_DIR_WIN32" --enable-static --disable-shared > "$TEMP_DIR/libogg-win32.log" 2>&1
check_error "libogg configuration for Win32 failed" "$TEMP_DIR/libogg-win32.log"
info "Compiling libogg for Win32..."
make clean >> "$TEMP_DIR/libogg-win32.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libogg-win32.log" 2>&1 && make install >> "$TEMP_DIR/libogg-win32.log" 2>&1
check_error "libogg compilation for Win32 failed" "$TEMP_DIR/libogg-win32.log"

# Compile libogg for Win64
info "Configuring libogg for Win64..."
./configure --host=x86_64-w64-mingw32 --prefix="$INSTALL_DIR_WIN64" --enable-static --disable-shared > "$TEMP_DIR/libogg-win64.log" 2>&1
check_error "libogg configuration for Win64 failed" "$TEMP_DIR/libogg-win64.log"
info "Compiling libogg for Win64..."
make clean >> "$TEMP_DIR/libogg-win64.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libogg-win64.log" 2>&1 && make install >> "$TEMP_DIR/libogg-win64.log" 2>&1
check_error "libogg compilation for Win64 failed" "$TEMP_DIR/libogg-win64.log"

# Copy required header files to libvorbisidec source directory
cd "$SRC_DIR/libvorbisidec_$LIBVORBISIDEC_VERSION"
section_header "Copying required header files to libvorbisidec source directory..."
cp "$ORIGINAL_DIR/libs/_G_config.h" ./
cp "$ORIGINAL_DIR/libs/os_types.h" ./
cp "$ORIGINAL_DIR/libs/ogg.h" ./
cp "$ORIGINAL_DIR/libs/config_types.h" ./
check_error "Failed to copy header files to libvorbisidec source directory"

# Generate configure for libvorbisidec (Tremor) using autogen.sh
info "Generating configure script for libvorbisidec..."
if [ ! -f "configure" ]; then
    ./autogen.sh > "$TEMP_DIR/libvorbisidec.log" 2>&1
    check_error "Failed to generate configure script for libvorbisidec (autogen.sh)" "$TEMP_DIR/libvorbisidec.log"
fi

# Compile libvorbisidec (Tremor) for Win32
info "Configuring libvorbisidec for Win32..."
CPPFLAGS="-I. -I$ORIGINAL_DIR/libs -I$INSTALL_DIR_WIN32/include" ./configure --host=i686-w64-mingw32 --prefix="$INSTALL_DIR_WIN32" --enable-static --disable-shared --with-ogg="$INSTALL_DIR_WIN32" > "$TEMP_DIR/libvorbisidec-win32.log" 2>&1
check_error "libvorbisidec configuration for Win32 failed" "$TEMP_DIR/libvorbisidec-win32.log"
info "Compiling libvorbisidec for Win32..."
make clean >> "$TEMP_DIR/libvorbisidec-win32.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libvorbisidec-win32.log" 2>&1 && make install >> "$TEMP_DIR/libvorbisidec-win32.log" 2>&1
check_error "libvorbisidec compilation for Win32 failed" "$TEMP_DIR/libvorbisidec-win32.log"

# Compile libvorbisidec (Tremor) for Win64
info "Configuring libvorbisidec for Win64..."
CPPFLAGS="-I. -I$ORIGINAL_DIR/libs -I$INSTALL_DIR_WIN64/include" ./configure --host=x86_64-w64-mingw32 --prefix="$INSTALL_DIR_WIN64" --enable-static --disable-shared --with-ogg="$INSTALL_DIR_WIN64" > "$TEMP_DIR/libvorbisidec-win64.log" 2>&1
check_error "libvorbisidec configuration for Win64 failed" "$TEMP_DIR/libvorbisidec-win64.log"
info "Compiling libvorbisidec for Win64..."
make clean >> "$TEMP_DIR/libvorbisidec-win64.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libvorbisidec-win64.log" 2>&1 && make install >> "$TEMP_DIR/libvorbisidec-win64.log" 2>&1
check_error "libvorbisidec compilation for Win64 failed" "$TEMP_DIR/libvorbisidec-win64.log"

# Compile libmad for Win32
cd "$SRC_DIR/libmad-$LIBMAD_VERSION"
info "Configuring libmad for Win32..."
./configure --host=i686-w64-mingw32 --prefix="$INSTALL_DIR_WIN32" --enable-static --disable-shared --disable-debugging > "$TEMP_DIR/libmad-win32.log" 2>&1
check_error "libmad configuration for Win32 failed" "$TEMP_DIR/libmad-win32.log"
info "Compiling libmad for Win32..."
make clean >> "$TEMP_DIR/libmad-win32.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libmad-win32.log" 2>&1 && make install >> "$TEMP_DIR/libmad-win32.log" 2>&1
check_error "libmad compilation for Win32 failed" "$TEMP_DIR/libmad-win32.log"

# Compile libmad for Win64
info "Configuring libmad for Win64..."
./configure --host=x86_64-w64-mingw32 --prefix="$INSTALL_DIR_WIN64" --enable-static --disable-shared --disable-debugging > "$TEMP_DIR/libmad-win64.log" 2>&1
check_error "libmad configuration for Win64 failed" "$TEMP_DIR/libmad-win64.log"
info "Compiling libmad for Win64..."
make clean >> "$TEMP_DIR/libmad-win64.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libmad-win64.log" 2>&1 && make install >> "$TEMP_DIR/libmad-win64.log" 2>&1
check_error "libmad compilation for Win64 failed" "$TEMP_DIR/libmad-win64.log"

# Copy libraries to output directory
section_header "Copying libraries to output directory..."
cp "$INSTALL_DIR_WIN32/lib/libogg.a" "$OUTPUT_DIR/libogg-win32.a"
cp "$INSTALL_DIR_WIN32/lib/libvorbisidec.a" "$OUTPUT_DIR/libvorbisidec-win32.a"
cp "$INSTALL_DIR_WIN32/lib/libmad.a" "$OUTPUT_DIR/libmad-win32.a"
cp "$INSTALL_DIR_WIN64/lib/libogg.a" "$OUTPUT_DIR/libogg-win64.a"
cp "$INSTALL_DIR_WIN64/lib/libvorbisidec.a" "$OUTPUT_DIR/libvorbisidec-win64.a"
cp "$INSTALL_DIR_WIN64/lib/libmad.a" "$OUTPUT_DIR/libmad-win64.a"
check_error "Failed to copy libraries to output directory"

# Copy libraries to the project's libs directory
section_header "Copying libraries to project's libs directory..."
cd "$ORIGINAL_DIR"
copy_libs "$OUTPUT_DIR" "libs" "windows"

# Cleanup: Remove temporary directories and logs, keeping only the libraries
section_header "Cleaning up temporary directories and logs..."
rm -rf "$INSTALL_DIR_WIN32" "$INSTALL_DIR_WIN64" "$SRC_DIR" "$TEMP_DIR"/*.log
success "Compilation completed! Libraries have been copied to the libs folder." 