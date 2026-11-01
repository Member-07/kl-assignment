run:
	go run cmd/api/main.go

ิีbuild:
	docker compose up -d

migrate:
	docker compose up -d migrate

k6:
	docker compose up -d k6