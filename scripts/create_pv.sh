#!/bin/bash

# Args: $1 - PV name
#       $2 - Storage qty
#       $3 - NFS server path

# The NFS Share path is hardcoded, so be sure that's where you want it!

mkdir /var/nfsshare/$1
chmod 777 /var/nfsshare/$1

cat << EOT > /tmp/more_storage.yaml
apiVersion: v1
kind: PersistentVolume
metadata: 
  name: $1
spec:
  capacity:
    storage: $2
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  nfs:
    path: /var/nfsshare/$1
    server: $3     
EOT

cat /tmp/more_storage.yaml

oc create -f /tmp/more_storage.yaml

oc get pv

