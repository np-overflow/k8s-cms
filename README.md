# K8s Contest Managment System
Make deploying [CMS](https://github.com/cms-dev/cms) great again.

## Intro
The Contest Managment System (CMS) is a great open source platform to host programming contests. 

Adopting CMS to run on kubernetes brings the following benefits:
- significantly simplifies the deployment process
- adds fault tolerance through automatic health checks and self recovery
- scales to larger contests:
    - supports running up to 24 workers
    - supports running multiple contest web servers

## Deploy
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
    - ingress controller optional for ingress support.
- Requires [Helm](https://helm.sh/docs/using_helm/#installing-helm) to be up and running. 
    - For security use [Tillerless Helm](https://github.com/rimusz/helm-tiller)

```sh
helm tiller start # tillerless helm only
helm install .
helm tiller stop # tillerless helm only
```

#### Optional Addons
Optionally configure addons the following in `values.yaml` before `helm install .`
- Expose services with ingress:
    - deploy nginx ingress controller setting `nginx-ingress.enabled` to `true`
    - set `ingress.enabled` to `true` and configure dns hosts
- Automatically provision TLS certificates for HTTPs:
    - expose services and test services with ingress first
    - deploy cert-manage by setting `cert-manager.enabled` to `true`
    - set `certGenerate.enabled` to `true` and configure email
- Mointoring with monitoring and alerts and prometheus and grafana 
    - deploy cert-manage by setting `prometheus-operato.enabled` to `true`
    - Port forward the grafana service to access monitoing dashboards.

### Docker-Compose
Runs CMS on a single machine. Suitable for testing:
- Only requires `docker` and `docker-compose`. No kubernetes required.
- Limited to only 2 workers.
```sh
docker-compose pull 
docker-compose up
```

## Design
Each CMS service to containerized by its own docker container:
- Database - Deploy using Postgres SQL container `cms-db`
- CMS - all services derive from base container `cms-base`
    - ResourceService - `cms-resource`
    - LogService - `cms-log`
    - EvaluationService - `cms-evaluation`
    - ScoringService - `cms-scoring`
    - ProxyService - `cms-proxy`
    - PrintingService - `cms-printing`
    - AdminWebServer - `cms-web-admin`
    - RankingWebServer - `cms-web-ranking`
    - Checker - `cms-checker`
    - ContestWebServer - `cms-web-contest`
    - Worker - `cms-worker` requires language support

> `cms-base` contains python runtime, copy of cms source code and `cms.conf`
>  and is used a a base to build the other services

### Security
Security Measures:
- internal service communicate on a virtual network are inaccessable to participants.
- Secrets are injected into the containers as environment variables via `.env` file.
- All services (except database) run as an unprivilleged user.

Security Concerns:
- `cms-worker` runs as a privileged container as the `isolate` sandbox requires 
    privileged access to the system.
- `helm`'s Tiller uses an exposes a insecure GRPC port with cluster wide admin 
    rights. Use [Tillerless Helm](https://github.com/rimusz/helm-tiller) to run
    Tiller locally for security.

### Limitations
Limitations:
- multiple contests - only supports running one contest at a time
- printing - hooking up printers to print stuff has not been implemented yet.
- importing contests - importing contests has not been  implmemented yet.
- scaling more than 24 instances - only supports scaling up to 24 worker instances
- requires the cluster to support privileged containers

## Contributing
Guidelines for contributors:
- proposed changes: `TODO.md`
- project changelog: `CHANGELOG.md`
 
Development setup for contributors:
1. Resolve submodules after cloning;
```sh
git submodule update --init --recursive
```
2. Fill `env` file with secrets
```sh
cp env .env
nano .env # use your favourite editor
```
3. Run the stack
```sh
docker-compose up # use docker-compose OR

# use kubernetes
helm tiller start # tillerless helm only
helm install .
helm tiller stop # tillerless helm only
```
