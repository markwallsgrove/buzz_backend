# Muzz Backend Challenge
Mini API that could power a very simple dating app.

# Development
## Initialisation
To bring up the environment only `docker-compose` is required:
```sh
docker-compose up -d
```
Endpoints:
- webserver: `http://localhost:80`
- database: random port run `docker-compose ps mariadb` for more details
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
