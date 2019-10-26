# Creating custom rootfs (with sshd)

Currently I have to use ubuntu, trying to use Centos proved unlucky, might try with fedora later.

```
$ truncate -s 1G rootfs

$ sudo mkfs.ext4 rootfs

$ sudo mount rootfs /mnt/rootfs

$ debootstrap --include openssh-server,vim bionic /mnt/rootfs http://archive.ubuntu.com/ubuntu/

$ sudo chroot /mnt/rootfs /bin/bash

$ passwd # set root password

$ sudo umount /mnt/rootfs
```

# Creating custom rootfs using a docker image
Using the python image for my example

Pull and create a container
```
$ docker pull python
$ docker create -t -i -name pycontainer python bash
```

Inside the container install sshd and enable systemd:
```
$ apt-get update
$ apt-get install vim openssh-server systemd # we need ssh and an init system
$ ln -s agetty /etc/init.d/agetty.ttyS0
$ echo ttyS0 > /etc/securetty
$ systemctl enable getty@ttyS0
$ passwd root # set root password

```

Create the rootfs file
```
$ docker export pycontainer -o pycon.tar
$ truncate -s 1500M rootfs
$ sudo mkfs.ext4 rootfs
$ sudo mount -o loop rootfs /mnt/
$ sudo tar -C /mnt/ -xvf pycon.tar
$ sudo umount /mnt
```
