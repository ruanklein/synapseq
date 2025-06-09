#!/bin/bash

# Source common library
. ./lib.sh

# Check if 32-bit executable exists
if [ ! -f dist/synapseq-win32.exe ]; then
    error "32-bit executable not found. Run ./windows-build-synapseq.sh first."
    exit 1
fi

SETUP_NAME="synapseq-windows-setup.exe"

# Remove the existing installer if it exists
rm -f dist/${SETUP_NAME}

section_header "Creating Windows Installer..."

# Set up wine environment
export WINEARCH=win32
export WINEPREFIX=/tmp/wineprefix
export WINEDEBUG=-all
export WINEDLLOVERRIDES="winemenubuilder.exe=d"
export DISPLAY=:99

# Increase Wine memory limits
export WINE_LARGE_ADDRESS_AWARE=1
export WINE_HEAP_MAXRESERVE=4096

# Get Xvfb PID
XVFB_PID=$(pgrep -f "Xvfb $DISPLAY -screen 0 1024x768x16")

# Start Xvfb to provide a virtual display if it's not already running
if [ -z "$XVFB_PID" ]; then
    rm -f /tmp/.X${DISPLAY/:/}-lock
    info "Starting Xvfb for headless Wine operation..."
    Xvfb $DISPLAY -screen 0 1024x768x16 & XVFB_PID=$!
    sleep 2  # Wait for Xvfb to start
fi

# Initialize Wine prefix if it doesn't exist
if [ ! -d "$WINEPREFIX" ]; then
    info "Initializing Wine prefix..."
    wineboot -i > /dev/null 2>&1
    # Wait for wineboot to complete
    sleep 5
fi

# Check if Inno Setup is installed in Wine
ISCC="$WINEPREFIX/drive_c/Program Files/Inno Setup 6/ISCC.exe"
if [ ! -f "$ISCC" ]; then
    info "Inno Setup not found. Downloading and installing..."
    
    # Create temp directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    # Download Inno Setup
    curl -L -o innosetup.exe -s "https://files.jrsoftware.org/is/6/innosetup-6.4.2.exe"

    if [ $? -ne 0 ]; then
        error "Failed to download Inno Setup"
        kill $XVFB_PID
        exit 1
    fi
    
    # Install Inno Setup silently
    info "Installing Inno Setup..."
    wine innosetup.exe /VERYSILENT /SUPPRESSMSGBOXES /NORESTART /SP- /NOICONS
    
    # Wait for installation to complete
    sleep 10
    wineserver -w
    
    # Clean up
    cd - > /dev/null
    rm -rf "$TEMP_DIR"
    
    if [ ! -f "$ISCC" ]; then
        error "Failed to install Inno Setup"
        kill $XVFB_PID
        exit 1
    fi
fi

# Create the installer
info "Creating installer..."

# Kill any hanging wine processes
wineserver -k

# For convert *.md to *.txt
create_dir_if_not_exists "build"

# Convert README.md to README.txt
pandoc -f markdown -t plain README.md -o build/README.txt 

# Convert USAGE.md to USAGE.txt
pandoc -f markdown -t plain USAGE.md -o build/USAGE.txt

# Convert RESEARCH.md to RESEARCH.txt
pandoc -f markdown -t plain RESEARCH.md -o build/RESEARCH.txt

# Run ISCC with increased memory limits and in silent mode
wine "$ISCC" /O+ /Q setup.iss

# Check if the installer was created successfully
if [ ! -f "dist/${SETUP_NAME}" ]; then
    error "Failed to create installer"

    # Kill any hanging processes
    wineserver -k
    kill $XVFB_PID
    exit 1
fi

success "Installer created successfully at dist/${SETUP_NAME}"

# Final cleanup
wineserver -w
rm -rf "$WINEPREFIX" build

# Kill Xvfb
kill $XVFB_PID

section_header "Build process completed!" 