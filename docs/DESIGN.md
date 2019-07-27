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
    - need to save VMs on the host to keep track (sqlite?)


curl -XPOST http://localhost:8888/hosts -d '{ "name":"hosto", "address": "192.168.122.5", "user": "root", "password": "centos"}'

curl -XPOST http://localhost:8888/vms/start -d '{ "id": "792a4940-49d1-4255-b31c-ed4a169dcc1c", "name": "hello", "status": 0, "host_id": "792a4940-49d1-4255-b31c-ed4a169dcc1c", "vcpu": 1, "memory":128}'