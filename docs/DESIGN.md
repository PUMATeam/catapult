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

# Creating custom rootfs (with sshd)

Currently I have to use ubuntu, trying to use Centos proved unlucky, might try with fedora later.

```
$ truncate -s 1G rootfs 

$ sudo mount rootfs /mnt/rootfs

$ debootstrap --include openssh-server,vim bionic /mnt/rootfs http://archive.ubuntu.com/ubuntu/

$ sudo chroot /mnt/rootfs /bin/bash

$ passwd # set root password

$ sudo umount /mnt/rootfs
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




