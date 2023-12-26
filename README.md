# CRUD quotes app

Simple CRUD app made with `go`, `chi` and `pgx`.

## Endpoints
- `GET localhost:8080` - health check
- `POST localhost:8080/quotes` - Create new quote
- `GET localhost:8080/quotes` - Read all quotes
- `PUT localhost:8080/quotes/{id}` - Update quote
- `DELETE localhost:8080/quotes/{id}` - Delete quote

## Installation
### Prerequisites
- [go](https://go.dev/dl/)
- [just](https://github.com/casey/just) - command runner
- docker
- `sqlx-cli` - `cargo install sqlx-cli`

### `.env` file
Create `.env` file containing:
``` sh
DATABASE_ADDRESS=localhost:5432
POSTGRES_PASSWORD=<YOUR PASSWORD>
POSTGRES_DB=quotes
DATABASE_URL=postgresql://postgres:<YOUR PASSWORD>@localhost:5432/quotes
```

### Migrating database
To setup the database, simply run
``` sh
docker-compose up -d
just migrate
```

## Running server
`just run`
