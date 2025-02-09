build:
	go build -o build/main ./cmd/api


install-deps:
	go mod tidy

infra-up:
	docker compose -p lumos-api -f ./dev/docker-compose.infra.yml up -d --wait

infra-down:
	docker compose -p lumos-api -f ./dev/docker-compose.infra.yml down

run:
	DOTENV_ENABLED=true DOTENV_CONFIG_PATH=./dev/.dev.env go run ./cmd/api

run-race:
	DOTENV_ENABLED=true DOTENV_CONFIG_PATH=./dev/.dev.env go run -race ./cmd/api