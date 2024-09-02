#!/bin/bash
set -e

exec /app/remote-dev
# (
#     echo "Starting rust"
#     /app/remote-dev &
# )

# echo "Starting lldb-server"
# exec lldb-server platform --listen *:1234 --log-file /tmp/lldb-server.log --log-channels "lldb all" --server
