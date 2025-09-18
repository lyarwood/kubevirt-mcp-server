# kubevirt-mcp-server

A simple Model Context Protocol server for KubeVirt.

## Architecture

The project is organized into modular packages:

- `main.go` - MCP server setup and registration
- `pkg/client/` - Shared KubeVirt client utilities
- `pkg/tools/` - MCP tool handlers for VM operations
- `pkg/resources/` - MCP resource handlers for structured data access
- `scripts/kubevirtci.sh` - Script for managing local kubevirtci development environment
- `scripts/sync.sh` - Script for building and running MCP server locally with kubevirtci access
- `Makefile` - Build automation and development tasks

## Features

### MCP Tools
- `list_vms` - List virtual machine names in a namespace
- `start_vm` - Start a virtual machine
- `stop_vm` - Stop a virtual machine
- `restart_vm` - Restart a virtual machine (handles both running and stopped VMs)
- `pause_vm` - Pause a virtual machine
- `unpause_vm` - Unpause a virtual machine
- `create_vm` - Create a virtual machine with specified container disk (supports OS name lookup), optional instancetype and preference
- `delete_vm` - Delete a virtual machine
- `patch_vm` - Apply JSON merge patch to modify VM configuration
- `list_instancetypes` - List available instance types
- `get_instancetype` - Get detailed information about a specific instance type
- `get_preference` - Get detailed information about a specific preference
- `get_vm_instancetype` - Get instance type for a VM
- `get_vm_status` - Get comprehensive VM status information
- `get_vm_conditions` - Get detailed VM condition information
- `get_vm_phase` - Get current VM phase and basic status
- `get_vm_disks` - Retrieve the list of disks attached to a virtual machine

### MCP Prompts
- `describe_vm` - Provide comprehensive VM description including configuration, status, and operational details
- `troubleshoot_vm` - Diagnose and analyze potential VM issues with actionable recommendations
- `health_check_vm` - Perform quick VM health check and report issues

### MCP Resources
- `kubevirt://{namespace}/vms` - JSON list of VMs with summary info
- `kubevirt://{namespace}/vm/{name}` - Complete VM specification
- `kubevirt://{namespace}/vm/{name}/status` - VM status and phase information
- `kubevirt://{namespace}/vm/{name}/console` - VM console connection details
- `kubevirt://{namespace}/vmis` - JSON list of VMIs with runtime info
- `kubevirt://{namespace}/vmi/{name}` - Complete VMI specification
- `kubevirt://{namespace}/vmi/{name}/guestosinfo` - VMI guest OS information
- `kubevirt://{namespace}/vmi/{name}/filesystems` - VMI filesystem information
- `kubevirt://{namespace}/vmi/{name}/userlist` - VMI user list information
- `kubevirt://{namespace}/datavolumes` - JSON list of DataVolumes with source and storage info
- `kubevirt://{namespace}/datavolume/{name}` - Complete DataVolume specification
- `kubevirt://{namespace}/instancetypes` - Namespaced instance types
- `kubevirt://{namespace}/preferences` - Namespaced VM preferences
- `kubevirt://cluster/instancetypes` - Cluster-wide instance types
- `kubevirt://cluster/preferences` - Cluster-wide VM preferences
- `kubevirt://cluster/instancetype/{name}` - Specific cluster instance type
- `kubevirt://cluster/preference/{name}` - Specific cluster preference

## Building

```bash
# Using Makefile (recommended)
make build

# Or directly with go
go build -o kubevirt-mcp-server .
```

## Development

### Available Make Targets
- `make build` - Build the binary (default)
- `make clean` - Clean build artifacts  
- `make test` - Run tests with Ginkgo framework
- `make coverage` - Generate test coverage report
- `make fmt` - Format Go code
- `make vet` - Run go vet
- `make lint` - Run golangci-lint
- `make deps` - Download and tidy dependencies
- `make run` - Build and run the server
- `make check` - Run fmt, vet, lint, and test
- `make cluster-up` - Start kubevirtci cluster for testing
- `make cluster-down` - Stop kubevirtci cluster
- `make cluster-sync` - Build and run MCP server locally with kubevirtci access
- `make test-functional` - Run functional tests against MCP server
- `make help` - Show help message

### Testing

The project uses the [Ginkgo](https://github.com/onsi/ginkgo) testing framework with [Gomega](https://github.com/onsi/gomega) assertions:

```bash
# Run unit tests
make test

# Generate test coverage report
make coverage

# Run functional tests against MCP server
make test-functional

# Run linter
make lint

# Run all quality checks
make check
```

### Local Development Environment

For functional testing with a real KubeVirt cluster, use the kubevirtci integration:

```bash
# Start a local kubevirtci cluster with KubeVirt
make cluster-up

# Stop the cluster when done
make cluster-down

# Build and run MCP server locally with cluster access
make cluster-sync
```

The kubevirtci integration includes:

**`scripts/kubevirtci.sh`** handles:
- Downloading and setting up kubevirtci
- Starting a local Kubernetes cluster with KubeVirt and CDI
- Configuring the environment for testing
- Providing access to kubectl, kubeconfig, and registry

**`scripts/sync.sh`** handles:
- Building the MCP server binary
- Setting up proper KUBECONFIG environment for kubevirtci access
- Providing instructions for local MCP server execution

### Test Structure

The project includes comprehensive test coverage:

- **Unit Tests** - Test individual components in isolation
  - `pkg/client/client_test.go` - KubeVirt client creation tests
  - `pkg/tools/tools_test.go` - MCP tool handler argument validation
  - `pkg/resources/resources_test.go` - MCP resource handler URI parsing

- **Functional Tests** - Test complete MCP server functionality
  - `tests/functional/functional_suite_test.go` - Test suite setup and KubeVirt cluster verification
  - `tests/functional/mcp_server_stdio_test.go` - Complete MCP server API coverage:
    - MCP server initialization and JSON-RPC communication
    - All MCP tools: list_vms, start_vm, stop_vm, restart_vm, create_vm, list_instancetypes, get_vm_instancetype, get_vm_disks
    - All MCP resources: kubevirt://namespace/vms, vm/name, vmis, vmi/name endpoints
    - Error handling for invalid tools, missing arguments, invalid URIs, and non-existent VMs

## Using with Claude CLI

This MCP server integrates seamlessly with the Claude CLI (Claude Code) to provide KubeVirt management capabilities directly within your development workflow.

### Configuration

#### Method 1: Project-Specific Configuration

Create a `.clauderc` file in your project directory:

```json
{
  "mcp": {
    "servers": {
      "kubevirt": {
        "command": "/path/to/kubevirt-mcp-server",
        "env": {
          "KUBECONFIG": "/path/to/your/kubeconfig"
        }
      }
    }
  }
}
```

#### Method 2: Global Configuration

Configure in your global Claude settings:

```bash
# Create/edit global Claude config
mkdir -p ~/.config/claude
cat > ~/.config/claude/config.json << 'EOF'
{
  "mcp": {
    "servers": {
      "kubevirt": {
        "command": "/path/to/kubevirt-mcp-server",
        "env": {
          "KUBECONFIG": "/path/to/your/kubeconfig"
        }
      }
    }
  }
}
EOF
```

#### Method 3: Environment Variables

Set up environment variables for dynamic configuration:

```bash
export KUBEVIRT_MCP_SERVER="/path/to/kubevirt-mcp-server"
export KUBECONFIG="/path/to/your/kubeconfig"

# Use in Claude config
{
  "mcp": {
    "servers": {
      "kubevirt": {
        "command": "$KUBEVIRT_MCP_SERVER",
        "env": {
          "KUBECONFIG": "$KUBECONFIG"
        }
      }
    }
  }
}
```

### Setup Steps

1. **Build the MCP Server**
   ```bash
   git clone <your-repo>
   cd kubevirt-mcp-server
   make build
   
   # Note the path to the binary
   echo "Binary location: $(pwd)/kubevirt-mcp-server"
   ```

2. **Verify KubeVirt Access**
   ```bash
   # Test your kubeconfig works
   kubectl get vms --all-namespaces
   
   # Test the MCP server directly
   export KUBECONFIG=/path/to/your/kubeconfig
   echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | ./kubevirt-mcp-server
   ```

3. **Configure Claude CLI**
   
   Use one of the configuration methods above, ensuring:
   - The `command` path points to your built binary
   - The `KUBECONFIG` environment variable points to your cluster config
   - The kubeconfig has appropriate permissions for VM operations

4. **Test the Integration**
   ```bash
   # Start Claude CLI and test MCP server connectivity
   claude --list-mcp-servers
   
   # Should show "kubevirt" server as available
   ```

### Usage Examples

Once configured, you can use Claude CLI with natural language to manage your VMs:

#### VM Management
```bash
# List and manage VMs
claude "List all VMs in the production namespace"
claude "Start the web-server VM in default namespace"
claude "Restart all stopped VMs in the staging namespace"
claude "Show me the configuration of the database VM"
```

#### Troubleshooting
```bash
# VM diagnostics
claude "The payment-service VM isn't responding, can you investigate?"
claude "Compare the instance types of VMs in prod vs staging"
claude "What VMs are currently running and what resources are they using?"
```

#### Bulk Operations
```bash
# Mass management
claude "Stop all VMs in the test namespace"
claude "List all VMs that don't have instance types assigned"
claude "Show me a summary of VM status across all namespaces"
```

#### Development Workflow
```bash
# Development tasks
claude "Start my development VMs (web-dev, db-dev, cache-dev)"
claude "Check if my feature branch VMs are ready for testing"
claude "Clean up any VMs from old feature branches"
```

### Available Capabilities

The MCP server provides Claude with these tools:

**VM Lifecycle:**
- `list_vms` - List VMs in a namespace
- `start_vm` - Start a specific VM
- `stop_vm` - Stop a specific VM
- `restart_vm` - Restart a VM
- `pause_vm` - Pause a VM
- `unpause_vm` - Unpause a VM
- `create_vm` - Create a new VM with container disk (supports OS names like "fedora", "ubuntu") and optional instancetype/preference
- `delete_vm` - Delete a VM
- `patch_vm` - Apply JSON merge patch to modify VM configuration

**VM Information:**
- `get_vm_status` - Get comprehensive VM status information
- `get_vm_conditions` - Get detailed VM condition information
- `get_vm_phase` - Get current VM phase and basic status
- `get_vm_instancetype` - Get VM's assigned instance type
- `get_vm_disks` - Retrieve the list of disks attached to a virtual machine

**Instance Types & Preferences:**
- `list_instancetypes` - List available instance types
- `get_instancetype` - Get detailed information about a specific instance type
- `get_preference` - Get detailed information about a specific preference

**Prompts:**
- `describe_vm` - Comprehensive VM description and analysis
- `troubleshoot_vm` - VM troubleshooting and diagnostics
- `health_check_vm` - Quick VM health checks

**Container Disk Lookup:**
The `create_vm` tool supports both full container image URLs and OS name shortcuts:
- **OS Names**: `"fedora"`, `"ubuntu"`, `"centos"`, `"debian"`, `"rhel"`, `"opensuse"`, `"alpine"`, `"cirros"`, `"windows"`, `"freebsd"`
- **Full URLs**: `"quay.io/containerdisks/fedora:latest"`, `"my-registry/my-image:tag"`
- **Auto-resolve**: Unknown OS names are resolved to `quay.io/containerdisks/{name}:latest`

**Structured Data:**
- `kubevirt://{namespace}/vms` - VM summary data
- `kubevirt://{namespace}/vm/{name}` - Complete VM specs
- `kubevirt://{namespace}/vmis` - VM instance runtime data
- `kubevirt://{namespace}/vmi/{name}` - Complete VMI specs

### Security Considerations

- **Permissions**: The MCP server uses your KUBECONFIG credentials
- **Scope**: Claude has the same KubeVirt permissions as your kubeconfig
- **Best Practice**: Consider using a dedicated service account with limited permissions:

```bash
# Create a dedicated service account for MCP operations
kubectl create serviceaccount kubevirt-mcp-user
kubectl create clusterrole kubevirt-mcp-role \
  --verb=get,list,create,update,patch,delete \
  --resource=virtualmachines,virtualmachineinstances
kubectl create clusterrolebinding kubevirt-mcp-binding \
  --clusterrole=kubevirt-mcp-role \
  --serviceaccount=default:kubevirt-mcp-user

# Generate kubeconfig for the service account
kubectl create token kubevirt-mcp-user > /path/to/mcp-kubeconfig
```

### Troubleshooting

**Connection Issues:**
```bash
# Test MCP server manually
export KUBECONFIG=/path/to/your/kubeconfig
echo '{"jsonrpc":"2.0","method":"initialize","params":{"capabilities":{},"clientInfo":{"name":"test","version":"1.0"}},"id":1}' | ./kubevirt-mcp-server
```

**Permission Issues:**
```bash
# Verify kubeconfig access
kubectl auth can-i get virtualmachines
kubectl auth can-i create virtualmachines
```

**Claude CLI Debug:**
```bash
# Enable verbose logging
claude --verbose "List my VMs"

# Check MCP server logs
CLAUDE_MCP_DEBUG=1 claude "List VMs in default namespace"
```

### Advanced Usage

**Project-Specific VM Management:**
```bash
# Create a .clauderc for your project
cat > .clauderc << 'EOF'
{
  "mcp": {
    "servers": {
      "kubevirt": {
        "command": "/usr/local/bin/kubevirt-mcp-server",
        "env": {
          "KUBECONFIG": "./k8s/kubeconfig",
          "DEFAULT_NAMESPACE": "myproject-dev"
        }
      }
    }
  },
  "context": {
    "project": "MyProject Development VMs",
    "defaultNamespace": "myproject-dev"
  }
}
EOF

# Now Claude understands your project context
claude "Start my development environment"
claude "Show me the status of project VMs"
```

## Demo

This short demo uses mcp-cli as a bridge between the kubevirt-mcp-server and LLM.

The model used by the demo is llama3.2 running locally under ollama.

![demo](demo.gif)

## Links 

- https://www.anthropic.com/news/model-context-protocol
- https://github.com/mark3labs/mcp-go
- https://github.com/chrishayuk/mcp-cli
