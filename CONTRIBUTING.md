## Contributing
Guidelines for contributors:
- Contribute to the repository by creating a pull request.
- Make suggestions for improvement by creating an issue.
- Be nice.


### Development Environment Setup
Development setup for contributors:
1. Resolve submodules after cloning
```sh
git submodule update --init --recursive
```
2. Fill `env` file with secrets
```sh
cp env .env
nano .env # use your favourite editor
```
3. Run the stack
```sh
# use docker-compose OR
docker-compose up 
# use kubernetes. you need to have a kubernetes cluster
# already up and running.
helm install k8s-cms .
```
### Project Structure
- project changelog - updated on new releases: `CHANGELOG.md`
- dockerfiles used to build the images: `containers/`
- Kubernetes manifests/helm chart: `chart/`
- `kcmscli` source code: `src/kcmscli/`
- End to end tests: `test/`
