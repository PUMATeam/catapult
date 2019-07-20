- Management Server
- Nodes
    - Server - communicates with the management server using gRPC
        API:
            - Start VM (cpu, memory, network)
            - Stop VM

- Upon adding a node firecracker and additional tools (gRPC client/server) will be installed

- VM
    - will be assigned a uuid
    - the uuid will be used for the unix socket path as well
    - create a tap device for the VM, same uuid will be used
