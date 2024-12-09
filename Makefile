APP_NAME := myapp
IMAGE_NAME := $(APP_NAME):latest
SIZE_LOG_FILE := build_size.log

build:
	@echo "Building Docker image..."
	@docker build -t $(IMAGE_NAME) .

size:
	@echo "Calculating image size..."
	@IMAGE_SIZE=$$(docker image inspect $(IMAGE_NAME) --format='{{.Size}}') && \
	HUMAN_SIZE=$$(numfmt --to=iec --suffix=B --format="%.2f" $${IMAGE_SIZE}) && \
	echo "Build size: $${HUMAN_SIZE}" && \
	echo "[$$(date +'%Y-%m-%dT%H:%M:%S')] Build size: $${HUMAN_SIZE}" > $(SIZE_LOG_FILE)

build-and-log-size: build size
