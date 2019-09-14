#!/bin/bash
echo "Defining VM..."

VM_NAME=${1:-"fc_host"}
LIBVIRT_NETWORK=${2:-"default"}
VM_IP=${3:-"192.168.122.45"}

if virsh list --all | grep -q "${VM_NAME}"; then
    echo "${VM_NAME} is already installed... "
else
    dom=$(virt-install --import --name "${VM_NAME}" \
        --memory 1024 --vcpus 1 --cpu host \
        --disk os.img,bus=virtio \
        --os-type=linux \
        --graphics spice \
        --noautoconsole \
        --network=default,model=virtio \
        --connect qemu:///system \
        --print-xml)
    echo $dom | virsh define /dev/stdin
fi

fc_host_status=$(virsh list | grep fc_host | tr -s \"[:blank:]\" | cut -d ' ' -f4)
if [  "${fc_host_status}" == 'running' ]; then
    echo "${VM_NAME} is already running"
    exit 0
fi

MAC_ADDRESS=$(virsh dumpxml "${VM_NAME}" | grep "mac address" | awk -F\' '{ print $2}')
echo "Setting IP address to ${VM_IP} for MAC address ${MAC_ADDRESS}"

xml_entry="<host mac=\"${MAC_ADDRESS}\" name=\"${VM_NAME}\" ip=\"${VM_IP}\"/>"
if virsh net-dumpxml "${LIBVIRT_NETWORK}" | grep -q "${VM_NAME}"; then
    echo "IP address is already configured"
else
    virsh net-update default add-last ip-dhcp-host "${xml_entry}" --live --config
fi

echo "starting ${VM_NAME}..."
virsh start "${VM_NAME}"
