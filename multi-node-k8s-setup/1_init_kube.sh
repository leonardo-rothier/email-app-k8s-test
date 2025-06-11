# Initizalize Kubeadm :
sudo kubeadm init --pod-network-cidr=10.244.0.0/16 --apiserver-advertise-address=192.168.0.100
# --pod-network-cidr the default cidr for calico is 192.168.0.0/16
# --apiserver-advertise-address guarantees that the API server announces itself 
# on the correct network address (use `ip addr` and get your master node ip)

# Configure kubectl access to your user:
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# Install the pods network (example - calico):
kubectl apply -f https://docs.projectcalico.org/manifests/calico.yaml
