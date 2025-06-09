#!/bin/bash

# Common library for SynapSeq build scripts
# Contains shared functions and variables for all build scripts

# Define colors for terminal output
export GREEN='\033[1;32m'
export BLUE='\033[1;34m'
export CYAN='\033[1;36m'
export YELLOW='\033[1;33m'
export RED='\033[1;31m'
export NC='\033[0m' # No Color

# Function to check for errors and exit if any
check_error() {
    if [ $? -ne 0 ]; then
        echo -e "${RED}==> Error: $1${NC}"
        if [ ! -z "$2" ]; then
            echo -e "${RED}==> Check the log file '$2' for details.${NC}"
        fi
        exit 1
    fi
}

# Function to create directory if it doesn't exist
create_dir_if_not_exists() {
    if [ ! -d "$1" ]; then
        echo -e "${YELLOW}==> Creating directory: $1${NC}"
        mkdir -p "$1"
        check_error "Failed to create directory: $1"
    fi
}

# Function to display section header
section_header() {
    echo -e "${CYAN}==> $1${NC}"
}

# Function to display warning
warning() {
    echo -e "${YELLOW}==> Warning: $1${NC}"
}

# Function to display error
error() {
    echo -e "${RED}==> Error: $1${NC}"
}

# Function to display success
success() {
    echo -e "${GREEN}==> $1${NC}"
}

# Function to display info
info() {
    echo -e "${BLUE}==> $1${NC}"
}

# Function to check if a command exists
command_exists() {
    # For simple commands without arguments
    if [[ "$1" != *" "* ]]; then
        command -v "$1" >/dev/null 2>&1
        return $?
    fi
    
    # For commands with arguments (like "gcc -m32")
    local cmd=$(echo "$1" | cut -d' ' -f1)
    
    # First check if the base command exists
    if ! command -v "$cmd" >/dev/null 2>&1; then
        return 1
    fi
    
    # For gcc with specific flags, try a simple check
    # without trying to compile, just checking if the command returns an error code different from 127
    if [[ "$cmd" == "gcc" || "$cmd" == "g++" || "$cmd" == "clang" ]]; then
        $1 --version >/dev/null 2>&1
        local ret=$?
        # 127 means command not found, any other code is acceptable
        if [ $ret -eq 127 ]; then
            return 1
        else
            return 0
        fi
    fi
    
    # For other commands with arguments, assume they work if the base command exists
    return 0
}

# Function to check if required tools are installed
check_required_tools() {
    for tool in "$@"; do
        # Special case for libtool which might be called libtoolize in some distributions
        if [ "$tool" = "libtool" ]; then
            if ! command_exists "libtool" && ! command_exists "libtoolize"; then
                error "Required tool '$tool' is not installed. Please install it and try again."
                info "On Debian/Ubuntu: sudo apt-get install libtool"
                info "On Fedora: sudo dnf install libtool"
                info "On Arch: sudo pacman -S libtool"
                exit 1
            fi
            continue
        fi
        
        # Check for other tools
        if ! command_exists "$tool"; then
            error "Required tool '$tool' is not installed. Please install it and try again."
            case "$tool" in
                gcc|g++)
                    info "On Debian/Ubuntu: sudo apt-get install build-essential"
                    info "On Fedora: sudo dnf install gcc gcc-c++"
                    info "On Arch: sudo pacman -S base-devel"
                    ;;
                clang|lipo)
                    info "On macOS: Install Xcode Command Line Tools with 'xcode-select --install'"
                    ;;
                curl)
                    info "On Debian/Ubuntu: sudo apt-get install curl"
                    info "On Fedora: sudo dnf install curl"
                    info "On Arch: sudo pacman -S curl"
                    info "On macOS: brew install curl"
                    ;;
                unzip)
                    info "On Debian/Ubuntu: sudo apt-get install unzip"
                    info "On Fedora: sudo dnf install unzip"
                    info "On Arch: sudo pacman -S unzip"
                    info "On macOS: brew install unzip"
                    ;;
                make|automake|autoconf)
                    info "On Debian/Ubuntu: sudo apt-get install build-essential automake autoconf"
                    info "On Fedora: sudo dnf install make automake autoconf"
                    info "On Arch: sudo pacman -S base-devel"
                    info "On macOS: brew install automake autoconf"
                    ;;
                i686-w64-mingw32-gcc|x86_64-w64-mingw32-gcc)
                    info "MinGW cross-compiler not found."
                    info "On Debian/Ubuntu: sudo apt-get install mingw-w64"
                    info "On Fedora: sudo dnf install mingw32-gcc mingw64-gcc"
                    info "On Arch: sudo pacman -S mingw-w64-gcc"
                    info "On macOS: brew install mingw-w64"
                    ;;
                *)
                    info "Please install $tool using your distribution's package manager."
                    ;;
            esac
            exit 1
        fi
    done
}

copy_or_skip() {
    local src_dir="$1"
    local dest_dir="$2"
    
    cp "$src_dir" "$dest_dir" 2>/dev/null
    if [ $? -eq 0 ]; then
        success "File copied to $dest_dir"
    else
        warning "File not found at $src_dir, skipping..."
    fi
}

# Function to copy libraries with proper naming
copy_libs() {
    local src_dir="$1"
    local dest_dir="$2"
    local prefix="$3"
    
    section_header "Copying compiled libraries to $dest_dir folder..."
    
    create_dir_if_not_exists "$dest_dir"
    
    # Copy libraries based on platform
    case "$prefix" in
        macos-universal)
            # For macOS, we have universal libraries
            copy_or_skip "$src_dir/libmad.a" "$dest_dir/${prefix}-libmad.a"
            copy_or_skip "$src_dir/libvorbisidec.a" "$dest_dir/${prefix}-libvorbisidec.a"
            copy_or_skip "$src_dir/libogg.a" "$dest_dir/${prefix}-libogg.a"
            ;;
        linux)
            # For Linux, we have x86, x86_64 and arm64 libraries
            copy_or_skip "$src_dir/libmad-x86.a" "$dest_dir/${prefix}-x86-libmad.a"
            copy_or_skip "$src_dir/libvorbisidec-x86.a" "$dest_dir/${prefix}-x86-libvorbisidec.a"
            copy_or_skip "$src_dir/libogg-x86.a" "$dest_dir/${prefix}-x86-libogg.a"
            
            copy_or_skip "$src_dir/libmad-x86_64.a" "$dest_dir/${prefix}-x86_64-libmad.a"
            copy_or_skip "$src_dir/libvorbisidec-x86_64.a" "$dest_dir/${prefix}-x86_64-libvorbisidec.a"
            copy_or_skip "$src_dir/libogg-x86_64.a" "$dest_dir/${prefix}-x86_64-libogg.a"
            
            copy_or_skip "$src_dir/libmad-arm64.a" "$dest_dir/${prefix}-arm64-libmad.a"
            copy_or_skip "$src_dir/libvorbisidec-arm64.a" "$dest_dir/${prefix}-arm64-libvorbisidec.a"
            copy_or_skip "$src_dir/libogg-arm64.a" "$dest_dir/${prefix}-arm64-libogg.a"
            ;;
        windows)
            # For Windows, we have both win32 and win64 libraries
            copy_or_skip "$src_dir/libmad-win32.a" "$dest_dir/${prefix}-win32-libmad.a"
            copy_or_skip "$src_dir/libvorbisidec-win32.a" "$dest_dir/${prefix}-win32-libvorbisidec.a"
            copy_or_skip "$src_dir/libogg-win32.a" "$dest_dir/${prefix}-win32-libogg.a"
            
            copy_or_skip "$src_dir/libmad-win64.a" "$dest_dir/${prefix}-win64-libmad.a"
            copy_or_skip "$src_dir/libvorbisidec-win64.a" "$dest_dir/${prefix}-win64-libvorbisidec.a"
            copy_or_skip "$src_dir/libogg-win64.a" "$dest_dir/${prefix}-win64-libogg.a"
            ;;
        *)
            # Generic case for other platforms
            copy_or_skip "$src_dir/libmad.a" "$dest_dir/${prefix}-libmad.a"
            copy_or_skip "$src_dir/libvorbisidec.a" "$dest_dir/${prefix}-libvorbisidec.a"
            copy_or_skip "$src_dir/libogg.a" "$dest_dir/${prefix}-libogg.a"
            ;;
    esac

    return 0
} 