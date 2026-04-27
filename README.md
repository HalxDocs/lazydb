# lazydb

A fast, keyboard-driven terminal UI for your database.
Think `lazygit` — but for Postgres, MySQL and SQLite.

![lazydb demo](https://raw.githubusercontent.com/HalxDocs/lazydb/main/assets/demo.gif)

---

## Features

- Browse tables and rows without leaving your terminal
- **Cell-level cursor** — navigate rows and columns independently
- **Schema view** — toggle column types and nullability with `s`
- **Row detail modal** — full values, wrapped, on `Enter`
- **Sort** by any column (asc → desc → none) with `o`
- **Filter** rows with substring match (`f`)
- **Sidebar table search** with `t`
- **Single-line query** (`/`) and **multi-line SQL editor** (`:`)
- **Query history** — up/down arrows in the query bar
- **CSV export** of visible rows (`e`)
- **Yank** current cell or row to the system clipboard (`y` / `Y`)
- **Vim-style paging** — `g`, `G`, `Ctrl+d`, `Ctrl+u`
- NULL, numeric and boolean values colored distinctly
- Async loading with spinner and live latency in the status bar
- Saved connections, single binary, no runtime deps

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

### Save and reuse a connection
```bash
lazydb --driver postgres --dsn "postgres://..." --save myapp
lazydb --conn myapp
```

### List saved connections
```bash
lazydb
```

---

## Keyboard shortcuts

### Navigation
| Key | Action |
|-----|--------|
| `↑` / `↓` / `k` / `j` | move row up / down |
| `←` / `→` / `h` / `l` | move column cursor |
| `g` / `G` | jump to first / last row |
| `Ctrl+d` / `Ctrl+u` | page down / up |
| `Tab` / `Shift+Tab` | next / prev table |
| `[` / `]` | next / prev table |

### Tables & data
| Key | Action |
|-----|--------|
| `Enter` | open row detail |
| `s` | toggle schema view |
| `o` | cycle sort on current column (asc → desc → none) |
| `f` | filter rows |
| `r` | refresh current table |
| `y` / `Y` | yank cell / yank row as TSV |
| `e` | export visible rows to CSV |

### Query
| Key | Action |
|-----|--------|
| `/` | single-line query |
| `:` | multi-line SQL editor |
| `↑` / `↓` (in query) | history previous / next |
| `Enter` | run query |
| `Ctrl+Enter` (editor) | run query |
| `Esc` | close / dismiss |

### App
| Key | Action |
|-----|--------|
| `t` | search tables |
| `?` | toggle help overlay |
| `q` / `Ctrl+c` | quit |

---

## Exports

CSV exports land in `~/.lazydb/exports/<table>-<timestamp>.csv`.

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
