#!/bin/bash

# This script is used to stand up a fresh okd4-services node. This node will handle HAProxy, named, and httpd services for the OKD setup process. Basically, modify the supplied configuration files, then run this script. It'll configure all the ports, install all the things, and copy all the files. When it's done, it'll print "==READY FOR IGNITION==" When you see that, use the ignite.sh script to write your ignition nodes.

# See https://getfedora.org/en/coreos/download?tab=metal_virtualized&stream=stable for new versions.
COREOS_VERSION="32.20201018.3.0"

# See https://github.com/openshift/okd/releases/ for the most up to date info.
OKD_VERSION="4.7.0-0.okd-2021-08-22-163618"

echo "ARE YOU SURE YOU WANT TO DO THIS? (press Enter if you do, ^C if you don't.)"
read

set -e

# Credits (Thanks, Craig)
# https://medium.com/swlh/guide-okd-4-5-single-node-cluster-832693cb752b
# https://itnext.io/guide-installing-an-okd-4-5-cluster-508a2631cbee

# This is for single node clusters right now. Specifically, this is to be run on the services VM's root account.

if [ $UID -ne 0 ]; then
    su root
fi

dnf update -y
dnf install -y git wget vim neovim haproxy httpd bind bind-utils
cd 
cd okd4_files

# TODO: Edit the named configs to work with your network.

# Configure named
rm /etc/named; mkdir /etc/named
/bin/cp -f named.conf /etc/named.conf
/bin/cp -f named.conf.local /etc/named/
mkdir /etc/named/zones
/bin/cp -f db* /etc/named/zones

systemctl enable named
systemctl start named
#systemctl status named

firewall-cmd --permanent --add-port=53/udp
firewall-cmd --reload
nmcli connection modify ens18 ipv4.dns "127.0.0.1"
systemctl restart NetworkManager

# Hope that worked!
dig postave.us
dig -x 10.10.33.70

# Configure HAProxy
/bin/cp -f haproxy.cfg /etc/haproxy/haproxy.cfg

setsebool -P haproxy_connect_any 1
systemctl enable haproxy
systemctl start haproxy
#systemctl status haproxy

firewall-cmd --permanent --add-port=6443/tcp
firewall-cmd --permanent --add-port=22623/tcp
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
firewall-cmd --reload

# Configure httpd
sed -i 's/Listen 80\n/Listen 8080\n/' /etc/httpd/conf/httpd.conf

setsebool -P httpd_read_user_content 1
systemctl enable httpd
systemctl start httpd
firewall-cmd --permanent --add-port=8080/tcp
firewall-cmd --reload

curl localhost:8080

# Install OKD 4
cd
OKD_CLIENT_LINUX="https://github.com/openshift/okd/releases/download/$OKD_VERSION/openshift-client-linux-$OKD_VERSION.tar.gz"
OKD_LINUX="https://github.com/openshift/okd/releases/download/$OKD_VERSION/openshift-install-linux-$OKD_VERSION.tar.gz"

wget $OKD_CLIENT_LINUX
wget $OKD_LINUX

tar -zxvf openshift-client-linux-$OKD_VERSION.tar.gz
tar -zxvf openshift-install-linux-$OKD_VERSION.tar.gz

mv kubectl oc openshift-install /usr/local/bin/
oc version
openshift-install version

# Setup your install config
#ssh-keygen -t ed25519

cd
rm -rf install_dir; mkdir install_dir
cp okd4_files/install-config.yaml install_dir

# TODO: Get pull secret
# cat pull_secret >> install_dir/install-config.yaml
# cat .ssh/id_*.pub >> install_dir/install-config.yaml
pullSecret_template='{\"auths\":{\"fake\":{\"auth\": \"bar\"}}}' 
pullSecret=$(cat pull_secret)
pubkey_template='ssh-ed25519\ AAAA...'
pubkey="sshKey: "
pubkey+=$(cat .ssh/id_ed25519.pub)
sed -i "s/$pullSecret_template/$pullSecret/g" install_dir/install-config.yaml
sed -i '$ d' install_dir/install-config.yaml
echo $pubkey >> install_dir/install-config.yaml
cp install_dir/install-config.yaml install_dir/install-config.yaml.bak


# Create manifests and ignition configs.
openshift-install create manifests --dir=install_dir/
# sed -i 's/mastersSchedulable: true/mastersSchedulable: False/' install_dir/manifests/cluster-scheduler-02-config.yml # Make masters unschedulable (Don't use this if you don't have any workers.)
openshift-install create ignition-configs --dir=install_dir/

rm -rf /var/www/html/okd4; mkdir /var/www/html/okd4

cp -R install_dir/* /var/www/html/okd4/
chown -R apache: /var/www/html/
chmod -R 755 /var/www/html/

curl localhost:8080/okd4/metadata.json

# Acquire and set up Fedora CoreOS ISOs and prepare for Ignition.

cd /var/www/html/okd4/

sudo wget https://builds.coreos.fedoraproject.org/prod/streams/stable/builds/$COREOS_VERSION/x86_64/fedora-coreos-$COREOS_VERSION-metal.x86_64.raw.xz
sudo wget https://builds.coreos.fedoraproject.org/prod/streams/stable/builds/$COREOS_VERSION/x86_64/fedora-coreos-$COREOS_VERSION-metal.x86_64.raw.xz.sig

sudo mv fedora-coreos-$COREOS_VERSION-metal.x86_64.raw.xz fcos.raw.xz
sudo mv fedora-coreos-$COREOS_VERSION-metal.x86_64.raw.xz.sig fcos.raw.xz.sig
sudo chown -R apache: /var/www/html/
sudo chmod -R 755 /var/www/html/

echo '== READY FOR IGNITION =='

exit 0
