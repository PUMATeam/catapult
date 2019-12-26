# Networking
Set forwarding
```
$ sysctl -w net.ipv4.conf.all.forwarding=1
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
