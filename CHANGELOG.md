# k8s-cms Changelog
## Version 0.2.1 beta - 2019-10-04
### Changed
- securing k8s-cms:
	- expose only required secrets instead of using envfrom exposing everything
    - run cms pods with cms-services service account, with automount token disable.

## Version 0.2.0 alpha - 2019-09-30
### Added
- database dependency - wait for database before starting CMS services
- kubernetes support - write k8s deployment for YAMLs all these:
    - Database :heavy_check_mark:
    - ResourceService :heavy_check_mark:
    - LogService :heavy_check_mark:
    - ScoringService :heavy_check_mark:
    - EvaluationService :heavy_check_mark:
    - ProxyService :heavy_check_mark:
    - AdminWebServer :heavy_check_mark:
    - PrintingService :heavy_check_mark:
    - RankingWebServer :heavy_check_mark:
    - Checker :heavy_check_mark:
    - ContestWebServer :heavy_check_mark:
    - Worker  :heavy_check_mark:
    - Ingress service to combine ranking, contest and admin web servers.

### Changed
- make cms docker images source cms configuration from /etc/
- split singular env file to config.env for config, .env for secrets

- securing k8s-cms:
    - run all (except db) services as non root user.

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
