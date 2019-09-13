#!/bin/python
import libvirt
import sys

from xml.etree import ElementTree as ET

def should_update(section, mac, ip):
    for node in section.findall(".//host"):
        if node.attrib["mac"] == mac and node.attrib["ip"] == ip:
            return False

    return True


def find_dhcp(root): 
    dhcp_section = root.find("./ip/dhcp")
    return dhcp_section


def update_network(dom_name, root, mac, ip, dhcp_section):
    host = ET.Element("host", mac=mac, name=dom_name, ip=ip)
    dhcp_section.append(host)

    return ET.tostring(root).decode()

# connect to libvirt
conn = libvirt.open("qemu:///system")

# lookup domain passed to us
dom_name = sys.argv[1]
try:
    dom = conn.lookupByName(dom_name)
except:
    sys.exit(-1)

ip_address = sys.argv[3]
mac_address = sys.argv[4]

network = conn.networkLookupByName(sys.argv[2])
network_xml = network.XMLDesc(0)
xml_root = ET.fromstring(network_xml)
dhcp_section = find_dhcp(xml_root)

if not should_update(dhcp_section, mac_address, ip_address):
    print("ip already configured")
    sys.exit(0)

updated_network = update_network(
        dom_name,
        xml_root,
        mac_address,
        ip_address,
        dhcp_section)

print(updated_network)
if network.isActive():
    network.destroy()

network = conn.networkDefineXML(updated_network)
network.setAutostart(True)
network.create()

