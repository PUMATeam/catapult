#!/bin/bash

IP=$(hostname -I)
IP_WITH_CIDR=$(ip addr show dev eth0 | grep inet | awk '{print $2}')
ROUTE=$(ip route list dev eth0 | grep -v default | cut -f 1 -d ' ')
DEFAULT_ROUTE=$(ip route list dev eth0 | grep default | awk '{print $3}')

ip link add fcbridge type bridge
ip link set fcbridge up
ip addr flush dev eth0
ip link set eth0 master fcbridge
ip addr add "${IP_WITH_CIDR}" brd + dev fcbridge
ip route add default via "${DEFAULT_ROUTE}"
ip route add "${ROUTE}"  dev fcbridge proto kernel scope link src "${IP}"
