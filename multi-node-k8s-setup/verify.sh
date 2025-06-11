# Verificar containerd
sudo crictl version
sudo crictl info

# Verificar módulos do kernel
lsmod | grep br_netfilter
lsmod | grep overlay

# Verificar configurações sysctl
sysctl net.bridge.bridge-nf-call-iptables
sysctl net.ipv4.ip_forward

# Verificar versão do Kubernetes
kubectl version --client
kubeadm version
