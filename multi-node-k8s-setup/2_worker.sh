# disable swap
sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab
swapoff -a

# configure SELINUX
setenforce 0
sed -i 's/^SELINUX=.*/SELINUX=permissive/' /etc/selinux/config

# Set dns local entries
tee -a /etc/hosts <<EOF
192.168.0.100 master
192.168.0.253 worker2
192.168.0.254 worker1
EOF

# Configure firewall (worker node ports)
systemctl enable --now firewalld
firewall-cmd --permanent --add-port=10250/tcp       # kubelet API
firewall-cmd --permanent --add-port=30000-32767/tcp # NodePort Services
firewall-cmd --permanent --add-port=179/tcp         # enable BGP mesh between nodes
sudo firewall-cmd --permanent --add-protocol=4
sudo firewall-cmd --permanent --add-rich-rule="rule family=ipv4 source address=10.244.0.0/16 accept"
sudo firewall-cmd --permanent --add-rich-rule="rule family=ipv4 source address=10.96.0.0/12 accept"
firewall-cmd --reload

# Load kernel modules
cat <<EOF | tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF

modprobe overlay
modprobe br_netfilter

# configure sysctl for Kubernetes network requirements
cat <<EOF | tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF

sysctl --system


# Setup containerd repository and install
dnf -y install dnf-plugins-core
dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
dnf install -y containerd.io

# Configure containerd
mkdir -p /etc/containerd
containerd config default | tee /etc/containerd/config.toml
sed -i 's/SystemdCgroup = false/SystemdCgroup = true/g' /etc/containerd/config.toml

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

dnf install -y kubelet kubeadm kubectl --disableexcludes=kubernetes

systemctl enable --now kubelet

echo "Worker node setup completed! Ready for 'kubeadm join'"
# if you losted the join token your can create a new one running:
# kubeadm token create --print-join-command (on the master node)

