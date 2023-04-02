.PHONY: build clean

IMAGE_NAME ?= ghcr.io/vaayne/notionboy
IMAGE_TAG ?= latest

build:
	go build -o ./bin/notionboy ./cmd/notionboy/main.go

init:
	if ! which pre-commit > /dev/null; then pip install pre-commit; fi
	pre-commit install

static: init
	golangci-lint run --fix --enable gofumpt
	pre-commit run --all-files

run:
	go run ./cmd/notionboy/main.go

local:
	cd .local && go run ./main.go

dev:
	cd .local && ENV=dev go run ./main.go

devui:
	cd webui && NODE_ENV="development" pnpm dev

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
	# set base path, so we don't need change it anymore
	# find webui/dist -type f -exec sed -i 's/\/_next\/static\//\.\/_next\/static\//g' {} +

grpc:
	rm -rf api/pb api/docs webui/src/lib.pb
	buf format -w
	buf generate

mdiff:
	atlas migrate diff init \
		--dir "file://db/ent/migrate/migrations" \
		--to "ent://db/ent/schema" \
		--dev-url "docker://mysql/8/ent"

mrehash:
	atlas migrate hash --dir "file://db/ent/migrate/migrations"

mapply:
	atlas migrate apply \
		--dir "file://db/ent/migrate/migrations" \
		--url "mysql://$(DB_NAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)" --baseline 20230402035028

mstatus:
	atlas migrate status \
		--dir "file://db/ent/migrate/migrations" \
		--url "mysql://$(DB_NAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)"

mlint:
	atlas migrate lint \
		--dir "file://db/ent/migrate/migrations" \
		--dev-url "mysql://$(DB_NAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)" \
		--latest 1
