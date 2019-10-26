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


# Example requests
```
curl -XPOST http://localhost:8888/hosts -d '{ "name":"hosto", "address": "192.168.122.5", "user": "root", "password": "centos"}'

curl -XPOST http://localhost:8888/vms/start -d '{ "id": "792a4940-49d1-4255-b31c-ed4a169dcc1c", "name": "hello", "status": 0, "host_id": "792a4940-49d1-4255-b31c-ed4a169dcc1c", "vcpu": 1, "memory":128}'
```


# Networking
Set
```
sysctl -w net.ipv4.conf.all.forwarding=1
```
Create a bridge
```
$ brctl addbr br0
$ ip addr flush dev eth0 # flush ip
$ brctl addif br0 eth0 # add eth0 to the bridge
$ dhclient -v br0
$ ip tuntap add dev fc-tap-$UUID mode tap
$ brctl addif br0 fc-tap-$UUID
$ ip link set fc-tap-$UUID up
 ```

# Stuff to think about
- How to load existing VMs when catapult-node starts?
  - Use an sqlite3 db and consider it the only source of truth
  - Scan the /var/vms
  - Combine both options above - but adding a vms from /var/vms to the db
might be challenging as it would require learning about the VM specifications
from the API (there's no "virsh dumpxml") as far as I know

# Handling host state
- Upon starting we should check all the present hosts state
  - If a host is UP, we need to try and health check it using grpc
    - if it failes to respond it might mean the node process is down (should we turn it into a systemd service?), move it into DOWN state
  - We should not accept any request until we determined the host states

- We should be able to de/activate a host - an ansible script to start/stop
  the catapult-node process
