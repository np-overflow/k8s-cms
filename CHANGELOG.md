# k8s-cms Changelog

## Version 0.2.0 alpha - Unreleased
### Added
- database dependency - wait for database before starting
- kubernetes support - write k8s deployment for YAMLs all these:
    - Database :heavy_check_marK:

### Changed
- make cms docker images source cms configuration from /etc/
- split singular env file to config.env for config, .env for secrets

## Version 0.1.0 alpha - 2019-09-21
### Added
- docker containers for CMS services :heavy_check_mark:
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
