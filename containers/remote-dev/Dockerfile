FROM rust:1.80 AS builder
WORKDIR /app

RUN apt-get update && \
    apt-get -y upgrade && \
    apt-get install -y \
    git \
    libclang-dev \
    pkg-config \
    curl \
    build-essential

COPY . .

RUN cargo build

# Use Ubuntu as the release image
FROM mcr.microsoft.com/devcontainers/rust:latest AS runtime
WORKDIR /app

RUN apt-get update && \
    apt-get install -y \
    curl \
    jq \
    iproute2 \
    git \
    lldb

COPY --from=builder /app/target/debug/remote-dev .
COPY entrypoint.sh .

ENTRYPOINT ["/app/entrypoint.sh"]
