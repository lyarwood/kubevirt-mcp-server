#!/bin/bash
#
# Copyright 2023 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex

# Configuration
SERVER_BINARY="kubevirt-mcp-server"
MCP_PID_FILE="/tmp/kubevirt-mcp-server.pid"

# Get base directory (project root)
BASEDIR=${BASEDIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}
KUBECTL=${KUBECTL:-kubectl}

# Get kubeconfig from kubevirtci
KUBECONFIG=${KUBECONFIG:-$("${BASEDIR}/scripts/kubevirtci.sh" kubeconfig)}

if [ ! -f "$KUBECONFIG" ]; then
    echo "Error: KUBECONFIG not found at $KUBECONFIG"
    echo "Please run 'make cluster-up' first to start the kubevirtci environment"
    exit 1
fi

echo "Using KUBECONFIG: $KUBECONFIG"

# Build the MCP server binary
echo "Building MCP server..."
cd "$BASEDIR"
make build

# Stop any existing MCP server process
echo "Stopping any existing MCP server..."
if [ -f "$MCP_PID_FILE" ]; then
    OLD_PID=$(cat "$MCP_PID_FILE")
    if kill -0 "$OLD_PID" 2>/dev/null; then
        echo "Stopping existing MCP server process (PID: $OLD_PID)"
        kill "$OLD_PID"
        # Wait for process to stop
        for i in {1..10}; do
            if ! kill -0 "$OLD_PID" 2>/dev/null; then
                break
            fi
            sleep 1
        done
        # Force kill if still running
        if kill -0 "$OLD_PID" 2>/dev/null; then
            echo "Force killing MCP server process"
            kill -9 "$OLD_PID" 2>/dev/null || true
        fi
    fi
    rm -f "$MCP_PID_FILE"
fi

# Explain MCP server execution model
echo "✅ MCP server binary built successfully!"
echo "📍 Binary location: $BASEDIR/$SERVER_BINARY"
echo "🔧 KUBECONFIG: $KUBECONFIG"
echo ""
echo "🏃 MCP servers use stdio communication and cannot run as background daemons."
echo "📡 To run the MCP server, execute it directly:"
echo ""
echo "  cd $BASEDIR"
echo "  export KUBECONFIG=$KUBECONFIG"
echo "  ./$SERVER_BINARY"
echo ""
echo "💡 The server will wait for JSON-RPC messages on stdin and respond on stdout."
echo "🧪 Use 'make test-functional' to test the server with automated stdio communication."