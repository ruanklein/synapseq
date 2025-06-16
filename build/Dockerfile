# Docker for building SynapSeq to Windows 64-bit
FROM debian:bookworm

# To work with Inno Setup
RUN dpkg --add-architecture i386 \
    && apt-get update

RUN apt-get install -y \
    build-essential \
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