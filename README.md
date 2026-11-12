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

You don't need Go, gcc, or anything else — `lazydb` ships as a single static
binary for every major OS.

### Windows

**Easiest — download and run:**

1. Go to the [latest release](https://github.com/HalxDocs/lazydb/releases/latest).
2. Download `lazydb-windows-amd64.exe`.
3. Rename it to `lazydb.exe` and put it in a folder on your `PATH`
   (e.g. `C:\Users\<you>\bin`), or run it from where you saved it.

**One-line install in PowerShell:**
```powershell
$dest = "$env:USERPROFILE\bin"
New-Item -ItemType Directory -Force -Path $dest | Out-Null
Invoke-WebRequest -Uri "https://github.com/HalxDocs/lazydb/releases/latest/download/lazydb-windows-amd64.exe" -OutFile "$dest\lazydb.exe"
# add $dest to PATH if it isn't already, then:
lazydb --help
```

**One-line install in Command Prompt:**
```cmd
mkdir "%USERPROFILE%\bin" 2>nul
curl -L -o "%USERPROFILE%\bin\lazydb.exe" https://github.com/HalxDocs/lazydb/releases/latest/download/lazydb-windows-amd64.exe
"%USERPROFILE%\bin\lazydb.exe" --help
```

### macOS

```bash
# Apple Silicon (M1/M2/M3/M4)
curl -L https://github.com/HalxDocs/lazydb/releases/latest/download/lazydb-darwin-arm64 -o /usr/local/bin/lazydb
chmod +x /usr/local/bin/lazydb

# Intel
curl -L https://github.com/HalxDocs/lazydb/releases/latest/download/lazydb-darwin-amd64 -o /usr/local/bin/lazydb
chmod +x /usr/local/bin/lazydb
```

> macOS Gatekeeper may block the binary the first time. Run it once, then go
> to *System Settings → Privacy & Security → "lazydb was blocked" → Open Anyway*,
> or run `xattr -dr com.apple.quarantine /usr/local/bin/lazydb`.

### Linux

```bash
# x86_64
curl -L https://github.com/HalxDocs/lazydb/releases/latest/download/lazydb-linux-amd64 -o /usr/local/bin/lazydb
sudo chmod +x /usr/local/bin/lazydb

# arm64 (Raspberry Pi 4/5, Ampere, etc.)
curl -L https://github.com/HalxDocs/lazydb/releases/latest/download/lazydb-linux-arm64 -o /usr/local/bin/lazydb
sudo chmod +x /usr/local/bin/lazydb
```

### Homebrew (macOS / Linux)
```bash
brew install HalxDocs/tap/lazydb
```

### Go install (for Go developers)
```bash
go install github.com/HalxDocs/lazydb@latest
```

### Build from source
```bash
git clone https://github.com/HalxDocs/lazydb.git
cd lazydb
go build -o lazydb .
./lazydb --help
```

### Verify the install
```bash
lazydb --help
```
You should see the usage banner. If "lazydb is not recognized," the binary
isn't on your `PATH` — either move it into a `PATH` folder or call it by its
full path (e.g. `C:\Users\you\bin\lazydb.exe`).

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
