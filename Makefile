.PHONY: build clean

IMAGE_NAME ?= ghcr.io/vaayne/notionboy
IMAGE_TAG ?= latest

build: build-wxgzh build-netlify

init:
	if ! which pre-commit > /dev/null; then sudo pip install pre-commit; fi
	pre-commit install

static: init
	pre-commit run --all-files

run:
	go run ./main.go

rund: build-docker
	docker run --rm ${IMAGE_NAME}:${IMAGE_TAG}

build-netlify:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/notion cmd/notion/main.go

build-wxgzh:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/wxgzh cmd/wxgzh/main.go

build-docker:
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .

clean:
	rm -rf ./bin ./vendor go.sum
