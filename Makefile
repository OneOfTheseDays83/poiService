#!/usr/bin/make

.PHONY: help
.SECONDEXPANSION:


################################################################################
# Set common variables
PROJECT                             := poi-service
BUILD_OUTPUT_DIR                    ?= dist
SERVICE_PORT						:= 8000
DATABASE_URL						:= mongodb://localhost:27017
DEST_HYDRA							:= /tmp/hydra/

################################################################################
help: ## Print this help message.
	@echo "Usable make targets:"
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1:\2/' | column -c2 -t -s : | sort)"

################################################################################
# Build, Package, Test and Code Quality Make Targets

download-deps:
	go mod download -xl
	docker pull mongo

build:
	docker build \
		--network=host \
		-f docker/build.Dockerfile \
		-t "$(PROJECT)" \
		.
build-local:
	cd ./cmd; go build -o ../dist/poiService

test-local:
	cd ./cmd; go test ./... -coverprofile ../dist/coverage.out


auth-server-start:
	cd /tmp/hydra
	docker-compose -f /tmp/hydra/quickstart.yml \
               -f /tmp/hydra/quickstart-postgres.yml \
               -f /tmp/hydra/quickstart-jwt.yml \
               up --build

auth-server-create-client:
	docker-compose -f /tmp/hydra/quickstart.yml exec hydra \
			   hydra clients create \
			   --endpoint http://127.0.0.1:4445/ \
			   --id my-client \
			   --secret secret \
			   -g client_credentials

auth-server-download:
	if ! git clone "https://github.com/ory/hydra.git" $(DEST_HYDRA) 2>/dev/null ; then echo "Hydra already cloned"; fi

mongodb:
	docker run  --rm --name mongo-db -p 27017:27017 -d mongo:latest

start-environment:
start-environment: mongodb
start-environment: auth-server-download
start-environment: auth-server-start
start-environment: auth-server-create-client

stop-environment:
	docker stop mongo-db
	docker stop hydra_hydra_1
	docker stop hydra_consent_1
	docker stop hydra_postgresd_1

start:
	docker run -p $(SERVICE_PORT):$(SERVICE_PORT) --env SERVICE_PORT=$(SERVICE_PORT) --env DATABASE_URL=$(DATABASE_URL) --network=host $(PROJECT)

start-local:
	DATABASE_URL=$(DATABASE_URL) SERVICE_PORT=$(SERVICE_PORT) ./dist/poiService

gen-mocks:
	mockgen -destination=cmd/handler/mongo_mock.go -package="handler" -source=cmd/handler/mongo.go
	mockgen -destination=cmd/handler/pois_mock.go -package="handler" -source=cmd/handler/pois.go
	mockgen -destination=cmd/auth/keyStore_mock.go -package="auth" -source=cmd/auth/keyStore.go
	mockgen -destination=cmd/download/http_mock.go -package="download" -source=cmd/download/http.go