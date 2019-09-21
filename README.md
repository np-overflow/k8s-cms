# K8s Contest Managment System
Make deploying [CMS](https://github.com/cms-dev/cms) great again.

## Intro
The Contest Managment System (CMS) is a great open source platform to host programming contests. 
However deploying it is also really hard.

By adapting CMS to be deployed using `kubernetes`, k8s-cms can make deploying 
CMS simple as:
```
kubectl apply -f https://raw.github.com.... TODO
```

## Setup

### Single Machine with Docker-Compose
Setups CMS on a single machine. Suitable for testing.
- Only requires `docker` and `docker-compose`. No kubernetes required.
1. Fill `env` file with secrets
```sh
curl https://raw.githubusercontent.com/np-overflow/k8s-cms/master/env -o env
cp env .env
nano .env # use your favourite editor
```
2. Run CMS on a single machine with `docker-compose`:
```sh
curl https://raw.githubusercontent.com/np-overflow/k8s-cms/master/docker-compose.yml -o docker-compose.yaml
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
Making k8s-cms secure:
- internal service communicate on a virtual network are inaccessable to participants.
- Secrets are injected into the containers as environment variables via `.env` file.

### Limitations
What does not work:
- multiple contests - only supports running one contest at a time
- printing - hooking up printers to print stuff has not been implemented yet.
- importing contests - importing contests has not been  implmemented yet.

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
docker-compose up
```
