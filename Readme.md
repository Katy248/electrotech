# Electrotech

> We all fucked up with nothing to do

## Development

**Requirements:**

- golang
- sqlite3

Run:

```bash
DEVEL=true go run cmd/server/main.go
```

Application requires several environment variables:

- `DB_CONNECTION` - just path to sqlite DB, for example _"/tmp/electrotech.sqlite3"_
- `DATA_DIR` path to directory with catalog data, default is `/data`. Fo development use `./example`
- `PORT` - HTTP port, default is _8080_
- `JWT_SECRET` - secret for JWT token

It can be configured via `.env` file or `electrotech-back.toml` file (`electrotech-back.devel.toml` for development)

## Contributing

Create fork, create new branch, commit changes and create Pull Request.

PR should contain full description of any changes.
