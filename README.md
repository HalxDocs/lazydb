# lazydb

A fast, keyboard-driven terminal UI for your database.  
Think `lazygit` — but for Postgres, MySQL and SQLite.

![lazydb demo](https://raw.githubusercontent.com/HalxDocs/lazydb/main/assets/demo.gif)

---

## Features

- Browse tables and rows without leaving your terminal
- Run raw SQL queries with instant results
- Row counts per table shown in the sidebar
- Highlighted row navigation with keyboard
- Save and reuse database connections by name
- Works with Postgres, MySQL and SQLite
- Single binary — no runtime, no dependencies, no install friction

---

## Install

### Homebrew (macOS / Linux)
```bash
brew install HalxDocs/tap/lazydb
```

### Go install
```bash
go install github.com/HalxDocs/lazydb@latest
```

### Download binary
Grab the latest binary for your platform from the [releases page](https://github.com/HalxDocs/lazydb/releases).

---

## Usage

### Connect directly
```bash
# Postgres
lazydb --driver postgres --dsn "postgres://user:pass@localhost:5432/mydb?sslmode=disable"

# MySQL
lazydb --driver mysql --dsn "user:pass@tcp(localhost:3306)/mydb"

# SQLite
lazydb --driver sqlite --dsn ./mydb.sqlite
```

### Save a connection
```bash
lazydb --driver postgres --dsn "postgres://user:pass@localhost:5432/mydb?sslmode=disable" --save myapp
```

### Use a saved connection
```bash
lazydb --conn myapp
```

### List saved connections
```bash
lazydb
```

---

## Keyboard shortcuts

| Key | Action |
|-----|--------|
| `↑` / `k` | Move row up |
| `↓` / `j` | Move row down |
| `←` / `h` | Previous table |
| `→` / `l` | Next table |
| `/` | Open query bar |
| `Enter` | Run query |
| `Esc` | Close query bar / dismiss error |
| `q` | Quit |

---

## Roadmap

- [ ] Expand row detail view on `Enter`
- [ ] Edit and delete rows
- [ ] Export query results to CSV
- [ ] Multiple simultaneous connections
- [ ] Indexes and schema view
- [ ] MySQL and SQLite full testing

---

## Why lazydb?

TablePlus costs money. DBeaver is slow Java bloat. pgAdmin is a browser app.  
Developers who live in the terminal have had nothing polished — until now.

`lazydb` is the missing tool. Single binary. Zero config. Works over SSH.  
It feels native because it is native.

---

## Contributing

Pull requests are welcome. For major changes please open an issue first.
```bash
git clone https://github.com/HalxDocs/lazydb.git
cd lazydb
go build ./...
go run main.go --driver sqlite --dsn ./test.sqlite
```

---

## License

MIT