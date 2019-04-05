#!/bin/bash
KUBE_VERSION=v1.13.4
KUBE_PAUSE_VERSION=3.1
ETCD_VERSION=3.2.24
DNS_VERSION=1.2.6
username=registry.cn-beijing.aliyuncs.com/syf-k8s

images=(kube-proxy:${KUBE_VERSION}
kube-scheduler:${KUBE_VERSION}
kube-controller-manager:${KUBE_VERSION}
kube-apiserver-amd64:${KUBE_VERSION}
kube-pause:${KUBE_PAUSE_VERSION}
kube-etcd:${ETCD_VERSION}
kube-coredns:${DNS_VERSION}
    )

for image in ${images[@]}
do
    docker pull ${username}/${image}
    docker tag ${username}/${image} k8s.gcr.io/${image}
    #docker tag ${username}/${image} gcr.io/google_containers/${image}
    docker rmi ${username}/${image}
done