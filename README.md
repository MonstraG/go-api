### Library docs

- https://gorm.io/docs/
- https://golang-jwt.github.io/jwt/usage/create/
- https://github.com/ytdl-org/youtube-dl

### Concept docs

- https://jwt.io/introduction

Generating secret for jwt tokens: `openssl rand -hex 32`

### Building

- Don't forget to start docker: `sudo systemctl start docker`
- run `build-and-log-size` in makefile

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