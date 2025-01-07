run:
	docker compose up

down:
	docker compose down

lint:
	golangci-lint run ./... -v