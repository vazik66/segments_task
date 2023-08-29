# Segments test task

## Installation

1. Set env vars
```sh
cat .env.example >> .env
```

<details>
<summary>Env defaults</summary>

```
PORT=8080
DB_PORT=5432
DB_PASSWORD=password
DB_USER=postgres
DB_DATABASE=postgres
DB_HOST=database
```
</details>

2. Launch

<details>
<summary>Docker</summary>

```sh
docker compose up -d
```

</details>

<details>
<summary>Locally</summary>

```sh
go mod download
go install github.com/swaggo/swag/cmd/swag@latest
```

Build docs
```sh
make docs
```

Run
```sh
make run
```
</details>


3. Migrate db
```sh
# using golang-migrate
make migrate
```

```sh
# using psql
psql -h localhost -U postgres -d postgres -a -f ./pkg/migrations/000001_init.up.sql
```
4. Swagger\
Go to http://localhost:8080/swagger/ to see all endpoints and to send requests.

> [!IMPORTANT]
> Method names for JSONRpc located in route description.

