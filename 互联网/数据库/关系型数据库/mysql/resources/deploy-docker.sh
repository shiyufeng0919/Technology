#! /bin/bash
#cose firewall
systemctl stop firewalld.service
systemctl disable firewalld.service

# disabled SELINUX
setenforce 0

cat <<EOF > /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF

modprobe br_netfilter
sysctl -p /etc/sysctl.d/k8s.conf

# install docker
yum install -y yum-utils device-mapper-persistent-data lvm2

yum-config-manager \
      --add-repo \
      https://download.docker.com/linux/centos/docker-ce.repo

yum install -y --setopt=obsoletes=0 \
        docker-ce-17.03.2.ce-1.el7.centos \
        docker-ce-selinux-17.03.2.ce-1.el7.centos

mkdir -p /etc/systemd/system/docker.service.d
cat <<EOF > /etc/systemd/system/docker.service.d/port.conf
ExecStartPost=/usr/sbin/iptables -P FORWARD ACCEPT
EOF

systemctl daemon-reload && systemctl restart docker