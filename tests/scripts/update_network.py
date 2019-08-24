#!/bin/python
import libvirt
import sys

from xml.etree import ElementTree as ET

def update_network(dom_name, xml, mac, ip):
    root = ET.fromstring(network_xml)
    dhcp_section = root.find('./ip/dhcp')
    host = ET.Element("host", mac=mac, name=dom_name, ip=ip)
    dhcp_section.append(host)

    return ET.tostring(root).decode()

# connect to libvirt
conn = libvirt.open('qemu:///system')

# lookup domain passed to us
dom_name = sys.argv[1]
try:
    dom = conn.lookupByName(dom_name)
except:
    sys.exit(-1)

mac_address = sys.argv[4]

network = conn.networkLookupByName(sys.argv[2])
network_xml = network.XMLDesc(0)
updated_network = update_network(dom_name, network_xml, mac_address, sys.argv[3])

if network.isActive():
    network.destroy()

network = conn.networkDefineXML(updated_network)
network.setAutostart(True)
network.create()

