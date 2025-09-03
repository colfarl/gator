# gator

A tiny CLI for following RSS/Atom feeds and storing posts in Postgres.

---

## Prerequisites

- **Go** 1.21+ → https://go.dev/dl/
- **PostgreSQL** 15+ → https://www.postgresql.org/download/

Verify installs:
- `go version`
- `psql --version`

Create a database (name it anything, e.g. `gator`):
- `createdb gator`  (or run `psql -c 'CREATE DATABASE gator;'`)
---

## Install

Install the CLI with `go install`:
- `go install github.com/colfarl/gator@latest`

Ensure your `$GOPATH/bin` (or `$GOBIN`) is on `PATH` so `gator` is runnable.

From source (optional):
```
git clone https://github.com/colfarl/gator
cd gator
go build -o gator .
```
---

## Configure

`gator` expects a config file **in the HOME directory** named `./gatorconfig.json`.

Create `gatorconfig.json`:
```json
{
  "db_url": "postgres://<connection string ?sslmode=disable>",
  "current_user_name": "default_user"
}
```
## Database Setup (Goose)

This project uses [Goose](https://github.com/pressly/goose) for SQL migrations (migrations live in `.sql/schema`)
to move to most recent version of db run:
```
goose postgres <connection-string > up-to 5 
```

### Commands

- `./gator register <name>` — Add a user to the database.
- `./gator login <name>` — Log in as an existing user.
- `./gator reset` — **Dangerous:** remove the entire database.
- `./gator users` — List all users; highlights the currently logged-in user.
- `./gator agg <duration>` — Poll on an interval (e.g., `1h`, `1m`, `30s`) to fetch new posts from the stalest feed.
- `./gator browse [limit]` — Show the most recent posts from followed feeds (default `2`).
- `./gator addfeed <name> <url>` — Add a feed; fails if it already exists.
- `./gator feeds` — List all feeds in the database.
- `./gator follow <url>` — Follow a feed for the current user.
- `./gator following` — List feeds the current user is following.
- `./gator unfollow <url>` — Unfollow a feed for the current user.
