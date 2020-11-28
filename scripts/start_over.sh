#!/bin/bash
echo 'ARE YOU TOTALLY SURE YOU WANT TO START OVER???? (Press enter to continue, ^C to abort)'
read

cd /root
rm -rf /var/www/html/okd4
rm -rf /root/install_dir

cd
mkdir install_dir
cp okd4_files/install-config.yaml install_dir

# TODO: Get pull secret from the internet
pullSecret_template='{\"auths\":{\"fake\":{\"auth\": \"bar\"}}}' 
pullSecret=$(cat pull_secret)
sed -i "s/$pullSecret_template/$pullSecret/g" install_dir/install-config.yaml

# Insert pubkey
pubkey="sshKey: "
pubkey+=\'$(cat .ssh/id_ed25519.pub)\'
sed -i '$ d' install_dir/install-config.yaml    # Simply delete the last line of the config ...
echo $pubkey >> install_dir/install-config.yaml # and insert a new one.

cat install_dir/install-config.yaml

# Create manifests and ignition configs.
openshift-install create manifests --dir=install_dir/
openshift-install create ignition-configs --dir=install_dir/

sudo mkdir /var/www/html/okd4

sudo cp -R install_dir/* /var/www/html/okd4/
sudo chown -R apache: /var/www/html/
sudo chmod -R 755 /var/www/html/

curl localhost:8080/okd4/metadata.json

cd /var/www/html/okd4/

COREOS_VERSION="32.20201018.3.0"

sudo wget https://builds.coreos.fedoraproject.org/prod/streams/stable/builds/$COREOS_VERSION/x86_64/fedora-coreos-$COREOS_VERSION-metal.x86_64.raw.xz
sudo wget https://builds.coreos.fedoraproject.org/prod/streams/stable/builds/$COREOS_VERSION/x86_64/fedora-coreos-$COREOS_VERSION-metal.x86_64.raw.xz.sig

sudo mv fedora-coreos-$COREOS_VERSION-metal.x86_64.raw.xz fcos.raw.xz
sudo mv fedora-coreos-$COREOS_VERSION-metal.x86_64.raw.xz.sig fcos.raw.xz.sig
sudo chown -R apache: /var/www/html/
sudo chmod -R 755 /var/www/html/
