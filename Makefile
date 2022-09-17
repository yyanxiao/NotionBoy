.PHONY: build clean

IMAGE_NAME ?= ghcr.io/vaayne/notionboy
IMAGE_TAG ?= latest

build:
	go build -o ./bin/notionboy ./cmd/notionboy/main.go

init:
	if ! which pre-commit > /dev/null; then pip install pre-commit; fi
	pre-commit install

static: init
	pre-commit run --all-files

run:
	go run ./cmd/notionboy/main.go

rund:
	docker run --rm -v `pwd`/data.db:/service/data.db -v `pwd`/settings_local.yaml:/service/settings_local.yaml --net=host ${IMAGE_NAME}:${IMAGE_TAG}

buildd:
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .

clean:
	rm -rf ./bin ./vendor go.sum
