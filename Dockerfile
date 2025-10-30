# alpine versions: https://alpinelinux.org/downloads/
ARG ALPINE_VERSION="3.22"
# golang versions https://go.dev/dl/
ARG GO_VERSION="1.25.3"

# specifies a parent image (image is alpine + all the stuff you need to build a golang application)
# and names this instance 'building-image'.
# cryptic source image names like 'alpine' explained in https://stackoverflow.com/a/59731596/11593686.
# official docker images for golang: https://hub.docker.com/_/golang/
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS building-image

# install git so that go will be able to write down current version when building
RUN apk add --no-cache git

# mkdir+cd into new directory, we are going to put everything there
WORKDIR /myapp

# copy entire project there (except for what's listed in .dockerignore)
COPY . .

# install all dependencies
RUN go mod download

# run go build, name the executable "go-api" and also disable CGO because people keep telling me that
RUN CGO_ENABLED=0 go build -o go-api

# switch to a new clean alpine without the golang stuff, much smaller
# General article about so called multi-stage patterns: https://medium.com/swlh/reducing-container-image-size-esp-for-go-applications-db7658e9063a
FROM alpine:${ALPINE_VERSION} AS running-image

# ensure sqlite is available on running-image
RUN apk add --no-cache sqlite

# install mailcap to add mime type support, https://stackoverflow.com/a/38033047
RUN apk add --no-cache mailcap

# notify docker we are going to be using port 8080
EXPOSE 8080

# cd to the folder again
# if this is not done, relative paths inside of the app code will start from root, which is not what we want
WORKDIR /myapp

 # Declare a volume for persistent data
VOLUME ["/myapp/data"]

# tell docker what to run
ENTRYPOINT ["/myapp/go-api",  "--config", "/myapp/data/config.json"]

# copy the executable into the running image, everything should be embedded
COPY --from=building-image /myapp/go-api /myapp/go-api
