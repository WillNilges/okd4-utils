#!/bin/bash

# Configures an NFS share for the OKD registry on your okd4-services node.

sudo dnf install -y nfs-utils
sudo systemctl enable nfs-server rpcbind
sudo systemctl start nfs-server rpcbind
sudo mkdir -p /var/nfsshare/registry
sudo chmod -R 777 /var/nfsshare
sudo chown -R nobody:nobody /var/nfsshare

SUBNET=10.10.51.0

echo '/var/nfsshare ' $SUBNET '/24(rw,sync,no_root_squash,no_all_squash,no_wdelay)' | sudo tee /etc/exports

sudo setsebool -P nfs_export_all_rw 1
sudo systemctl restart nfs-server
sudo firewall-cmd --permanent  --add-service mountd
sudo firewall-cmd --permanent  --add-service rpc-bind
sudo firewall-cmd --permanent  --add-service nfs
#sudo firewall-cmd --permanent --port
sudo firewall-cmd --reload
