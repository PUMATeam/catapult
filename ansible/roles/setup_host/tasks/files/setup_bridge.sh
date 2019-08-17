#!/bin/bash

brctl addbr fcbridge
ip addr flush dev eth0
brctl addif fcbridge eth0
dhclient -v fcbridge
ip link set fcbridge up
