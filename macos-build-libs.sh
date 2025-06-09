#!/bin/bash

# Source common library
. ./lib.sh

# Store the original directory
ORIGINAL_DIR="$PWD"

# Define variables for temporary paths
TEMP_DIR="/tmp/cross_compile"
INSTALL_DIR_ARM64="$TEMP_DIR/arm64"
INSTALL_DIR_X86_64="$TEMP_DIR/x86_64"
SRC_DIR="$PWD/build"
UNIVERSAL_DIR="$TEMP_DIR/universal"
LIBOGG_VERSION="1.3.5"         
LIBVORBISIDEC_VERSION="1.0.2+svn16259" 
LIBMAD_VERSION="0.15.1b"       

# Define macOS SDK path and minimum version (Big Sur 11.0 for ARM64 support)
SDK_PATH=$(xcrun --sdk macosx --show-sdk-path)
MACOS_VERSION_MIN="11.0"

# Check for required tools
check_required_tools clang curl lipo make automake autoconf libtool

# Create temporary directories
section_header "Creating temporary directories..."
create_dir_if_not_exists "$INSTALL_DIR_ARM64"
create_dir_if_not_exists "$INSTALL_DIR_X86_64"
create_dir_if_not_exists "$SRC_DIR"
create_dir_if_not_exists "$UNIVERSAL_DIR"

# Toolchain settings
ARM64_CC="clang -arch arm64 -isysroot $SDK_PATH -mmacosx-version-min=$MACOS_VERSION_MIN"
X86_64_CC="clang -arch x86_64 -isysroot $SDK_PATH -mmacosx-version-min=$MACOS_VERSION_MIN"

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

# Update config.sub and config.guess for libmad to support modern macOS
cd "$SRC_DIR/libmad-$LIBMAD_VERSION"
section_header "Updating config.sub and config.guess for libmad..."
# Copy config files from libs directory instead of downloading
cp "$PWD/../../libs/config.sub" ./config.sub
cp "$PWD/../../libs/config.guess" ./config.guess
chmod +x config.sub config.guess
check_error "Failed to update config.sub and config.guess for libmad"

# Patch libmad configure to remove unsupported clang flags
section_header "Patching libmad configure script..."
sed -i '' 's/-fforce-mem//g' configure
sed -i '' 's/-fthread-jumps//g' configure
sed -i '' 's/-fcse-follow-jumps//g' configure
sed -i '' 's/-fcse-skip-blocks//g' configure
sed -i '' 's/-fregmove//g' configure
sed -i '' 's/-fexpensive-optimizations//g' configure
sed -i '' 's/-fschedule-insns2//g' configure
check_error "Failed to patch libmad configure script"

# Compile libogg for arm64
cd "$SRC_DIR/libogg-$LIBOGG_VERSION"
info "Configuring libogg for arm64..."
./configure CC="$ARM64_CC" --prefix="$INSTALL_DIR_ARM64" --enable-static --disable-shared > "$TEMP_DIR/libogg.log" 2>&1
check_error "libogg configuration for arm64 failed" "$TEMP_DIR/libogg.log"
info "Compiling libogg for arm64..."
make clean >> "$TEMP_DIR/libogg.log" 2>&1 && make -j$(sysctl -n hw.ncpu) >> "$TEMP_DIR/libogg.log" 2>&1 && make install >> "$TEMP_DIR/libogg.log" 2>&1
check_error "libogg compilation for arm64 failed" "$TEMP_DIR/libogg.log"

# Compile libogg for x86_64
info "Configuring libogg for x86_64..."
./configure CC="$X86_64_CC" --prefix="$INSTALL_DIR_X86_64" --enable-static --disable-shared >> "$TEMP_DIR/libogg.log" 2>&1
check_error "libogg configuration for x86_64 failed" "$TEMP_DIR/libogg.log"
info "Compiling libogg for x86_64..."
make clean >> "$TEMP_DIR/libogg.log" 2>&1 && make -j$(sysctl -n hw.ncpu) >> "$TEMP_DIR/libogg.log" 2>&1 && make install >> "$TEMP_DIR/libogg.log" 2>&1
check_error "libogg compilation for x86_64 failed" "$TEMP_DIR/libogg.log"

# Generate configure for libvorbisidec (Tremor) using autogen.sh
cd "$SRC_DIR/libvorbisidec_$LIBVORBISIDEC_VERSION"
info "Generating configure script for libvorbisidec..."
if [ ! -f "configure" ]; then
    ./autogen.sh > "$TEMP_DIR/libvorbisidec.log" 2>&1
    check_error "Failed to generate configure script for libvorbisidec (autogen.sh)" "$TEMP_DIR/libvorbisidec.log"
fi

# Compile libvorbisidec (Tremor) for arm64
info "Configuring libvorbisidec for arm64..."
./configure CC="$ARM64_CC" --prefix="$INSTALL_DIR_ARM64" --enable-static --disable-shared --with-ogg="$INSTALL_DIR_ARM64" >> "$TEMP_DIR/libvorbisidec.log" 2>&1
check_error "libvorbisidec configuration for arm64 failed" "$TEMP_DIR/libvorbisidec.log"
info "Compiling libvorbisidec for arm64..."
make clean >> "$TEMP_DIR/libvorbisidec.log" 2>&1 && make -j$(sysctl -n hw.ncpu) >> "$TEMP_DIR/libvorbisidec.log" 2>&1 && make install >> "$TEMP_DIR/libvorbisidec.log" 2>&1
check_error "libvorbisidec compilation for arm64 failed" "$TEMP_DIR/libvorbisidec.log"

# Compile libvorbisidec (Tremor) for x86_64
info "Configuring libvorbisidec for x86_64..."
./configure CC="$X86_64_CC" --prefix="$INSTALL_DIR_X86_64" --enable-static --disable-shared --with-ogg="$INSTALL_DIR_X86_64" >> "$TEMP_DIR/libvorbisidec.log" 2>&1
check_error "libvorbisidec configuration for x86_64 failed" "$TEMP_DIR/libvorbisidec.log"
info "Compiling libvorbisidec for x86_64..."
make clean >> "$TEMP_DIR/libvorbisidec.log" 2>&1 && make -j$(sysctl -n hw.ncpu) >> "$TEMP_DIR/libvorbisidec.log" 2>&1 && make install >> "$TEMP_DIR/libvorbisidec.log" 2>&1
check_error "libvorbisidec compilation for x86_64 failed" "$TEMP_DIR/libvorbisidec.log"

# Compile libmad for arm64
cd "$SRC_DIR/libmad-$LIBMAD_VERSION"
info "Configuring libmad for arm64..."
./configure CC="$ARM64_CC" --prefix="$INSTALL_DIR_ARM64" --enable-static --disable-shared --disable-debugging > "$TEMP_DIR/libmad.log" 2>&1
check_error "libmad configuration for arm64 failed" "$TEMP_DIR/libmad.log"
info "Compiling libmad for arm64..."
make clean >> "$TEMP_DIR/libmad.log" 2>&1 && make -j$(sysctl -n hw.ncpu) >> "$TEMP_DIR/libmad.log" 2>&1 && make install >> "$TEMP_DIR/libmad.log" 2>&1
check_error "libmad compilation for arm64 failed" "$TEMP_DIR/libmad.log"

# Compile libmad for x86_64
info "Configuring libmad for x86_64..."
./configure CC="$X86_64_CC" --prefix="$INSTALL_DIR_X86_64" --enable-static --disable-shared --disable-debugging >> "$TEMP_DIR/libmad.log" 2>&1
check_error "libmad configuration for x86_64 failed" "$TEMP_DIR/libmad.log"
info "Compiling libmad for x86_64..."
make clean >> "$TEMP_DIR/libmad.log" 2>&1 && make -j$(sysctl -n hw.ncpu) >> "$TEMP_DIR/libmad.log" 2>&1 && make install >> "$TEMP_DIR/libmad.log" 2>&1
check_error "libmad compilation for x86_64 failed" "$TEMP_DIR/libmad.log"

# Create universal libraries with lipo
cd "$UNIVERSAL_DIR"
section_header "Creating universal library for libogg..."
lipo -create "$INSTALL_DIR_ARM64/lib/libogg.a" "$INSTALL_DIR_X86_64/lib/libogg.a" -output "libogg.a"
check_error "Creation of universal libogg.a failed"
section_header "Creating universal library for libvorbisidec..."
lipo -create "$INSTALL_DIR_ARM64/lib/libvorbisidec.a" "$INSTALL_DIR_X86_64/lib/libvorbisidec.a" -output "libvorbisidec.a"
check_error "Creation of universal libvorbisidec.a failed"
section_header "Creating universal library for libmad..."
lipo -create "$INSTALL_DIR_ARM64/lib/libmad.a" "$INSTALL_DIR_X86_64/lib/libmad.a" -output "libmad.a"
check_error "Creation of universal libmad.a failed"

# Verify architectures in universal libraries
# section_header "Verifying architectures in libogg.a:"
# lipo -info "libogg.a"
# section_header "Verifying architectures in libvorbisidec.a:"
# lipo -info "libvorbisidec.a"
# section_header "Verifying architectures in libmad.a:"
# lipo -info "libmad.a"

# Copy libraries to the project's libs directory
section_header "Copying universal libraries to project's libs directory..."
cd "$ORIGINAL_DIR"
copy_libs "$UNIVERSAL_DIR" "libs" "macos-universal"

# Cleanup: Remove temporary directories and logs, keeping only universal .a files
section_header "Cleaning up temporary directories and logs..."
rm -rf "$INSTALL_DIR_ARM64" "$INSTALL_DIR_X86_64" "$SRC_DIR" "$TEMP_DIR/libogg.log" "$TEMP_DIR/libvorbisidec.log" "$TEMP_DIR/libmad.log"
success "Compilation completed! Universal libraries have been copied to the libs folder."