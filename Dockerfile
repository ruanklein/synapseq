# Docker for building SynapSeq to Linux and Windows
FROM debian:bookworm

RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    libasound2-dev \
    libmad0-dev \
    libvorbis-dev \
    libogg-dev \
    curl \
    automake \
    autoconf \
    libtool \
    # For Windows
    mingw-w64 \
    wine \
    xvfb \
    pandoc

RUN apt-get clean && rm -rf /var/lib/apt/lists/*

CMD ["tail", "-f", "/dev/null"]