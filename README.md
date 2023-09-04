# Go Kanban

Trello inspired Kanban board created with Go.

![Go Kanban Screenshot](https://github.com/hunterwilkins2/go-kanban/blob/master/img/go-kanban.png)

## Getting Started

### Requirements

- Go v1.20+
- sqlite3
- Tailwindcss
- sqlc

### Run

1. Create the database with `make generate-sql`
2. Run the application with live reloading with `make run/live`
3. Open http://localhost:4000 to view the application

### Build

1. Create the database with `make generate-sql`
2. Build with `make build`
3. Run the binary with `./bin/go-kanban`
4. Open http://localhost:4000 to view the application
