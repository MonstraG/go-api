### Library docs

- https://gorm.io/docs/
- https://htmx.org/docs/

### Concept docs

Generating secret for auth tokens:

```shell
openssl rand -hex 32
```

### Building

1. Remember to start docker:

```shell
sudo systemctl start docker
```

2. build the docker image

```shell
make
```

3. run the image

```shell
docker run --publish 8080:8080 --volume ./data:/myapp/data go-api
```

### Other notes

The default address is 0.0.0.0, not
localhost [because docker](https://serverfault.com/questions/1084915/still-confused-why-docker-works-when-you-make-a-process-listen-to-0-0-0-0-but-no).

When I, inevitably, would want to stop docker *container* and run the app straight:

```shell
docker container stop go-api
```

Will stop the server

```shell
docker system prune -a --volumes
```

Will delete all build artifacts from disk.

And, just in case, list containers:

```shell
docker container list --all
```

sh into a built image to inspect it

```shell
docker run -it --entrypoint sh go-api
```

### Updating

1. Update locally installed go version
2. Update dependencies in go.mod
3. Update alpine/go versions in Dockerfile
4. Update vendored htmx version
