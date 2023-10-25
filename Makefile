ENV ?= development
-include .env
-include .env.$(ENV)
export

build:
	go build -C cmd -o ../bin/main

run:build
	./bin/main

migrate-up:
	migrate -database "$(POSTGRES_URL)" -path migrations up

migrate-down:
	migrate -database "$(POSTGRES_URL)" -path migrations down

migrate-force:
	migrate -database "$(POSTGRES_URL)" -path migrations force 1

migrate-drop:
	migrate -database "$(POSTGRES_URL)" -path migrations drop

connect-db:
	psql -d "host=localhost port=$(POSTGRES_PORT) password=$(POSTGRES_PASSWORD) user=$(POSTGRES_USER)"

test:
	docker compose -p testing up -d
	sleep 2
	-migrate -database "$(POSTGRES_URL)" -path migrations up
	-POSTGRES_URL=$(POSTGRES_URL) go test ./... -p 1 -cover -v -count 1
	docker rm -f test
