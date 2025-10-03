set dotenv-load := true
set dotenv-required := true

postgres_dsn := "postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable"

[group('database migration')]
new-migration name:
    migrate create -ext sql -dir migrations -seq {{name}}

[group('database migration')]
migrate-up step:
    migrate -database {{postgres_dsn}} -path migrations up {{step}}

[group('database migration')]
migrate-down step:
    migrate -database {{postgres_dsn}} -path migrations down {{step}}

[group('database migration')]
migration-version:
    migrate -database {{postgres_dsn}} -path migrations version

[group('database migration')]
set-migration-version version:
    migrate -database {{postgres_dsn}} -path migrations force {{version}}

[group('backend')]
dev-application:
    cd backend && air -c .air.application.toml

dev-worker:
    cd backend && air -c .air.worker.toml