# Docker for building SynapSeq to Linux and Windows
FROM debian:bookworm

# To Inno Setup
RUN dpkg --add-architecture i386 \
    && apt-get update

RUN apt-get install -y \
    build-essential \
    pkg-config \
    libasound2-dev \
    libmad0-dev \
    libvorbis-dev \
    libogg-dev \
    # For building Windows .exe and installer
    curl \
    automake \
    autoconf \
    libtool \
    mingw-w64 \
    wine \
    xvfb \
    pandoc

RUN apt-get clean && rm -rf /var/lib/apt/lists/*

CMD ["tail", "-f", "/dev/null"]