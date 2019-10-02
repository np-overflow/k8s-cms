# k8s-cms Roadmap
## Version 0.2.1 alpha
### Added
- securing k8s-cms:
	- expose only required secrets instead of using env from exposing everything

### Todo
- securing k8s-cms:
    - tighten cluster security using RBAC
        - with service accounts and role bindings
        - assign read only for secrets permissions to service accounts
    - HTTPs for RankingWebServer,AdminWebServer,ContestWebServer.
        - set up lets encrypt container to perform acme challenge for certificate
	- read only filesystem
- Quality of Service
    - limit worker resources to ensure quality of service.
    - set limits and requests for pods in kubernetes
		- requires profiling of container performance. 
	- ddos protection - conncurrent connections and rate limiting through nginx ingress.

## Version 0.3.0 alpha
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

## Future Releases
- multiple contests support
    - contests can be obtained from DB via `get_contest_list()`
    - spawn multiple `cms-proxy` to serve multiple contests

- printing support.
