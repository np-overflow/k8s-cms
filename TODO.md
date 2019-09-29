# k8s-cms Roadmap

## Version 0.2.0 alpha

## Version 0.3.0 alpha
- cms-master 
    - master service manages current cms status
    - exposes REST api used to control cms with `kcmscli` 
    - updates configuration and restarts containers to reload configuration when scaling.
- importing contests
    - importing contests in the italian filesystem format with `kcmscli`
- making k8s-cms scalable:
    - scaling `ContestWebServer` to cater to more participants
- making k8s-cms `Workers` scalable:
    - scaling `Workers` to cater to more participants.
    - regenerate `cms.conf` using kubernetes deployment/docker-compose file.
    - restart `Checker` and `EvaluationService` to load rescaled workers
    - lightweight version of `cms-worker` with limited language support.

- securing k8s-cms:
    - HTTPs for RankingWebServer,AdminWebServer,ContestWebServer.
    - read only filesystem for docker images


## Future Releases
- multiple contests support
    - contests can be obtained from DB via `get_contest_list()`
    - make `cms-proxy` run without an active contest
    - spawn multiple `cms-proxy` to serve multiple contests

- printing support.
