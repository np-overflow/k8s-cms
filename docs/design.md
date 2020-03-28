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

