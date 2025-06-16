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
