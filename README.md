### Library docs

- https://gorm.io/docs/
- https://golang-jwt.github.io/jwt/usage/create/
- https://htmx.org/docs/

### Concept docs

- https://jwt.io/introduction

Generating secret for jwt tokens:

```shell
openssl rand -hex 32
```

### Building

- Don't forget to start docker:

```shell
sudo systemctl start docker
```

- run `build-and-get-size` from makefile
  (or just run `make`)

### Other notes

Default address is 0.0.0.0 not
localhost [because docker](https://serverfault.com/questions/1084915/still-confused-why-docker-works-when-you-make-a-process-listen-to-0-0-0-0-but-no).

When I, inevitably, would want to stop docker *container* and run the app straight:

```shell
docker container stop go-api
```

Will stop the server

```shell
docker system prune -a --volumes
```

Will delete all build artefacts from disk.

And, just in case, list containers:

```shell
docker container list
```

### Updating

1. Update locally installed go version
2. Update dependencies in go.mod
3. Update github action versions
4. Update alpine/go versions in Dockerfile
5. Update vendored htmx version