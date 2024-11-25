MIGRATE_CMD = go run ./cmd/migrate
MIGRATIONS_DIR=./db/migrations

build:
	@go build -ldflags="-s -w" -o ./tmp/psionicalch ./cmd/web

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

load-content:
	@go run -tags content_loader ./cmd/content

generate-file-key:
	@go run ./cmd/filekey

new-admin:
	@go run ./cmd/add_admin_user -name="$(name)" -surname="$(surname)" -email="$(email)" -password="$(password)"

clean:
	rm -rf ./tmp
