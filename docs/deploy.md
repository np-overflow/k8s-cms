# Deploy
Deployment is supported on the following platforms:
- Kubernetes
- Docker-Compose

Common instructions:
1. Clone or Download the repository
2. Fill `env` file with credentials
```sh
cp env .env
nano .env # use your favourite editor
```

### Kubernetes
Runs CMS on Kubernetes cluster. Suitable for hosting actual contests:
- Requires Kubernetes Cluster with the following:
    - preconfigured default storage class (check with `kubectl get sc`)
    - ingress (nginx-ingress available as optional addon)

```sh
helm install k8s-cms chart
```

> Additional deployment configuration in `values.yaml`:
> - change `replicas` to control the no. of replicas created for autoscaling

#### K8s: Optional Addons
Optionally configure addons the following in `values.yaml` before `helm install .`:

- Import addon charts before enabling addons:

```sh
helm repo add stable https://kubernetes-charts.storage.googleapis.com
helm repo add jetstack https://charts.jetstack.io
helm dep update
```

- Expose services with ingress
    - deploy nginx ingress controller setting `nginx-ingress.enabled` to `true`
    - set `ingress.enabled` to `true` and configure dns hosts
- Automatically provision TLS certificates for HTTPs:
    - expose services and test services with ingress first
    - deploy cert-manage by setting `cert-manager.enabled` to `true`
    - set `certGenerate.enabled` to `true` and configure email
- Mointoring with monitoring and alerts and prometheus and grafana 
    - deploy cert-manage by setting `prometheus-operator.enabled` to `true`
    - Port forward the grafana service to access monitoing dashboards.

### Docker-Compose
Runs CMS on a single machine. Suitable for testing:
- Only requires `docker` and `docker-compose`. No kubernetes required.
- Limited to only 2 workers.
```sh
docker-compose pull 
docker-compose up
```
