package routes

import (
	"database/sql"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hunterwilkins2/go-kanban/internal/config"
	"github.com/hunterwilkins2/go-kanban/internal/data"
	"github.com/hunterwilkins2/go-kanban/internal/templates"
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
	app.server.Use(logger.New())
	app.server.Use(helmet.New())
	app.server.Use(recover.New())

	if app.config.HotReload {
		app.server.Get("/hot-reload", websocket.New(app.hotReload))
		app.server.Get("/hot-reload/ready", app.testAlive)
	}

	app.server.Static("/static", "ui/static")

	app.server.Get("/", app.homepage)
	app.server.Get("/new-board", app.newBoard)

	app.server.Post("/board", app.createBoard)
	app.server.Get("/board/:slug", app.kanban)
	app.server.Get("/edit/:slug", app.editBoard)
	app.server.Patch("/board/:slug", app.updateBoard)
	app.server.Delete("/board/:slug", app.deleteBoard)

	app.server.Get("/board/:slug/new", app.newColumn)
	app.server.Post("/board/:slug/column", app.createColumn)
	app.server.Delete("/board/:slug/column/:id", app.deleteColumn)
	app.server.Get("/board/:slug/column/:id", app.editColumn)
	app.server.Patch("/board/:slug/column/:id", app.updateColumn)
	app.server.Post("/columns", app.updateColumnOrder)
	app.server.Post("/columns/:id", app.createItem)
	app.server.Post("/items", app.updateItemOrder)
	app.server.Delete("/items/:id", app.deleteItem)
	app.server.Get("/items/:id/edit", app.editItem)
	app.server.Patch("/items/:id", app.updateItem)

	app.server.Get("/signup", app.registerPage)
	app.server.Get("/login", app.loginPage)
	app.server.Post("/signup", app.register)
	app.server.Post("/login", app.login)
}

func (app *Application) Render(c *fiber.Ctx, name string, bind map[string]interface{}, layouts ...string) error {
	return templates.Render(app.config, c, name, bind, layouts...)
}
