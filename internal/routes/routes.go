package routes

import (
	"database/sql"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/hunterwilkins2/go-kanban/internal/config"
	"github.com/hunterwilkins2/go-kanban/internal/data"
	"github.com/hunterwilkins2/go-kanban/internal/templates"
)

type Application struct {
	config  *config.Config
	server  *fiber.App
	queries *data.Queries
	store   *session.Store
}

func New(server *fiber.App, db *sql.DB, cfg *config.Config) *Application {
	return &Application{
		config:  cfg,
		server:  server,
		queries: data.New(db),
		store:   session.New(),
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

	app.server.Get("/signup", app.registerPage)
	app.server.Get("/login", app.loginPage)
	app.server.Post("/signup", app.register)
	app.server.Post("/login", app.login)
	app.server.Post("/logout", app.logout)

	auth := app.server.Group("/", func(c *fiber.Ctx) error {
		sess, err := app.store.Get(c)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		id := sess.Get("id")
		if id == nil {
			return fiber.ErrForbidden
		}

		return c.Next()
	})

	auth.Get("/", app.homepage)
	auth.Get("/new-board", app.newBoard)

	auth.Post("/board", app.createBoard)
	auth.Get("/board/:slug", app.kanban)
	auth.Get("/edit/:slug", app.editBoard)
	auth.Patch("/board/:slug", app.updateBoard)
	auth.Delete("/board/:slug", app.deleteBoard)

	auth.Get("/board/:slug/new", app.newColumn)
	auth.Post("/board/:slug/column", app.createColumn)
	auth.Delete("/board/:slug/column/:id", app.deleteColumn)
	auth.Get("/board/:slug/column/:id", app.editColumn)
	auth.Patch("/board/:slug/column/:id", app.updateColumn)
	auth.Post("/columns", app.updateColumnOrder)
	auth.Post("/columns/:id", app.createItem)
	auth.Post("/items", app.updateItemOrder)
	auth.Delete("/items/:id", app.deleteItem)
	auth.Get("/items/:id/edit", app.editItem)
	auth.Patch("/items/:id", app.updateItem)
}

func (app *Application) Render(c *fiber.Ctx, name string, bind map[string]interface{}, layouts ...string) error {
	sess, err := app.store.Get(c)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	fullname := sess.Get("name")
	if bind == nil {
		bind = map[string]interface{}{}
	}

	if fullname != nil {
		bind["UserName"] = fullname.(string)
	}

	return templates.Render(app.config, c, name, bind, layouts...)
}
