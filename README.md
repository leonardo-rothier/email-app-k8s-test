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
Create the Gmail credentials secret that the backend will use:
```bash
kubectl create secret generic gmail-credentials \
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

helm install prometheus prometheus-community/kube-prometheus-stack
```

### EXTRA: EKS cluster
If you pretend to run in a Hosted (Managed) solution as EKS, it will still works but with a slight change on the email-site-service, change the type to `LoadBalancer` and remove the `nodePort`.  
Even better, use an ingress, is more cost-effective for multiple services. For this, your email-site-service will become a ClusterIP(default type) as we made with the backend, and you will need to apply:
```bash
kubectl apply -f k8s-specifications/ingress/email-ingress.yml
```