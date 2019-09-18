# K8s Contest Managment System
Make deploying [CMS](https://github.com/cms-dev/cms) great again.

## Intro
The Contest Managment System (CMS) is a great open source platform to host programming contests. 
However deploying it is a certified pain. 

By adapting to be deployed using `kubernetes`, we can make deploying CMS as:
```
kubectl apply -f https://raw.github.com.... TODO
```

## Setup
This repository contains submodules - run after cloning;
```
git submodule update --init --recursive
```

## Design
![k8s-cms Design](./assets/k8s_cms_design.jpg)

Containers:
- Database - Deploy using Postgres SQL container
- CMS - all services derive from base container `cms-base`
    - ResourceService - `cms-resource`
    - LogService - `cms-log`
    - ScoringService - `cms-scoring`
    - ProxyService - `cms-proxy`
    - PrintingService - `cms-printing`
    - AdminWebServer - `cms-web-admin`
    - RankingWebServer - `cms-ranking`
    - Checker - `cms-checker`
    - ContestWebServer - `cms-web-contest`
    - Worker - `cms-worker` requires language support

> `cms-base` contains python runtime, copy of cms source code and `cms.conf`

## Project Structure
- deps/ - project dependencies/submodules
- containers/ - docker containers
- k8s/ - kubernetes deployemennt

## Roadmap
- dockerizing all these:
    - Database 
    - ResourceService
    - LogService
    - ScoringService
    - ProxyService
    - PrintingService
    - AdminWebServer
    - RankingWebServer
    - Checker
    - ContestWebServer
    - Worker - requires language support
        - C C++ Java Pascal Python with zip executable PHP Rust C# 
- write k8s deployment YAMLs all these:
    - Database
    - ResourceService
    - LogService
    - ScoringService
    - ProxyService
    - PrintingService
    - AdminWebServer
    - RankingWebServer
    - Checker
    - ContestWebServer
    - Worker - requires language support
        - C C++ Java Pascal Python with zip executable PHP Rust C# 
- figure out contest & deployment
- making k8s-cms scalable:
    - generate `cms.conf` using kubernetes deployment/docker-compose  file.
    - scaling `ContestWebServer` and other components live.
- securing k8s-cms:
    - data storage encryption
    - k8s communication encryption.
    - HTTPs for RankingWebServer,AdminWebServer,ContestWebServer.
