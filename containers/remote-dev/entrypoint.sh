#!/bin/bash
set -eu

(
    echo "Starting rust"
    /app/remote-dev &
)
RUST_PID=$!
echo "Rust program started with PID: $RUST_PID"

echo "Starting lldb-server"
exec lldb-server platform --listen *:1234 --log-file /tmp/lldb-server.log --log-channels "lldb all" --server
