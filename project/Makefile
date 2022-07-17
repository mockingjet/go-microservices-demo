AUTH_BINARY=authApp
BROKER_BINARY=brokerApp
LOGGER_BINARY=loggerApp
MAILER_BINARY=mailerApp
LISTENER_BINARY=listenerApp
FRONT_END_BINARY=frontApp

define run
	docker run -it --rm \
	-v $(CURDIR)/../$(1)/:/app \
 	-w /app \
	-e CGO_ENABLED=0 \
	--name go-builder \
	golang:1.18-alpine \
 	$(2)
endef

## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images...
	docker-compose up -d
	@echo Docker images started!

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_auth build_broker build_logger build_mailer build_listener build_front
	@echo Stopping docker images (if running...)
	docker-compose down
	@echo Building (when required) and starting docker images...
	docker-compose up --build -d
	@echo Docker images built and started!

## down: stop docker compose
down:
	@echo Stopping docker compose...
	docker-compose down
	@echo Done!

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo Building auth binary...
	$(call run,auth-service,go build -o ${AUTH_BINARY} ./cmd/api)
	@echo Done!

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo Building broker binary...
	$(call run,broker-service,go build -o ${BROKER_BINARY} ./cmd/api)
	@echo Done!

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo Building logger binary...
	$(call run,logger-service,go build -o ${LOGGER_BINARY} ./cmd/api)
	@echo Done!

## build_mail: builds the mail binary as a linux executable
build_mailer:
	@echo Building mailer binary...
	$(call run,mailer-service,go build -o ${MAILER_BINARY} ./cmd/api)
	@echo Done!

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo Building listener binary...
	$(call run,listener-service,go build -o ${LISTENER_BINARY} ./)
	@echo Done!


## build_front: builds the frone end binary
build_front:
	@echo Building front end binary...
	$(call run,front-end,go build -o ${FRONT_END_BINARY} ./cmd/web)
	@echo Done!