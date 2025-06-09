# Docker for building SynapSeq to Linux and Windows
FROM debian:bookworm

RUN apt-get update && apt-get install -y \
    build-essential \
    libasound2-dev \
    libtool \
    automake \
    autoconf \
    curl

RUN if [ "$(uname -m)" = "x86_64" ]; then \
        # Compile for 32-bit and for Windows 32-bit/64-bit
        dpkg --add-architecture i386 && \
        apt-get update && \
        apt-get install -y \
            libasound2-dev:i386 \
            gcc-multilib \
            mingw-w64 \
            wine \
            xvfb \
            pandoc; \
    fi

RUN apt-get clean && rm -rf /var/lib/apt/lists/*

CMD ["tail", "-f", "/dev/null"]