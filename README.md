# Email App K8s Test Environment 📧
A simple Kubernetes project demonstrating communication between frontend and backend pods. This is an educational environment designed to help understand core Kubernetes concepts like Deployments, Services, Secrets, and pod-to-pod communication.

## Project Overview
This project simulates a basic email application with:

**Frontend:** A simple web interface (email-site)
**Backend:** An email service API (email-service) that integrates with Gmail  

> *Note: This is a test environment for explore Kubernetes. The code is intentionally simple and not production-ready.*

## Architecture
```
Internet
    |
    v
NodePort (30000)
    |
    v
email-site-service
    |
    v
email-site pods (2 replicas)
    |
    v
email-service (ClusterIP)
    |
    v
email-service pods (3 replicas)
    |
    v
Gmail API
```

## Setup Instructions
### Create Kubernetes Secret
Create the Gmail credentials secret that email-service backend will use (depends on yours senders):
```bash
kubectl create secret generic gmail-sender-credentials \
--from-literal=username=your-email@gmail.com \
--from-literal=password=your-app-password
```

### Clone repo
```bash
git clone https://github.com/leonardo-rothier/email-app-k8s-test.git
cd email-app-k8s-test
```

### Apply k8s definitions
```bash
kubectl apply -f k8s-specifications/deployments/
kubectl apply -f k8s-specifications/services/
```

### Get the status of deployed components

```bash
kubectl get deployments,svc,pods
```

### Monitoring

For monitoring we are going to use prometheus, for this, with helm installed, run:
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm install prometheus prometheus-community/kube-prometheus-stack --namespace monitoring --create-namespace
```

#### Expose Prometheus/Grafana
To expose the pod instances of prometheus and grafana change its services to NodePort editing:
```bash
kubectl edit svc/prometheus-kube-prometheus-prometheus
kubectl edit svc/prometheus-grafana
```
#### Necessary configuration to scrap metrics
KubeProxy to be acessible ouside its localhost, if you are using multiples nodes:
```bash
# edit the metricsBindAddress as needed for your case
kubectl edit configmap -n kube-system kube-proxy
```
The same with the etcd config file at /etc/kubernetes/manifests/etcd.yaml:

```bash
# Change the value inside the --listen-metrics-urls as needed
sudo vim /etc/kubernetes/manifests/etcd.yaml
```

And edit the kube-scheduler manifest:
```bash
# Change the value inside the --bind-address as needed
sudo vim /etc/kubernetes/manifests/kube-scheduler.yaml
```

#### Additional scrape configs (for your applications)
Create a ServiceMonitor and apply it, Service Monitors define a set of targets for prometheus to monitor and scrape:
```bash
kubectl apply -f k8s-specifications/monitoring/
```
The release label is necessary and you get this information under:
```bash
# This label allows prometheus to find service monitors in the cluster
kubectl get prometheuses.monitoring.coreos.com -o yaml -n monitoring | grep -A3 serviceMonitorSelector:
```

#### Grafana get credentials
```bash
kubectl get secret prometheus-grafana -n default -o jsonpath="{.data.admin-user}" | base64e --decode; echo
kubectl get secret prometheus-grafana -n default -o jsonpath="{.data.admin-password}" | base64 --decode; echo
```/