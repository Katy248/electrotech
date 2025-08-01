# Running the app

- Build
- Create `.env`
- Migrate
- Run

## Building

```bash
go build -o build/migrator cmd/migrator/main.go
go build -o build/server cmd/server/main.go
```

## Environment variables

- `PORT` - port to run the server on (required)
- `DATA_DIR` - directory with all xml datafile (required)
- `DB_CONNECTION`  - connection string. Currently we use SQLite so it must be path to db-file, for example: _"/var/electrotech/db.sqlite3"_ (required)

All these variables must be set in `.env` file.

## Migration

```bash
./build/migrator
```

## Running

```bash
./build/server
```
