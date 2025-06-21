# Docker for building SynapSeq to Windows 64-bit
FROM debian:bookworm

RUN apt-get update \
    && apt-get install -y \
    build-essential \
    curl \
    automake \
    autoconf \
    libtool \
    mingw-w64

RUN apt-get clean && rm -rf /var/lib/apt/lists/*

CMD ["tail", "-f", "/dev/null"]