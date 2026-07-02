# Trabalho de Grafos - 2026
# Atalhos para compilar e rodar as tres unidades.
# Uso: `make`, `make u1`, `make u2`, `make u3`, `make build`, etc.

BIN_DIR := bin
U1_SRC  := ./cmd/unidade1
U2_SRC  := ./cmd/unidade2
U3_SRC  := ./cmd/unidade3
U1_BIN  := $(BIN_DIR)/unidade1
U2_BIN  := $(BIN_DIR)/unidade2
U3_BIN  := $(BIN_DIR)/unidade3

.DEFAULT_GOAL := help

## help: mostra esta ajuda
.PHONY: help
help:
	@echo "Trabalho de Grafos - 2026"
	@echo ""
	@echo "Alvos disponiveis:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'

## u1: roda a Unidade 1 (sem compilar) -> outputs/
.PHONY: u1
u1:
	go run $(U1_SRC)

## u2: roda a Unidade 2 - Eulerianos (sem compilar) -> outputs_u2/
.PHONY: u2
u2:
	go run $(U2_SRC)

## u3: roda a Unidade 3 - Caixeiro Viajante (sem compilar) -> outputs_u3/
.PHONY: u3
u3:
	go run $(U3_SRC)

## web: sobe a aplicacao web do Projeto 3 (EcoUrbano) em http://localhost:8080
.PHONY: web
web:
	go run ./cmd/servidor

## run: roda as tres unidades em sequencia
.PHONY: run
run: u1 u2 u3

## build: compila as tres unidades em bin/
.PHONY: build
build: build-u1 build-u2 build-u3

## build-u1: compila a Unidade 1 em bin/unidade1
.PHONY: build-u1
build-u1:
	@mkdir -p $(BIN_DIR)
	go build -o $(U1_BIN) $(U1_SRC)

## build-u2: compila a Unidade 2 em bin/unidade2
.PHONY: build-u2
build-u2:
	@mkdir -p $(BIN_DIR)
	go build -o $(U2_BIN) $(U2_SRC)

## build-u3: compila a Unidade 3 em bin/unidade3
.PHONY: build-u3
build-u3:
	@mkdir -p $(BIN_DIR)
	go build -o $(U3_BIN) $(U3_SRC)

## vet: roda go vet em todo o modulo
.PHONY: vet
vet:
	go vet ./...

## fmt: formata todo o codigo com gofmt
.PHONY: fmt
fmt:
	gofmt -w .

## check: build de tudo + vet (sanidade antes de commitar)
.PHONY: check
check:
	go build ./...
	go vet ./...

## clean: remove binarios e saidas geradas
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
	rm -f outputs/*.dot outputs_u2/*.dot outputs_u3/*.dot
