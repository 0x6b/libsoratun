FROM golang:1.24.2-bullseye
RUN apt-get update && apt-get install -y make
RUN apt install -y gcc-multilib gcc-mingw-w64 binutils-mingw-w64
RUN rm -rf /var/lib/apt/lists/*