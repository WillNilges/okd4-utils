#!/bin/bash

# Use ydotool to type in the ignition information. Screw doing that by hand!

NODE_TYPE=bootstrap # bootstrap, or master, or worker
NODE_TYPE=master
NODE_TYPE=worker
NODE_IP=10.10.33.76

SERVICES=129.21.49.25
NAMESERVER=10.10.33.70
NODE_ROUTE=10.10.33.1
NODE_MASK=255.255.255.0

IGNITION=" coreos.inst.install_dev=/dev/sda coreos.inst.image_url=http://$SERVICES:8080/okd4/fcos.raw.xz coreos.inst.ignition_url=http://$SERVICES:8080/okd4/$NODE_TYPE.ign"

#IGNITION=" ip=$NODE_IP::$NODE_ROUTE:$NODE_MASK:::none nameserver=$NAMESERVER coreos.inst.install_dev=/dev/sda coreos.inst.image_url=http://$SERVICES:8080/okd4/fcos.raw.xz coreos.inst.ignition_url=http://$SERVICES:8080/okd4/$NODE_TYPE.ign"

sleep 3
ydotool type "$IGNITION"