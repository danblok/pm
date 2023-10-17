FILE=.env
POSTGRES_URL=`cat $(FILE) | grep POSTGRES_URL= | sed -e s/^POSTGRES_URL=//`

build:
	go build -C cmd -o ../bin/main

run:build
	go ./bin/main

migrate-up:
	migrate -database "$(POSTGRES_URL)" -path migrations up

migrate-down:
	migrate -database "$(POSTGRES_URL)" -path migrations down

migrate-force:
	migrate -database "$(POSTGRES_URL)" -path migrations force 1

migrate-drop:
	migrate -database "$(POSTGRES_URL)" -path migrations drop

