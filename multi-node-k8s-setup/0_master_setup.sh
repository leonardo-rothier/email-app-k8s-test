#!/bin/bash
set -e

# disable swap (exigência do Kubernetes)
swapoff -a
sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab

# set SELinux for permissive during setup
setenforce 0
sed -i 's/^SELINUX=.*/SELINUX=permissive/' /etc/selinux/config

# add all nodes to hosts (Adjust the ips following your vms ips)
tee -a /etc/hosts <<EOF
192.168.0.100 master
192.168.0.253 worker2
192.168.0.254 worker1
EOF

# enable modules for k8s
tee /etc/modules-load.d/k8s.conf <<EOF
overlay
br_netfilter
EOF
modprobe overlay
modprobe br_netfilter

# configure sysctl for Kubernetes network requirements
tee /etc/sysctl.d/k8s.conf <<EOF
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF
sysctl --system

# Install and configure containerd (k8s default runtime)
dnf -y install dnf-plugins-core
dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
dnf install -y containerd.io
mkdir -p /etc/containerd
containerd config default | tee /etc/containerd/config.toml
sed -i 's/SystemdCgroup = false/SystemdCgroup = true/' /etc/containerd/config.toml

systemctl enable --now containerd

# Configure crictl
tee /etc/crictl.yaml <<EOF
runtime-endpoint: unix:///run/containerd/containerd.sock
image-endpoint: unix:///run/containerd/containerd.sock
timeout: 10
EOF

# Add Kubernetes repo
cat <<EOF | tee /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://pkgs.k8s.io/core:/stable:/v1.33/rpm/
enabled=1
gpgcheck=1
gpgkey=https://pkgs.k8s.io/core:/stable:/v1.33/rpm/repodata/repomd.xml.key
exclude=kubelet kubeadm kubectl cri-tools kubernetes-cni
EOF

# install kubeadm, kubelet e kubectl
dnf install -y kubelet kubeadm kubectl --disableexcludes=kubernetes
systemctl enable --now kubelet

# FIREWALL (adjust for security or remove for labs)
systemctl enable --now firewalld

# Master recommended ports:
firewall-cmd --permanent --add-port=6443/tcp      # Kubernetes API server
firewall-cmd --permanent --add-port=2379-2380/tcp # etcd server client API
firewall-cmd --permanent --add-port=10250/tcp     # kubelet API
firewall-cmd --permanent --add-port=10251/tcp     # kube-scheduler
firewall-cmd --permanent --add-port=10252/tcp     # kube-controller-manager
firewall-cmd --reload


#### kubeadm join 192.168.1.31:6443 --token ********* command need to be add"
