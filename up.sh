 #!/usr/bin/env sh
 
git pull

export GIT_COMMIT_SHA="$(git rev-parse --short HEAD)"
docker compose up --build --detach