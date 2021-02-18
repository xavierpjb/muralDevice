.PHONY: dev test prd help
.DEFAULT_GOAL: help

default: help

help: ## Output available commands
	@echo "Available commands:"
	@echo
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

multi:## Build and push a project docker image for multiple platforms
	@docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t waduphaitian/mural_dev:multi --push .