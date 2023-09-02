package main

import (
	"flag"

	"github.com/hunterwilkins2/go-kanban/internal/config"
	"github.com/hunterwilkins2/go-kanban/internal/server"
)

var version = "0.1.0"

func main() {
	cfg := config.Config{}

	flag.IntVar(&cfg.Port, "port", 4000, "Server port")
	flag.StringVar(&cfg.DSN, "dsn", "file:db/go-kanban.db", "Sqlite3 Database DSN")
	flag.StringVar(&cfg.Version, "version", version, "Package version")
	flag.StringVar(&cfg.Env, "env", "development", "Server environment (development|production|testing)")
	flag.BoolVar(&cfg.HotReload, "hot-reload", false, "Hot reloads the browers when a change is made")
	flag.Parse()

	server := server.New(&cfg)
	server.Start()
}
