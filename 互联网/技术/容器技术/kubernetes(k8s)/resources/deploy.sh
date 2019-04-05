#! /bin/bash

# close firewall
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

# install kubeadm & kubelet
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg

EOF

yum install -y kubelet-1.13.3 kubectl-1.13.3 kubernetes-cni-0.6.0 kubeadm-1.13.3

# deal swap
swapoff -a
echo 'vm.swappiness=0' >>/etc/sysctl.d/k8s.conf
sysctl -p /etc/sysctl.d/k8s.conf
systemctl enable kubelet.service && systemctl start kubelet.service


===================================================================

### 总结：本机kubeadm安装kubernetes集群完整版

#### 1.硬件环境: centos7  [vagrantfile]

    机器--注意：需要与本机IP端端一致，否则物理主机与虚拟机间无法通信

    11.12.81.11 k8s1
    11.12.81.12 k8s2
    11.12.81.13 k8s3

#### 2.ssh免登录:

    $ vim /etc/ssh/sshd_config #修改3处

        PermitRootLogin yes
        PubkeyAuthentication yes
        PasswordAuthentication yes

    $ ssh-keygen #生成密钥

    $ systemctl restart sshd #重启

    **上述均在各节点上执行**

    $ ssh-copy-id root@k8s2   #相互拷备公钥

    $ ssh k8s2   #当前k8s1身份即可直接登录k8s2

#### 3.部署流程--各节点均需执行

[部署流程](resources/deploy.sh)

#### 4.[master节点初始化集群](resources/initk8s.sh)

生成：kubeadm join 11.12.81.11:6443 --token 7i7ld2.ot25zym047k9szh9 --discovery-token-ca-cert-hash sha256:8aafb6689a539c07d46f4b1ebb37e99dfbb211cdaadc49f12b3af2d5dd7759c8

#### 5.k8s2节点加入集群

    $ ssh k8s2
    $ kubeadm join 11.12.81.11:6443 --token 7i7ld2.ot25zym047k9szh9 --discovery-token-ca-cert-hash sha256:8aafb6689a539c07d46f4b1ebb37e99dfbb211cdaadc49f12b3af2d5dd7759c8

    $ ssh k8s1

    #查看Nod节点状态
    [root@k8s1 k8s]# kubectl get nodes
    NAME        STATUS   ROLES    AGE   VERSION
    k8s1.node   Ready    master   13m   v1.13.3
    k8s2.node   Ready    <none>   12m   v1.13.3

    #确保所有pod都处于Running状态
    [root@k8s1 k8s]# kubectl get pod --all-namespaces -o wide
    NAMESPACE     NAME                                       READY   STATUS    RESTARTS   AGE     IP              NODE        NOMINATED NODE   READINESS GATES
    kube-system   calico-kube-controllers-55df754b5d-jw8wh   1/1     Running   0          8m14s   192.168.197.3   k8s2.node   <none>           <none>
    kube-system   calico-node-6dw8n                          1/1     Running   0          8m14s   10.0.2.15       k8s2.node   <none>           <none>
    kube-system   calico-node-wthzc                          1/1     Running   0          8m14s   10.0.2.15       k8s1.node   <none>           <none>
    kube-system   coredns-78d4cf999f-t7xn5                   1/1     Running   0          13m     192.168.197.2   k8s2.node   <none>           <none>
    kube-system   coredns-78d4cf999f-z5lkz                   1/1     Running   0          13m     192.168.197.1   k8s2.node   <none>           <none>
    kube-system   etcd-k8s1.node                             1/1     Running   0          13m     10.0.2.15       k8s1.node   <none>           <none>
    kube-system   kube-apiserver-k8s1.node                   1/1     Running   0          13m     10.0.2.15       k8s1.node   <none>           <none>
    kube-system   kube-controller-manager-k8s1.node          1/1     Running   0          12m     10.0.2.15       k8s1.node   <none>           <none>
    kube-system   kube-proxy-jrw68                           1/1     Running   0          13m     10.0.2.15       k8s1.node   <none>           <none>
    kube-system   kube-proxy-vh6x7                           1/1     Running   0          13m     10.0.2.15       k8s2.node   <none>           <none>
    kube-system   kube-scheduler-k8s1.node                   1/1     Running   0          13m     10.0.2.15       k8s1.node   <none>           <none>
