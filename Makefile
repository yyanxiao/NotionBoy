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

local:
	cd .local && go run ./main.go

rund:
	docker run --rm -v `pwd`/data.db:/service/data.db -v `pwd`/settings_local.yaml:/service/settings_local.yaml --net=host ${IMAGE_NAME}:${IMAGE_TAG}

buildd:
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .

clean:
	rm -rf ./bin ./vendor

ent:
	go generate ./db/ent/

webui: clean
	cd webui && npm run export
	find webui/dist -type f -exec sed -i 's/\/_next\/static\//\.\/_next\/static\//g' {} +

grpc:
	rm -rf api/pb api/docs webui/src/lib.pb
	buf format -w
	buf generate
