# Decolecta Sunat Ruc

```sh
cd ~/deco/decolecta-ruc
export DATABASE_URI="postgres://juan:createdat2024@localhost:5432/decolecta_rucs?sslmode=disable"
export REDIS_URI="localhost:6379"
go run .
```

## Tests

Run the happy-path suite with race detection and coverage:

```sh
go test -race -cover ./...
```

Tests cover all registered HTTP routes, cached and database-backed company
lookups, the full V2 JSON contract, address parsing, response conversions, and
RUC utilities. Redis runs in memory and PostgreSQL queries use mocks, so no
external services or environment variables are required. The suite does not
validate connectivity or schema compatibility against a live database.
