include .env
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

test:
	docker compose -f docker-compose-test.yml -p testing up -d
	sleep 2
	-migrate -database "$(POSTGRES_URL_TEST)" -path migrations up
	-POSTGRES_URL_TEST=$(POSTGRES_URL_TEST) go test -p 1 ./... -cover -v -count 1
	docker rm -f test
