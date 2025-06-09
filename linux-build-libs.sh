#!/bin/bash

# Source common library
. ./lib.sh

# Store the original directory
ORIGINAL_DIR="$PWD"

# Define variables for temporary paths
TEMP_DIR="/tmp/cross_compile"
INSTALL_DIR_X86="$TEMP_DIR/x86"
INSTALL_DIR_X86_64="$TEMP_DIR/x86_64"
INSTALL_DIR_ARM64="$TEMP_DIR/arm64"
SRC_DIR="$PWD/build"
OUTPUT_DIR="$TEMP_DIR/linux"
LIBOGG_VERSION="1.3.5"
LIBVORBISIDEC_VERSION="1.0.2+svn16259"
LIBMAD_VERSION="0.15.1b"

# Check for required tools
check_required_tools gcc g++ curl make automake autoconf libtool

# Create temporary directories
section_header "Creating temporary directories..."
create_dir_if_not_exists "$INSTALL_DIR_X86"
create_dir_if_not_exists "$INSTALL_DIR_X86_64"
create_dir_if_not_exists "$INSTALL_DIR_ARM64"
create_dir_if_not_exists "$SRC_DIR"
create_dir_if_not_exists "$OUTPUT_DIR"

# Detect host architecture
HOST_ARCH=$(uname -m)
info "Detected host architecture: $HOST_ARCH"

# Toolchain settings for cross-compilation
X86_CC="gcc -m32"
X86_64_CC="gcc -m64"
ARM64_CC="gcc"

SKIP_32BIT=0
SKIP_ARM64=0

if [ "$HOST_ARCH" = "aarch64" ]; then
    info "Running on ARM64, using native gcc for 64-bit compilation"
    # Skip 32-bit compilation on ARM64 as it's not typically supported
    SKIP_32BIT=1
    warning "32-bit compilation is not supported on ARM64, skipping..."
fi

# Avoid ARM64 compilation if not on ARM64 platform
if [ "$HOST_ARCH" != "aarch64" ]; then
    SKIP_ARM64=1
    info "Not on ARM64 platform, skipping ARM64 compilation..."
fi

# Check if 32-bit compilation is supported (only if not on ARM64)
if [ $SKIP_32BIT = 0 ]; then
    section_header "Checking for 32-bit compilation support..."
    if command_exists "gcc"; then
        # Create a simple test program - using single line to avoid newline issues
        echo 'int main() { return 0; }' > /tmp/test32.c
        
        # Try to compile with -m32
        if gcc -m32 -o /tmp/test32 /tmp/test32.c >/dev/null 2>&1; then
            success "32-bit compilation is supported!"
            rm -f /tmp/test32.c /tmp/test32
        else
            warning "32-bit compilation is not supported. Installing multilib packages might be required."
            warning "On Debian/Ubuntu: sudo apt-get install gcc-multilib g++-multilib"
            warning "On Fedora: sudo dnf install glibc-devel.i686 libstdc++-devel.i686"
            warning "On Arch: sudo pacman -S multilib-devel"
            
            # Ask if the user wants to continue with 64-bit only
            echo -e "${YELLOW}==> Do you want to continue with 64-bit compilation only? (y/n)${NC}"
            read -r response
            if [[ "$response" =~ ^([nN][oO]|[nN])$ ]]; then
                error "Aborting compilation. Please install the required packages and try again."
                exit 1
            else
                warning "Continuing with 64-bit compilation only..."
                # Set a flag to skip 32-bit compilation
                SKIP_32BIT=1
            fi
        fi
    else
        error "GCC not found. Please install GCC and try again."
        exit 1
    fi
fi

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

# Update config.sub and config.guess for libmad to support modern Linux
cd "$SRC_DIR/libmad-$LIBMAD_VERSION"
section_header "Updating config.sub and config.guess for libmad..."
# Copy config files from libs directory instead of downloading
cp "$PWD/../../libs/config.sub" ./config.sub
cp "$PWD/../../libs/config.guess" ./config.guess
chmod +x config.sub config.guess
check_error "Failed to update config.sub and config.guess for libmad"

# Patch libmad configure to remove unsupported gcc flags
section_header "Patching libmad configure script..."
sed -i 's/-fforce-mem//g' configure
sed -i 's/-fthread-jumps//g' configure
sed -i 's/-fcse-follow-jumps//g' configure
sed -i 's/-fcse-skip-blocks//g' configure
sed -i 's/-fregmove//g' configure
sed -i 's/-fexpensive-optimizations//g' configure
sed -i 's/-fschedule-insns2//g' configure
check_error "Failed to patch libmad configure script"

# Compile libogg for x86
cd "$SRC_DIR/libogg-$LIBOGG_VERSION"
if [ $SKIP_32BIT = 0 ]; then
    info "Configuring libogg for x86..."
    ./configure CC="$X86_CC" --prefix="$INSTALL_DIR_X86" --enable-static --disable-shared > "$TEMP_DIR/libogg-x86.log" 2>&1
    check_error "libogg configuration for x86 failed" "$TEMP_DIR/libogg-x86.log"
    info "Compiling libogg for x86..."
    make clean >> "$TEMP_DIR/libogg-x86.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libogg-x86.log" 2>&1 && make install >> "$TEMP_DIR/libogg-x86.log" 2>&1
    check_error "libogg compilation for x86 failed" "$TEMP_DIR/libogg-x86.log"
else
    warning "Skipping 32-bit compilation for libogg..."
fi



# Compile libogg for arm64 or x86_64
if [ $SKIP_ARM64 = 0 ]; then
    info "Configuring libogg for arm64..."
    ./configure CC="$ARM64_CC" --prefix="$INSTALL_DIR_ARM64" --enable-static --disable-shared > "$TEMP_DIR/libogg-arm64.log" 2>&1
    check_error "libogg configuration for arm64 failed" "$TEMP_DIR/libogg-arm64.log"
    info "Compiling libogg for arm64..."
    make clean >> "$TEMP_DIR/libogg-arm64.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libogg-arm64.log" 2>&1 && make install >> "$TEMP_DIR/libogg-arm64.log" 2>&1
    check_error "libogg compilation for arm64 failed" "$TEMP_DIR/libogg-arm64.log"
else 
    info "Configuring libogg for x86_64..."
    ./configure CC="$X86_64_CC" --prefix="$INSTALL_DIR_X86_64" --enable-static --disable-shared > "$TEMP_DIR/libogg-x86_64.log" 2>&1
    check_error "libogg configuration for x86_64 failed" "$TEMP_DIR/libogg-x86_64.log"
    info "Compiling libogg for x86_64..."
    make clean >> "$TEMP_DIR/libogg-x86_64.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libogg-x86_64.log" 2>&1 && make install >> "$TEMP_DIR/libogg-x86_64.log" 2>&1
    check_error "libogg compilation for x86_64 failed" "$TEMP_DIR/libogg-x86_64.log"
fi

# Generate configure for libvorbisidec (Tremor) using autogen.sh
cd "$SRC_DIR/libvorbisidec_$LIBVORBISIDEC_VERSION"
info "Generating configure script for libvorbisidec..."
if [ ! -f "configure" ]; then
    ./autogen.sh > "$TEMP_DIR/libvorbisidec.log" 2>&1
    check_error "Failed to generate configure script for libvorbisidec (autogen.sh)" "$TEMP_DIR/libvorbisidec.log"
fi

# Compile libvorbisidec (Tremor) for x86
if [ $SKIP_32BIT = 0 ]; then
    info "Configuring libvorbisidec for x86..."
    ./configure CC="$X86_CC" --prefix="$INSTALL_DIR_X86" --enable-static --disable-shared --with-ogg="$INSTALL_DIR_X86" > "$TEMP_DIR/libvorbisidec-x86.log" 2>&1
    check_error "libvorbisidec configuration for x86 failed" "$TEMP_DIR/libvorbisidec-x86.log"
    info "Compiling libvorbisidec for x86..."
    make clean >> "$TEMP_DIR/libvorbisidec-x86.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libvorbisidec-x86.log" 2>&1 && make install >> "$TEMP_DIR/libvorbisidec-x86.log" 2>&1
    check_error "libvorbisidec compilation for x86 failed" "$TEMP_DIR/libvorbisidec-x86.log"
else
    warning "Skipping 32-bit compilation for libvorbisidec..."
fi

# Compile libvorbisidec (Tremor) for arm64 or x86_64
if [ $SKIP_ARM64 = 0 ]; then
    info "Configuring libvorbisidec for arm64..."
    ./configure CC="$ARM64_CC" --prefix="$INSTALL_DIR_ARM64" --enable-static --disable-shared --with-ogg="$INSTALL_DIR_ARM64" > "$TEMP_DIR/libvorbisidec-arm64.log" 2>&1
    check_error "libvorbisidec configuration for arm64 failed" "$TEMP_DIR/libvorbisidec-arm64.log"
    info "Compiling libvorbisidec for arm64..."
    make clean >> "$TEMP_DIR/libvorbisidec-arm64.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libvorbisidec-arm64.log" 2>&1 && make install >> "$TEMP_DIR/libvorbisidec-arm64.log" 2>&1
    check_error "libvorbisidec compilation for arm64 failed" "$TEMP_DIR/libvorbisidec-arm64.log"
else
    info "Configuring libvorbisidec for x86_64..."
    ./configure CC="$X86_64_CC" --prefix="$INSTALL_DIR_X86_64" --enable-static --disable-shared --with-ogg="$INSTALL_DIR_X86_64" > "$TEMP_DIR/libvorbisidec-x86_64.log" 2>&1
    check_error "libvorbisidec configuration for x86_64 failed" "$TEMP_DIR/libvorbisidec-x86_64.log"
    info "Compiling libvorbisidec for x86_64..."
    make clean >> "$TEMP_DIR/libvorbisidec-x86_64.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libvorbisidec-x86_64.log" 2>&1 && make install >> "$TEMP_DIR/libvorbisidec-x86_64.log" 2>&1
    check_error "libvorbisidec compilation for x86_64 failed" "$TEMP_DIR/libvorbisidec-x86_64.log"
fi

# Compile libmad for x86
cd "$SRC_DIR/libmad-$LIBMAD_VERSION"
if [ $SKIP_32BIT = 0 ]; then
    info "Configuring libmad for x86..."
    ./configure CC="$X86_CC" --prefix="$INSTALL_DIR_X86" --enable-static --disable-shared --disable-debugging > "$TEMP_DIR/libmad-x86.log" 2>&1
    check_error "libmad configuration for x86 failed" "$TEMP_DIR/libmad-x86.log"
    info "Compiling libmad for x86..."
    make clean >> "$TEMP_DIR/libmad-x86.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libmad-x86.log" 2>&1 && make install >> "$TEMP_DIR/libmad-x86.log" 2>&1
    check_error "libmad compilation for x86 failed" "$TEMP_DIR/libmad-x86.log"
else
    warning "Skipping 32-bit compilation for libmad..."
fi

# Compile libmad for arm64 or x86_64
if [ $SKIP_ARM64 = 0 ]; then
    info "Configuring libmad for arm64..."
    ./configure CC="$ARM64_CC" --prefix="$INSTALL_DIR_ARM64" --enable-static --disable-shared --disable-debugging > "$TEMP_DIR/libmad-arm64.log" 2>&1
    check_error "libmad configuration for arm64 failed" "$TEMP_DIR/libmad-arm64.log"
    info "Compiling libmad for arm64..."
    make clean >> "$TEMP_DIR/libmad-arm64.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libmad-arm64.log" 2>&1 && make install >> "$TEMP_DIR/libmad-arm64.log" 2>&1
    check_error "libmad compilation for arm64 failed" "$TEMP_DIR/libmad-arm64.log"
else 
    info "Configuring libmad for x86_64..."
    ./configure CC="$X86_64_CC" --prefix="$INSTALL_DIR_X86_64" --enable-static --disable-shared --disable-debugging > "$TEMP_DIR/libmad-x86_64.log" 2>&1
    check_error "libmad configuration for x86_64 failed" "$TEMP_DIR/libmad-x86_64.log"
    info "Compiling libmad for x86_64..."
    make clean >> "$TEMP_DIR/libmad-x86_64.log" 2>&1 && make -j$(nproc) >> "$TEMP_DIR/libmad-x86_64.log" 2>&1 && make install >> "$TEMP_DIR/libmad-x86_64.log" 2>&1
    check_error "libmad compilation for x86_64 failed" "$TEMP_DIR/libmad-x86_64.log"
fi

# Copy libraries to output directory
section_header "Copying libraries to output directory..."
if [ $SKIP_32BIT = 0 ]; then
    cp "$INSTALL_DIR_X86/lib/libogg.a" "$OUTPUT_DIR/libogg-x86.a"
    cp "$INSTALL_DIR_X86/lib/libvorbisidec.a" "$OUTPUT_DIR/libvorbisidec-x86.a"
    cp "$INSTALL_DIR_X86/lib/libmad.a" "$OUTPUT_DIR/libmad-x86.a"
fi

if [ $SKIP_ARM64 = 0 ]; then
    cp "$INSTALL_DIR_ARM64/lib/libogg.a" "$OUTPUT_DIR/libogg-arm64.a"
    cp "$INSTALL_DIR_ARM64/lib/libvorbisidec.a" "$OUTPUT_DIR/libvorbisidec-arm64.a"
    cp "$INSTALL_DIR_ARM64/lib/libmad.a" "$OUTPUT_DIR/libmad-arm64.a"
else
    cp "$INSTALL_DIR_X86_64/lib/libogg.a" "$OUTPUT_DIR/libogg-x86_64.a"
    cp "$INSTALL_DIR_X86_64/lib/libvorbisidec.a" "$OUTPUT_DIR/libvorbisidec-x86_64.a"
    cp "$INSTALL_DIR_X86_64/lib/libmad.a" "$OUTPUT_DIR/libmad-x86_64.a"
fi
check_error "Failed to copy libraries to output directory"

# Copy libraries to the project's libs directory
section_header "Copying libraries to project's libs directory..."
cd "$ORIGINAL_DIR"
copy_libs "$OUTPUT_DIR" "libs" "linux"

# Cleanup: Remove temporary directories and logs, keeping only the libraries
section_header "Cleaning up temporary directories and logs..."
rm -rf "$INSTALL_DIR_X86" "$INSTALL_DIR_X86_64" "$INSTALL_DIR_ARM64" "$SRC_DIR" "$TEMP_DIR"/*.log
success "Compilation completed! Libraries have been copied to the libs folder." 