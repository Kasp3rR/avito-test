MIGRATION_FOLDER = $(CURDIR)/internal/db/migrations
POSTGRES_SETUP_TEST := user=avito_test password=test dbname=avito_test host=localhost port=5432 sslmode=disable

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down -v

migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down