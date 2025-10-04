# TODO - KubeVirt MCP Server Improvements

This document tracks potential improvements and future features for the KubeVirt MCP server project.

## High Priority

### MCP Tools & Resources (Based on KubeVirt API Analysis)

#### New MCP Tools
- [x] `delete_vm` - Delete a virtual machine
- [ ] `update_vm` - Update VM configuration (memory, CPU, etc.)
- [x] `pause_vm` - Pause a running VM
- [x] `unpause_vm` - Unpause a paused VM
- [ ] `addvolume_vm` - Add volume to VM
- [ ] `removevolume_vm` - Remove volume from VM
- [ ] `migrate_vm` - Migrate VM to different node
- [ ] `get_vm_console` - Get VM console connection info
- [ ] `get_vm_guestinfo` - Get guest OS information
- [ ] `get_vm_filesystems` - List VM filesystems
- [ ] `get_vm_userlist` - Get VM user list
- [ ] `clone_vm` - Clone VM from existing VM or template

#### New MCP Resources
- [x] `kubevirt://{namespace}/datavolumes` - List DataVolumes
- [x] `kubevirt://{namespace}/datavolume/{name}` - Get specific DataVolume
- [x] `kubevirt://{namespace}/vm/{name}/status` - Get VM status and phase
- [x] `kubevirt://{namespace}/vmi/{name}/guestosinfo` - Get guest OS info
- [x] `kubevirt://{namespace}/vmi/{name}/filesystems` - Get filesystem info
- [x] `kubevirt://{namespace}/vmi/{name}/userlist` - Get user list
- [x] `kubevirt://{namespace}/vm/{name}/console` - Get console connection details
- [x] `kubevirt://{namespace}/instancetypes` - List namespaced instance types
- [x] `kubevirt://{namespace}/preferences` - List namespaced VM preferences
- [x] `kubevirt://cluster/instancetypes` - List cluster-wide instance types
- [x] `kubevirt://cluster/preferences` - List cluster-wide VM preferences
- [x] `kubevirt://cluster/instancetype/{name}` - Get specific cluster instance type
- [x] `kubevirt://cluster/preference/{name}` - Get specific cluster preference

### Core Functionality
- [x] Add VM restart functionality to MCP tools
- [ ] Implement VM instance type modification tool
- [ ] Add VM status filtering in resource handlers (running, stopped, etc.)
- [ ] Add label-based filtering for VM and VMI resources
- [ ] Implement VM cloning/template functionality
- [x] Add VM deletion tool (delete_vm)
- [ ] Add VM update/modification tool (update_vm)
- [x] Implement VM pause/unpause functionality
- [ ] Add VM addvolume/removevolume tools for disk management

### Error Handling & Validation
- [ ] Improve error messages with more specific context
- [ ] Add input validation for VM names (Kubernetes naming rules)
- [ ] Add namespace existence validation before operations
- [ ] Implement retry logic for transient Kubernetes API errors

### Security & Authentication
- [ ] Add proper RBAC validation for operations
- [ ] Implement service account token handling
- [ ] Add namespace access control validation
- [ ] Consider adding audit logging for operations

## Medium Priority

### Resource Management
- [ ] Add VM disk management tools (attach/detach volumes)
- [ ] Implement VM network interface management
- [ ] Add VM snapshot creation and management
- [ ] Support for VM migration between nodes
- [ ] Add VM resource usage metrics (CPU, memory, storage)
- [ ] Add VMI (VirtualMachineInstance) specific tools and resources
- [ ] Implement Data Volume (DV) management tools
- [ ] Add Virtual Machine Instance Preset (VMIP) support
- [ ] Support for VM expanding/shrinking (hot-plug resources)

### User Experience
- [ ] Add VM console access capabilities (VNC, serial, guest agent)
- [ ] Implement VM event streaming/monitoring
- [ ] Add bulk operations (start/stop multiple VMs)
- [ ] Provide VM scheduling and lifecycle management
- [ ] Add VM backup and restore functionality
- [ ] Add VM guest OS info retrieval tool
- [ ] Implement VM filesystem listing tool
- [ ] Add VM userlist and guest agent info tools
- [ ] Support for VM screenshot/image capture
- [ ] Add VM performance metrics and monitoring tools

### Performance & Scalability
- [ ] Implement caching for frequently accessed resources
- [ ] Add pagination support for large VM lists
- [ ] Optimize Kubernetes API calls with field selectors
- [ ] Add connection pooling for KubeVirt client

## Low Priority

### Development & Tooling
- [ ] Add integration tests with real KubeVirt cluster
- [ ] Implement mock KubeVirt client for testing
- [ ] Add benchmarking tests for performance monitoring
- [ ] Create Docker/container image for deployment
- [ ] Add Helm chart for Kubernetes deployment

### Documentation & Examples
- [ ] Create detailed API documentation
- [ ] Add usage examples for each MCP tool and resource
- [ ] Create troubleshooting guide
- [ ] Add architecture documentation with diagrams
- [ ] Create deployment guide for different environments

### Monitoring & Observability
- [ ] Add structured logging with levels
- [ ] Implement metrics collection (Prometheus compatible)
- [ ] Add health check endpoints
- [ ] Create dashboard for monitoring server status
- [ ] Add distributed tracing support

## Future Considerations

### Advanced Features
- [ ] Support for VM GPU passthrough management
- [ ] Implement VM auto-scaling based on metrics
- [ ] Add multi-cluster support
- [ ] Integration with CI/CD pipelines for VM management
- [ ] Support for VM templating and golden images
- [ ] Add VM clone capabilities with different storage classes
- [ ] Implement VM live migration tools
- [ ] Support for VM memory dump and analysis
- [ ] Add VM VSOCK (virtio socket) management
- [ ] Implement VM SELinux and security context management

### Integration
- [ ] Add support for other virtualization platforms
- [ ] Integrate with cloud provider VM services
- [ ] Add support for bare metal provisioning
- [ ] Integration with configuration management tools

## Technical Debt

### Code Quality
- [ ] Add more comprehensive test coverage (target >80%)
- [ ] Implement proper dependency injection
- [ ] Add interface abstractions for better testability
- [ ] Refactor large functions into smaller, focused ones
- [ ] Add comprehensive API documentation with examples

### Architecture
- [ ] Consider implementing plugin architecture for extensibility
- [ ] Add configuration file support for server settings
- [ ] Implement graceful shutdown handling
- [ ] Add request/response middleware support
- [ ] Consider adding API versioning strategy

---

## Contributing

When working on items from this TODO list:

1. Create a feature branch following the naming convention: `feat/todo-item-name`
2. Update relevant tests and documentation
3. Follow the Conventional Commits standard for commit messages
4. Update both AGENTS.md and README.md if the changes affect project structure or usage
5. Include `Assisted-By: Claude <noreply@anthropic.com>` in commits when AI assistance is used

## Priority Guidelines

- **High Priority**: Core functionality improvements that directly impact user experience
- **Medium Priority**: Feature enhancements that add significant value
- **Low Priority**: Quality of life improvements and developer experience
- **Future Considerations**: Ideas for major features or architectural changes
- **Technical Debt**: Code quality and maintainability improvements

Items can be promoted between priority levels based on user feedback and project needs.