# Electrotech

> We all fucked up with nothing to do

## Development

**Requirements:**

- golang
- sqlc
- sqlite3

Run:

```bash
sqlc generate
go run cmd/server/main.go
```

Application requires two environment variables:

- `DB_CONNECTION` - just path to sqlite DB, for example _"/tmp/electrotech.sqlite3"_
- `DATA_DIR` path to directory with catalog data, leave just _"./example"_

## Contributing

Create fork, create new branch, commit changes and create Pull Request.

PR should contain full description of any changes.
