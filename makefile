ifneq (,$(wildcard ./.env))
    include .env
    export
endif


BINARY=engine
test: clean documents generate
    go test -v -cover -covermode=atomic ./...

coverage: clean documents generate
	bash coverage.sh --html

dev: generate
	go run github.com/air-verse/air

run: generate
	go run .

build:
	GOFLAGS=-buildvcs=false go build -o $123BINARY125 .

clean:
	@if [ -f $123BINARY125 ] ; then rm $123BINARY125 ; fi
	@find . -name *mock* -delete
	@rm -rf .cover wire_gen.go docs

docker_build:
	docker build -t evv-service  .

docker_start:
	docker-compose up --build

docker_stop:
	docker-compose down

lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

generate:
	GOFLAGS=-buildvcs=false go generate ./...

# SQLite Migrations
migrate-up:
	migrate -path migrations -database "sqlite://${DB_SQLITE_PATH}" up

create-migration:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -format unix $$name

migrate-down:
	@read -p "Enter number of migrations to revert (default 1): " num; \
	migrate -path migrations -database "sqlite3://${DB_SQLITE_PATH}" down $${num:-1}

migrate-fix:
	@read -p "Enter migration version: " version; \
	migrate -path migrations -database "sqlite3://${DB_SQLITE_PATH}" force $$version

# Legacy PostgreSQL migrations (commented out)
# migrate-up-pg:
# 	migrate -path migrations -database "postgres://${DB_POSTGRES_WRITE_USER}:${DB_POSTGRES_WRITE_PASSWORD}@${DB_POSTGRES_WRITE_HOST}:${DB_POSTGRES_WRITE_PORT}/${DB_POSTGRES_WRITE_NAME}?sslmode=disable" up
#
# migrate-down-pg:
# 	@read -p "Enter number of migrations to revert (default 1): " num; \
# 	migrate -path migrations -database "postgres://${DB_POSTGRES_WRITE_USER}:${DB_POSTGRES_WRITE_PASSWORD}@${DB_POSTGRES_WRITE_HOST}:${DB_POSTGRES_WRITE_PORT}/${DB_POSTGRES_WRITE_NAME}?sslmode=disable" down $${num:-1}
#
# migrate-fix-pg:
# 	@read -p "Enter migration version: " version; \
# 	migrate -path migrations -database "postgres://${DB_POSTGRES_WRITE_USER}:${DB_POSTGRES_WRITE_PASSWORD}@${DB_POSTGRES_WRITE_HOST}:${DB_POSTGRES_WRITE_PORT}/${DB_POSTGRES_WRITE_NAME}?sslmode=disable" force $$version

format:
	gofmt -s -w .

.PHONY: test coverage engine clean build docker run stop lint-prepare lint documents generate