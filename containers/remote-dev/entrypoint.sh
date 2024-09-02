#!/bin/bash
set -e

/app/remote-dev &
PID=$!

exec lldb-server gdbserver 0.0.0.0:1234 --attach $PID

# (
#     echo "Starting rust"
#     /app/remote-dev &
# )

# echo "Starting lldb-server"
# exec lldb-server platform --listen *:1234 --log-file /tmp/lldb-server.log --log-channels "lldb all" --server
