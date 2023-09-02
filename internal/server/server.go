package server

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/template/html/v2"
	"github.com/hunterwilkins2/go-kanban/internal/config"
	"github.com/hunterwilkins2/go-kanban/internal/routes"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	config *config.Config
	server *fiber.App
	db     *sql.DB
}

func New(cfg *config.Config) *Server {
	db, err := openDB(cfg.DSN)
	if err != nil {
		log.Fatalf("could not connect to the database. %s", err.Error())
	}

	engine := html.New("ui/html", ".html")
	server := fiber.New(fiber.Config{
		Views:                 engine,
		DisableStartupMessage: true,
	})

	ser := &Server{
		config: cfg,
		server: server,
		db:     db,
	}

	routes := routes.New(server, db, cfg)
	routes.Register()

	return ser
}

func (s *Server) Start() {
	defer s.db.Close()

	shutdownErr := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		log.Infow("shutting down server", "signal", sig.String())

		shutdownErr <- s.server.Shutdown()
	}()

	log.Infow("starting server", "port", s.config.Port)
	err := s.server.Listen(fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		log.Fatalf("uncaught error occurred. %s", err.Error())
	}

	if err = <-shutdownErr; err != nil {
		log.Fatalf("error shutting down server. %s", err.Error())
	}

	log.Info("stopped server")
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}