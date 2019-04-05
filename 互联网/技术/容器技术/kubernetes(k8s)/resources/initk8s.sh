kubeadm init \
    --apiserver-advertise-address=10.13.80.11 \
    --image-repository registry.aliyuncs.com/google_containers \
    --kubernetes-version v1.13.3 \
    --service-cidr=10.1.0.0/16 \
    --pod-network-cidr=10.244.0.0/16