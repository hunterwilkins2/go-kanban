package routes

import (
	"github.com/gofiber/fiber/v2"
)

func (app *Application) homepage(c *fiber.Ctx) error {
	return app.Render(c, "pages/home", nil)
}
