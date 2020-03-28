# K8s Contest Managment System
Make deploying [CMS](https://github.com/cms-dev/cms) great again.

## Intro
`k8s-cms` ports the [Contest Management System](https://github.com/cms-dev/cms) to Kubernetes

Features:
- no more having to follow long deployment instructions :tada:
- adds fault tolerance through automatic health checks and self healing
- autoscales to larger contests:
    - supports running up to 24 workers
    - supports running multiple contest web servers
- includes `kcmscli` CLI to simply contest setup process 
    (ie importing users from CSV, importing contests)

Limitations:
- multiple contests - only supports running one contest at a time
- printing - hooking up printers to print stuff has not been implemented yet.
- scaling more than 24 worker instances - only supports scaling up to 24 worker instances
- K8s only: requires the cluster to support privileged containers

## Quickstart
Deploy k8s-cms:
1. Clone or Download the repository
2. Fill `env` file with credentials
```sh
cp env .env
nano .env # use your favourite editor
```
3. Deploy
- On Kubernetes 
```sh
helm install k8s-cms .
```
- On Docker-Compose
```sh
docker-compose pull
docker-compose up
```

> [More Deployment Details](./docs/deploy.md)
