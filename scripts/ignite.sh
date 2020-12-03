#!/bin/bash

# Use ydotool to type in the ignition information. Screw doing that by hand!

NODE_TYPE=bootstrap # bootstrap, or master, or worker
#NODE_TYPE=master
#NODE_TYPE=worker
NODE_IP= # If you don't have Static DHCP

SERVICES= # This must be specified.
NAMESERVER=
NODE_ROUTE=
NODE_MASK=255.255.255.0

# For static DHCP
IGNITION=" coreos.inst.install_dev=/dev/sda coreos.inst.image_url=http://$SERVICES:8080/okd4/fcos.raw.xz coreos.inst.ignition_url=http://$SERVICES:8080/okd4/$NODE_TYPE.ign"

# For non-static DHCP
#IGNITION=" ip=$NODE_IP::$NODE_ROUTE:$NODE_MASK:::none nameserver=$NAMESERVER coreos.inst.install_dev=/dev/sda coreos.inst.image_url=http://$SERVICES:8080/okd4/fcos.raw.xz coreos.inst.ignition_url=http://$SERVICES:8080/okd4/$NODE_TYPE.ign"

sleep 3
ydotool type "$IGNITION"
