# TODO:

3. generate tokens and require those in query to talk
4. deploy on server
5. see if it works :)

### Library docs

- https://gorm.io/docs/
- https://golang-jwt.github.io/jwt/usage/create/

### Concept docs

- https://jwt.io/introduction

Generating secret for jwt tokens: `openssl rand -hex 32`

### Building

- Don't forget to start docker: `sudo systemctl start docker`
- run `build-and-log-size` in makefile