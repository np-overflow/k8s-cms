# K8s Contest Managment System
Make deploying [CMS](https://github.com/cms-dev/cms) great again.

## Intro
The Contest Managment System (CMS) is a great open source platform to host programming contests. 
However deploying it is a certified pain. 

By adapting to be deployed using `kubernetes`, we can make deploying CMS as:
```
kubectl apply -f https://raw.github.com.... TODO
```

### CMS itself
Services in the CMS:
- LogService: collects all log messages in a single place;
- ResourceService: collects data about the services running on the same server, and takes care of starting all of them with a single command;
- Checker: simple heartbeat monitor for all services;
- EvaluationService: organizes the queue of the submissions to compile or evaluate on the testcases, and dispatches these jobs to the workers;
- Worker: actually runs the jobs in a sandboxed environment;
- ScoringService: collects the outcomes of the submissions and computes the score;
- ProxyService: sends the computed scores to the rankings;
- PrintingService: processes files submitted for printing and sends them to a printer;
- AdminWebServer: the webserver to control and modify the parameters of the contests.
- RankingWebServer: displays contests ranking.
- ContestWebServer: the webserver that the contestants will be interacting with;


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
