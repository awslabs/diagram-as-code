# ─── Platform detection ───────────────────────────────────────────────────────
ifeq ($(OS),Windows_NT)
  BINARY    := awsdac.exe
  RM        := del /Q
  RMDIR     := rmdir /S /Q
  SEP       := \\
  NULL      := nul
  OPEN      := start
  NPM       := npm.cmd
else
  BINARY    := awsdac
  RM        := rm -f
  RMDIR     := rm -rf
  SEP       := /
  NULL      := /dev/null
  OPEN      := open
  NPM       := npm
  UNAME_S   := $(shell uname -s)
  ifeq ($(UNAME_S),Linux)
    OPEN    := xdg-open
  endif
endif

WEB_DIR := web
FILE    ?= examples/alb-ec2.yaml

.PHONY: build test test-func run dev web install clean open help

## build: compila o CLI Go
build:
	go build -o $(BINARY) ./cmd/awsdac

## test: roda todos os testes (unit + functional)
test:
	go test ./...

## test-func: roda só os testes funcionais (gera PNGs em /tmp/results/)
test-func:
	go test ./test/...

## run: roda o CLI (uso: make run FILE=examples/alb-ec2.yaml)
run: build
	.$(SEP)$(BINARY) $(FILE)

## dev: sobe o frontend Next.js em modo desenvolvimento (porta 3001)
dev:
	cd $(WEB_DIR) && $(NPM) run dev

## web: instala dependências e sobe o frontend
web:
	cd $(WEB_DIR) && $(NPM) install && $(NPM) run dev

## install: instala dependências do frontend
install:
	cd $(WEB_DIR) && $(NPM) install

## open: abre o frontend no browser (requer 'make dev' rodando)
open:
	$(OPEN) http://localhost:3001

## clean: remove binário compilado e PNGs temporários de teste
clean:
	$(RM) $(BINARY) 2>$(NULL) || true
	$(RM) /tmp/results/*.png 2>$(NULL) || true

## help: lista todos os comandos disponíveis
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'
