#! /bin/bash

# close firewall
apt-get install ufw
ufw enable
ufw default deny

# install docker
sudo apt-get update
sudo apt-get -y install apt-transport-https ca-certificates curl software-properties-common
curl -fsSL http://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] http://mirrors.aliyun.com/docker-ce/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get -y update
apt -y install docker-ce=17.03.1~ce-0~ubuntu-xenial
systemctl start docker