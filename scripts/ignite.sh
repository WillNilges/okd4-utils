#!/bin/bash

# Use xdotool to type in the ignition information. Screw doing that by hand!
# (Requires using X11. Sorry.)

help () {
echo 'Please provide arguments:'
echo '  $1: Node Type [ bootstrap, master, or worker ]'
echo '  $2: Node IP [x.x.x.x]'
echo
echo 'For additional configuration, edit the script.'
}

set -e

if [ "$1" == '' ] || [ "$2" == '' ]
then
    help
    exit 1
fi

if [ "$1" != 'master' ] && [ "$1" != 'bootstrap' ] && [ "$1" != 'worker' ]
then
    echo 'Node Type not recognized.'
    help
    exit 1
fi

# If you use static DHCP, fill out only these fields
NODE_TYPE=$1 # bootstrap, or master, or worker
SERVICES=10.10.51.70 # IP address of your services VM

# If you don't use static DHCP, fill these out as well
NODE_IP=$2
NAMESERVER=10.10.51.70
NODE_ROUTE=10.10.51.1
NODE_MASK=255.255.255.0

# For static DHCP
#IGNITION=" coreos.inst.install_dev=/dev/sda coreos.inst.image_url=http://$SERVICES:8080/okd4/fcos.raw.xz coreos.inst.ignition_url=http://$SERVICES:8080/okd4/$NODE_TYPE.ign"

# For non-static DHCP
IGNITION=" ip=$NODE_IP::$NODE_ROUTE:$NODE_MASK:::none nameserver=$NAMESERVER coreos.inst.install_dev=/dev/sda coreos.inst.image_url=http://$SERVICES:8080/okd4/fcos.raw.xz coreos.inst.ignition_url=http://$SERVICES:8080/okd4/$NODE_TYPE.ign"

echo 'Igniting in 3 seconds...'
sleep 3
xdotool type "$IGNITION"
