# k8s-cms Roadmap
## Version 0.3.0
- migrate kubernetes yaml to helm v2 
- use ingress chart to setup nginx ingress
- migrate postgres db service from in house manifest to postgres chart
- securing k8s-cms:
    - setup cert manager chart to use lets encrypt to obtain certificates 
    - HTTPs for RankingWebServer,AdminWebServer,ContestWebServer.
        - set up lets encryptontainer to perform acme challenge for certificate
- Quality of Service
    - limit worker resources to ensure quality of service.
    - set limits and requests for pods in kubernetes
		- requires profiling of container performance. 
	- ddos protection - conncurrent connections and rate limiting through nginx ingress.

## Version 0.4.0
- k8s-cms-master 
    - bridge between CMS and k8s
    - exposes REST api used to control cms with `kcmscli` 
    - updates configuration and restarts containers to reload configuration when scaling.
- importing contests
    - importing contests in the italian filesystem format with `kcmscli`
    - loads k8s-cms.yml which contains most options
- making k8s-cms scalable:
    - auto scaling `ContestWebServer` to cater to more participants
- making k8s-cms `Workers` scalable:
    - scaling `Workers` to cater to more participants.
    - regenerate `cms.conf` using kubernetes deployment/docker-compose file.
    - restart `Checker` and `EvaluationService` to load rescaled workers
    - lightweight version of `cms-worker` with limited language support.

- multiple contests support
    - contests can be obtained from DB via `get_contest_list()`
    - spawn multiple `cms-proxy` to serve multiple contests

- printing support.

