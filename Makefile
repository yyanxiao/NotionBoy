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

# generate migration files
mFiles:
	@echo "create migration files, please input table name: "; \
	read tableName; \
	up_file=$$(printf "db/migrations/%s_%s.up.sql" "$$(date +%Y%m%d%H%M%S)" "$$tableName"); \
	down_file=$$(printf "db/migrations/%s_%s.down.sql" "$$(date +%Y%m%d%H%M%S)" "$$tableName"); \
	echo "" > $$up_file; \
	echo "" > $$down_file; \
	echo "Created $$up_file and $$down_file"

# migrate database
migrate:
	go run ./cmd/migratedb
