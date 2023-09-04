package routes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hunterwilkins2/go-kanban/internal/data"
)

func (app *Application) homepage(c *fiber.Ctx) error {
	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	boards, err := app.queries.GetBoards(timeout)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return app.Render(c, "pages/home", map[string]interface{}{"Title": "Go Kanban", "Boards": boards}, "base")
}

func (app *Application) newBoard(c *fiber.Ctx) error {
	return app.Render(c, "partials/new-board", nil)
}

func generateSlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}

func (app *Application) createBoard(c *fiber.Ctx) error {
	name := c.FormValue("name")
	slug := generateSlug(name)

	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	err := app.queries.CreateBoard(timeout, data.CreateBoardParams{
		Name: name,
		Slug: slug,
	})

	if err != nil {
		return fiber.ErrInternalServerError
	}

	return app.Render(c, "board", map[string]interface{}{
		"Name": name,
		"Slug": slug,
	})
}

func (app *Application) deleteBoard(c *fiber.Ctx) error {
	slug := c.Params("slug")

	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	err := app.queries.DeleteBoard(timeout, slug)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}

func (app *Application) editBoard(c *fiber.Ctx) error {
	slug := c.Params("slug")

	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	board, err := app.queries.GetBoard(timeout, slug)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return app.Render(c, "partials/edit-board", map[string]interface{}{
		"Name": board.Name,
		"Slug": board.Slug,
	})
}

func (app *Application) updateBoard(c *fiber.Ctx) error {
	name := c.FormValue("name")
	oldSlug := c.Params("slug")
	slug := generateSlug(name)

	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	err := app.queries.UpdateBoard(timeout, data.UpdateBoardParams{
		Name:   name,
		Slug:   slug,
		Slug_2: oldSlug,
	})
	fmt.Println(err)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	c.Set("Hx-Redirect", fmt.Sprintf("/board/%s", slug))
	return nil
}
