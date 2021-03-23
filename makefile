.PHONY: dev test prd help
.DEFAULT_GOAL: help

default: help

help: ## Output available commands
	@echo "Available commands:"
	@echo
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

multi:## Build and push a project docker image for multiple platforms
	@docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t waduphaitian/mural_dev:multi --push .

build:## Build a docker container and tag as mvral
	@docker build --build-arg BUILD_VERSION=1.0.0 -t waduphaitian/mural_dev:latest .

run:## Build a docker container and tag as mvral
# docker run -v <Path to host dir>:/containerFiles -p <host port to use>:42069 -it mvral
	@docker run -v ${PWD}/containerFiles:/containerFiles -p 42069:42069 -it waduphaitian/mural_dev:latest .


tar:## Build a docker container and tag as mvral
	@docker save waduphaitian/mural_dev:latest | gzip > mural_dev.tar.gz

buildTar:## Build a docker container and tag as mvral
	@make build
	@make tar

buildArm:## Build a go executable to arm64 computer
	@env GOOS=linux GOARCH=arm64 go build -ldflags "-X main.BuildVersion=1.0.0"
	@tar czf muraldevice.tar.gz muraldevice