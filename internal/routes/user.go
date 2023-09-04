package routes

import "github.com/gofiber/fiber/v2"

func (app *Application) registerPage(c *fiber.Ctx) error {
	return app.Render(c, "page/register", nil)
}

func (app *Application) register(c *fiber.Ctx) error {
	return nil
}

func (app *Application) loginPage(c *fiber.Ctx) error {
	return app.Render(c, "page/login", nil)
}

func (app *Application) login(c *fiber.Ctx) error {
	return nil
}
