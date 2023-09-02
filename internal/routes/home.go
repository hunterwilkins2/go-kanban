package routes

import (
	"github.com/gofiber/fiber/v2"
)

func (app *Application) homepage(c *fiber.Ctx) error {
	return c.Render("pages/home", app.props(nil))
}
