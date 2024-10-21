MIGRATE_CMD = go run ./cmd/migrate
MIGRATIONS_DIR=./db/migrations

build:
	@go build -o ./tmp/psionicalch ./cmd/web

run: build
	./tmp/psionicalch

migrate-up:
	@$(MIGRATE_CMD) up

migrate-down:
	@$(MIGRATE_CMD) down

rollback:
	@rm ./db/db.*

new-migration:
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

new-keys:
	@go run ./cmd/keys

clean:
	rm -rf ./tmp
