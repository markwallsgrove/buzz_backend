# Muzz Backend Challenge
Mini API that could power a very simple dating app.

# Demonstrated qualities
- Experience
- Diversity (dev, platform, security, sre)
- Quality/robust code (normally aim full test coverage but ran out of time for this project)
- Testing (unit, integration, browser)
- Startup experience (currently at my third startup)
- Experience with many languages (Node, Java, Python, Golang, and even Assembly)
- Docker, Docker Compose, Kubernetes (EKS), and ECS Fargate experience
- Focus on quality of implementation along with cost of running it
- New Relic advocate (ten years experience)
- Always finds root cause of the issue
- Micro service environment experience / debugging

# Development
## Initialisation
To bring up the environment only `docker-compose` is required:
```sh
docker-compose up -d
```
Endpoints:
- webserver: `http://localhost:80`
- database: `localhost:3306`
## Required Tooling
- golang 1.18+
- [mockery](https://github.com/vektra/mockery)
- docker
- docker compose
## Tests
To run the unit tests run the following command:
```sh
docker-compose run unittests
```
The integration tests can also be run by executing:
```sh
docker-compose run integrations
```
`docker-compose` isn't mandatory to use. To run the commands by hand see the [docker-compose.yml](./docker-compose.yml) for more details on how to run the tests.
