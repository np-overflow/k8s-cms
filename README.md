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
This repository contains submodules so:
```
git submodule update --init --recursive
```
after cloning.

## Roadmap
- dockerizing all these:
    - Database
    - ResourceService, LogService, ScoringService,ProxyService, PrintingService, AdminWebServer, RankingWebServer, Checker
    - ContestWebServer
    - Worker - requires language support
        - C
        - C++
        - Java
        - Pascal
        - Python with zip executable
        - PHP
        - Rust 
        - C# 
- write k8s deployment YAMLs all these:
    - Database
    - ResourceService, LogService, ScoringService,ProxyService, PrintingService,  AdminWebServer, RankingWebServer, Checker
    - ContestWebServer
    - Worker - requires language support
        - C
        - C++
        - Java
        - Pascal
        - Python with zip executable
        - PHP
        - Rust 
        - C# 
- securing k8s-cms:
    - data storage encryption
    - k8s communication encryption.
