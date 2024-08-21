VERSION := $(shell git describe --tags | sed -e 's/^v//g' | awk -F "-" '{print $$1}')

SHELL := /bin/bash

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

LOCAL_BIN:=$(CURDIR)/bin
GOLANGCI_TAG:=1.58.0
PATH:=$(PATH):$(LOCAL_BIN)

default: help

.PHONY: help
help: # Показывает информацию о каждом рецепте в Makefile
	@grep -E '^[a-zA-Z0-9 _-]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

.PHONY: .bin_deps
.bin_deps: # Устанавливает зависимости необходимые для работы приложения
	$(info Installing binary dependencies...)
	mkdir -p $(LOCAL_BIN)
	GOPROXY=direct GOBIN=$(LOCAL_BIN) go install github.com/vektra/mockery/v2@v2.42.0
	GOBIN=$(LOCAL_BIN) go install github.com/rubenv/sql-migrate/...@v1.5.2
	GOPROXY=direct GOBIN=$(LOCAL_BIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: .app_deps
.app_deps: # Устанавливает необходимые go пакеты
	GOPROXY=$(GOPROXY) GOPRIVATE=$(GOPRIVATE) go mod tidy

.PHONY: go_get
go_get: # Хелпер для go get
	GOPROXY=$(GOPROXY) GOPRIVATE=$(GOPRIVATE) go get $(RUN_ARGS)

.PHONY: .install_linter
.install_linter: # Устанавливает линтер
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info Downloading golangci-lint v$(GOLANGCI_TAG))
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_TAG)
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
endif

.PHONY: install
install: .install_linter .bin_deps .app_deps # Устанавливает все зависимости для работы приложения

.PHONY: tests
tests: # Запускает юнит тесты с ковереджем
	go test -race -coverprofile=coverage.out ./...

.PHONY: tests_integration
tests_integration: # Запуск интеграционных тестов (пока только локально)
	go test -race -tags=integration -coverprofile=coverage.out ./...

.PHONY: show_cover
show_cover: # Открывает ковередж юнит тестов
	go tool cover -html=coverage.out

.PHONY: linter
linter: # Запуск линтеров
	$(LOCAL_BIN)/golangci-lint cache clean && \
	$(LOCAL_BIN)/golangci-lint run

.PHONY: linter_fix
linter_fix: # Запуск линтеров с фиксом где возможно
	$(LOCAL_BIN)/golangci-lint cache clean && \
	$(LOCAL_BIN)/golangci-lint run --fix

.PHONY: quality
quality: linter tests_integration # Запуск линтеров и интеграционных тестов

.PHONY: clean_cache
clean_cache: # Очистка кеша go
	go clean -cache

binary_name=FamilySyncHub
.PHONY: build_binary
build_binary: # Компиляция проекта
	go build -ldflags="-X github.com/Yaroher2442/FamilySyncHub/build.Version=$(VERSION)" -v -o $(binary_name) ./cmd

.PHONY: build_docker
build_docker: # Сборка докер образа
	docker build --build-arg VERSION=$(VERSION) -t ptaf-integrations-manager:$(VERSION) -f docker/Dockerfile .

.PHONY: mockery
mockery: # Создание моков
	$(LOCAL_BIN)/mockery --name $(name) --dir $(dir) --output $(dir)/mocks


# Add mocking interface as make mockery name=Interface dir=./path/to/interface/dir
.PHONY: mock
mock:

.PHONY: generate_go_sql
generate_go_sql: # Генерация SQL структуры из миграций
	$(LOCAL_BIN)/sqlc generate

.PHONY: app_run
app_run: # Запуск приложения
	./FamilySyncHub


SQL_MIGRATE_CONFIG:=./sql-migrate.yml
SQL_MIGRATE_ENV:=db
.sql-migrate = $(LOCAL_BIN)/sql-migrate $1 -config=$(SQL_MIGRATE_CONFIG) -env=$(SQL_MIGRATE_ENV) $2

# sql migrations

.PHONY: migrate_status
migrate_status: # Статус миграций
	$(call .sql-migrate,status)

# make migrate_new NEW_MIGRATION_NAME
.PHONY: migrate_new
migrate_new: # Создание миграции
	$(call .sql-migrate,new,$(RUN_ARGS))

# make migrate_up MIGRATION_ID
# if MIGRATION_ID is empty -> migrate to the latest version
.PHONY: migrate_up
migrate_up: # Применение миграции
	$(call .sql-migrate,up,$(RUN_ARGS))

# make migrate_down MIGRATION_ID
# if MIGRATION_ID is empty -> migrate to the first version
.PHONY: migrate_down
migrate_down: # Откат миграции
	$(call .sql-migrate,down,$(RUN_ARGS))

