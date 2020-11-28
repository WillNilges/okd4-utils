#!/bin/bash

set -e

# Credits (Thanks, Craig)
# https://medium.com/swlh/guide-okd-4-5-single-node-cluster-832693cb752b
# https://itnext.io/guide-installing-an-okd-4-5-cluster-508a2631cbee

# This is for single node clusters right now. Specifically, this is to be run on the services VM's root account.

# If this ever fires, you need to re-evaluate your life choices.
if [ $UID -ne 0 ]; then
    su root
fi

dnf update -y
dnf install -y git wget vim haproxy httpd bind bind-utils
cd ~
cd okd4_files

# TODO: Edit the named configs to work with your network.

mkdir /etc/named
/bin/cp -f named.conf /etc/named.conf
/bin/cp -f named.conf.local /etc/named/
mkdir /etc/named/zones
/bin/cp -f db* /etc/named/zones

systemctl enable named
systemctl start named
systemctl status named

firewall-cmd --permanent --add-port=53/udp
firewall-cmd --reload
nmcli connection modify ens18 ipv4.dns "127.0.0.1"
systemctl restart NetworkManager

# Hope that worked!
dig okd.local
dig –x 192.168.60.240

cp haproxy.cfg /etc/haproxy/haproxy.cfg

setsebool -P haproxy_connect_any 1
systemctl enable haproxy
systemctl start haproxy
systemctl status haproxy

firewall-cmd --permanent --add-port=6443/tcp
firewall-cmd --permanent --add-port=22623/tcp
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
firewall-cmd --reload

sed -i 's/Listen 80/Listen 8080/' /etc/httpd/conf/httpd.conf

setsebool -P httpd_read_user_content 1
systemctl enable httpd
systemctl start httpd
firewall-cmd --permanent --add-port=8080/tcp
firewall-cmd --reload

curl localhost:8080

cd
wget https://github.com/openshift/okd/releases/download/4.5.0-0.okd-2020-07-29-070316/openshift-client-linux-4.5.0-0.okd-2020-07-29-070316.tar.gz
wget https://github.com/openshift/okd/releases/download/4.5.0-0.okd-2020-07-29-070316/openshift-install-linux-4.5.0-0.okd-2020-07-29-070316.tar.gz

tar -zxvf openshift-client-linux-4.5.0-0.okd-2020-07-29-070316.tar.gz
tar -zxvf openshift-install-linux-4.5.0-0.okd-2020-07-29-070316.tar.gz

mv kubectl oc openshift-install /usr/local/bin/
oc version
openshift-install version

ssh-keygen

cd
mkdir install_dir
cp okd4_files/install-config.yaml ./install_dir

# TODO: sed pull secret
echo "Go put your pull secret in ~/pull_secret. Press ENTER to continue."
read
cat pull_secret >> install_dir/install-config.yaml
cat .ssh/id_rsa.pub >> install_dir/install-config.yaml
vim ./install_dir/install-config.yaml # Fix whatever needs fixing.
cp ./install_dir/install-config.yaml ./install_dir/install-config.yaml.bak

openshift-install create manifests --dir=install_dir/
openshift-install create ignition-configs --dir=install_dir/

sudo mkdir /var/www/html/okd4

sudo cp -R install_dir/* /var/www/html/okd4/
sudo chown -R apache: /var/www/html/
sudo chmod -R 755 /var/www/html/

curl localhost:8080/okd4/metadata.json

cd /var/www/html/okd4/
sudo wget https://builds.coreos.fedoraproject.org/prod/streams/stable/builds/32.20200715.3.0/x86_64/fedora-coreos-32.20200715.3.0-metal.x86_64.raw.xz
sudo wget https://builds.coreos.fedoraproject.org/prod/streams/stable/builds/32.20200715.3.0/x86_64/fedora-coreos-32.20200715.3.0-metal.x86_64.raw.xz.sig
sudo mv fedora-coreos-32.20200715.3.0-metal.x86_64.raw.xz fcos.raw.xz
sudo mv fedora-coreos-32.20200715.3.0-metal.x86_64.raw.xz.sig fcos.raw.xz.sig
sudo chown -R apache: /var/www/html/
sudo chmod -R 755 /var/www/html/

# Alright, go ignite your cluster. When you get to the API part, run this command:
# oc patch etcd cluster -p='{"spec": {"unsupportedConfigOverrides": {"useUnsupportedUnsafeNonHANonProductionUnstableEtcd": true}}}' --type=merge
# That allows the use of only one ctrl plane.

# Or, if you've got multiple, then here is what your ignition entries will look like:
# that is, unless you're using a router or something smart like that.

# Bootstrap
#ip=10.10.3.37::10.10.3.1:255.255.255.0:::none nameserver=10.10.3.33 coreos.inst.install_dev=/dev/sda coreos.inst.image_url=http://10.10.3.33:8080/okd4/fcos.raw.xz coreos.inst.ignition_url=http://10.10.3.33:8080/okd4/bootstrap.ign

# Ctrl01
#ip=10.10.3.34::10.10.3.1:255.255.255.0:::none nameserver=10.10.3.33 coreos.inst.install_dev=/dev/sda coreos.inst.image_url=http://10.10.3.33:8080/okd4/fcos.raw.xz coreos.inst.ignition_url=http://10.10.3.33:8080/okd4/master.ign

# Ctrl02
#ip=10.10.3.35::10.10.3.1:255.255.255.0:::none nameserver=10.10.3.33 coreos.inst.install_dev=/dev/sda coreos.inst.image_url=http://10.10.3.33:8080/okd4/fcos.raw.xz coreos.inst.ignition_url=http://10.10.3.33:8080/okd4/master.ign

# Ctrl03
#ip=10.10.3.36::10.10.3.1:255.255.255.0:::none nameserver=10.10.3.33 coreos.inst.install_dev=/dev/sda coreos.inst.image_url=http://10.10.3.33:8080/okd4/fcos.raw.xz coreos.inst.ignition_url=http://10.10.3.33:8080/okd4/master.ign

exit 0