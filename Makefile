.PHONY: help migrate server build run-prod clean seed

# Comando padrão
help:
	@echo "Comandos disponíveis:"
	@echo "  make migrate    - Executa migrations do banco de dados"
	@echo "  make server     - Inicia o servidor HTTP"
	@echo "  make seed       - Popula o banco de dados com dados iniciais (somente desenvolvimento)"
	@echo "  make build      - Compila o binário para produção"
	@echo "  make run-prod   - Executa migrate + server em produção"
	@echo "  make clean      - Remove binários compilados"

# Desenvolvimento
migrate:
	@echo "Executando migrations..."
	go run ./cmd/api migrate

server:
	@echo "Iniciando servidor..."
	go run ./cmd/api server

seed:
	@echo "Populando banco de dados com dados iniciais..."
	go run ./cmd/api seed

# Produção
build:
	@echo "Compilando binário..."
	go build -o bin/api ./cmd/api
	@echo "Binário criado em: bin/api"

run-prod: build
	@echo "Executando migrations..."
	./bin/api migrate
	@echo "Iniciando servidor..."
	./bin/api server

clean:
	@echo "Limpando binários..."
	rm -rf bin/
	@echo "Limpeza concluída"