# k8s-cms Roadmap

## Version 0.1.0 alpha
- dockerizing all these:
    - Database  :heavy_check_mark:
    - ResourceService :heavy_check_mark:
    - LogService :heavy_check_mark:
    - ScoringService :heavy_check_mark:
    - ProxyService - with single contest support limitation :heavy_check_mark:
    - EvaluationService :heavy_check_mark:
    - PrintingService :heavy_check_mark:
    - AdminWebServer :heavy_check_mark:
    - RankingWebServer :heavy_check_mark:
    - Checker :heavy_check_mark:
    - ContestWebServer - with single contest limitation :heavy_check_mark:
    - Worker - requires language support :heavy_check_mark:
        - C C++ Java Pascal Python with zip executable PHP Rust C# 
- docker-compose single machine support :heavy_check_mark:

## Version 0.2.0 alpha
- kubernetes support - write k8s deployment YAMLs all these:
    - Database
    - ResourceService
    - LogService
    - ScoringService
    - EvaluationService 
    - ProxyService
    - PrintingService
    - AdminWebServer
    - RankingWebServer
    - Checker
    - ContestWebServer
    - Worker - requires language support
        - C C++ Java Pascal Python with zip executable PHP Rust C# 

## Version 0.3.0 alpha
- importing contests
    - importing contests in the italian filesystem format
- making k8s-cms scalable:
    - scaling `ContestWebServer` to cater to more participants

## Future Releases
- securing k8s-cms:
    - data storage encryption
    - k8s communication encryption.
    - HTTPs for RankingWebServer,AdminWebServer,ContestWebServer.

- multiple contests support
    - contests can be obtained from DB via `get_contest_list()`
    - make `cms-proxy` run without an active contest
    - spawn multiple `cms-proxy` to serve multiple contests

- making k8s-cms `Workers` scalable:
    - scaling `Workers` to cater to more participants.
    - regenerate `cms.conf` using kubernetes deployment/docker-compose file.
    - restart `Checker` and `EvaluationService` to load rescaled workers
    - lightweight version of `cms-worker` with limited language support.
