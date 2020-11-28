#!/bin/bash
cd /root
rm -rf /var/www/html/okd4
rm -rf /root/install_dir
mkdir /root/install_dir
mkdir /var/www/html/okd4
cp /root/install-config.yaml install_dir
openshift-install create manifests --dir=install_dir/
openshift-install create ignition-configs --dir=install_dir/
sudo cp -R install_dir/* /var/www/html/okd4/
sudo chown -R apache: /var/www/html/
sudo chmod -R 755 /var/www/html/
echo $(curl localhost:8080/okd4/metadata.json)
cd /var/www/html/okd4/
wget https://builds.coreos.fedoraproject.org/prod/streams/stable/builds/32.20201004.3.0/x86_64/fedora-coreos-32.20201004.3.0-metal.x86_64.raw.xz
wget https://builds.coreos.fedoraproject.org/prod/streams/stable/builds/32.20201004.3.0/x86_64/fedora-coreos-32.20201004.3.0-metal.x86_64.raw.xz.sig
sudo mv fedora-coreos-32.20201004.3.0-metal.x86_64.raw.xz fcos.raw.xz
sudo mv fedora-coreos-32.20201004.3.0-metal.x86_64.raw.xz.sig fcos.raw.xz.sig
sudo chown -R apache: /var/www/html/
sudo chmod -R 755 /var/www/html/
