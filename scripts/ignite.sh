#!/bin/bash

# Use ydotool to type in the ignition information. Screw doing that by hand!

# If you use static DHCP, fill out only these fields
NODE_TYPE= # bootstrap, or master, or worker
SERVICES= # IP address of your services VM

# If you don't use static DHCP, fill these out as well
NODE_IP=
NAMESERVER=
NODE_ROUTE=
NODE_MASK=

# For static DHCP
IGNITION=" coreos.inst.install_dev=/dev/sda coreos.inst.image_url=http://$SERVICES:8080/okd4/fcos.raw.xz coreos.inst.ignition_url=http://$SERVICES:8080/okd4/$NODE_TYPE.ign"

# For non-static DHCP
#IGNITION=" ip=$NODE_IP::$NODE_ROUTE:$NODE_MASK:::none nameserver=$NAMESERVER coreos.inst.install_dev=/dev/sda coreos.inst.image_url=http://$SERVICES:8080/okd4/fcos.raw.xz coreos.inst.ignition_url=http://$SERVICES:8080/okd4/$NODE_TYPE.ign"

echo 'Igniting in 3 seconds...'
sleep 3
ydotool type "$IGNITION"
