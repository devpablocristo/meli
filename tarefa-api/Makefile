SHELL:=/bin/bash -O extglob
BINARY=taregaapi  # Cambiar el nombre del binario si es necesario
VERSION=0.1
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# Directorio raíz del proyecto
ROOT_DIR := $(shell pwd)

# Ruta del archivo .env en el directorio raíz
ENV_FILE := $(ROOT_DIR)/.env

# go
build-bin:
	go build ${LDFLAGS} -o ${BINARY} cmd/api/main.go

run:
	@go run cmd/api/main.go

test:
	@go test ./...

# docker-compose
build-con:
	docker-compose up --build

up:
	docker-compose up

down:
	docker-compose down --remove-orphans