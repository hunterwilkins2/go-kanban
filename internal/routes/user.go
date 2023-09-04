package routes

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hunterwilkins2/go-kanban/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) registerPage(c *fiber.Ctx) error {
	return app.Render(c, "pages/register", map[string]interface{}{"Title": "Sign Up"}, "base")
}

func (app *Application) register(c *fiber.Ctx) error {
	sess, err := app.store.Get(c)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	id, err := app.queries.CreateUser(timeout, data.CreateUserParams{
		Fullname:     name,
		Email:        email,
		PasswordHash: string(hash),
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}

	sess.Set("name", name)
	sess.Set("id", id)
	if err = sess.Save(); err != nil {
		panic(err)
	}

	c.Set("Hx-Redirect", "/")
	return nil
}

func (app *Application) loginPage(c *fiber.Ctx) error {
	return app.Render(c, "pages/login", map[string]interface{}{"Title": "Login"}, "base")
}

func (app *Application) login(c *fiber.Ctx) error {
	sess, err := app.store.Get(c)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	user, err := app.queries.GetUser(timeout, email)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return fiber.ErrForbidden
	}

	sess.Set("name", user.Fullname)
	sess.Set("id", user.ID)
	if err = sess.Save(); err != nil {
		panic(err)
	}

	c.Set("Hx-Redirect", "/")
	return nil
}

func (app *Application) logout(c *fiber.Ctx) error {
	sess, err := app.store.Get(c)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	err = sess.Destroy()
	if err != nil {
		panic(err)
	}

	c.Set("Hx-Redirect", "/login")
	return nil
}
