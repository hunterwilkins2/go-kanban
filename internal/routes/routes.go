package routes

import (
	"database/sql"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hunterwilkins2/go-kanban/internal/config"
	"github.com/hunterwilkins2/go-kanban/internal/data"
)

type Application struct {
	config  *config.Config
	server  *fiber.App
	queries *data.Queries
}

func New(server *fiber.App, db *sql.DB, cfg *config.Config) *Application {
	return &Application{
		config:  cfg,
		server:  server,
		queries: data.New(db),
	}
}

func (app *Application) Register() {
	app.server.Use(csrf.New())
	app.server.Use(logger.New())
	app.server.Use(helmet.New())

	app.server.Static("/static", "ui/static")

	if app.config.HotReload {
		app.server.Get("/hot-reload", websocket.New(app.hotReload))
		app.server.Get("/hot-reload/ready", app.testAlive)
	}

	app.server.Get("/", app.homepage)
}

func (app *Application) Render(c *fiber.Ctx, name string, bind map[string]interface{}, layouts ...string) error {
	if bind == nil {
		bind = map[string]interface{}{}
	}

	bind["HotReload"] = app.config.HotReload
	bind["CurrentYear"] = time.Now().Year()

	return c.Render(name, bind, layouts...)
}
