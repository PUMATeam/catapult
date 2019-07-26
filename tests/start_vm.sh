#!/bin/bash 
echo "Defining centos VM..."
if virsh list --all | grep -q centos7; then
    echo "centos7 is already installed... "
else
    dom=$(virt-install --import --name centos7 \
        --memory 1024 --vcpus 1 --cpu host \
        --disk /home/libvirt/images/CentOS-7-x86_64-GenericCloud.qcow2,bus=virtio \
        --os-type=linux \
        --os-variant=centos7.0 \
        --graphics spice \
        --noautoconsole \
        --connect qemu:///system \
        --print-xml)
    echo $dom | virsh define /dev/stdin
fi

#ip_address=192.168.122.45
#mac_address=$(virsh dumpxml centos7 | grep "mac address" | awk -F\' '{ print $2}')
#echo "Setting ip address to ${ip_address} for MAC address ${mac_address}"
#./update_network.py centos7 default ${ip_address} $mac_address
#virsh start centos7
