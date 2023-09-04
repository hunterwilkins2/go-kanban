package routes

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hunterwilkins2/go-kanban/internal/data"
)

type kanbanColumn struct {
	Slug  string
	Items []data.GetItemsRow
	data.GetColumnsRow
}

func (app *Application) kanban(c *fiber.Ctx) error {
	slug := c.Params("slug")

	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	board, err := app.queries.GetBoard(timeout, slug)
	cancel()
	if err != nil {
		return app.Render(c, "pages/not-found", map[string]interface{}{"Title": "404 - Not Found"}, "base")
	}

	timeout, cancel = context.WithTimeout(c.UserContext(), 5*time.Second)
	columns, err := app.queries.GetColumns(timeout, board.ID)
	cancel()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	var kanbanColumns []kanbanColumn
	for _, column := range columns {
		timeout, cancel = context.WithTimeout(c.UserContext(), 5*time.Second)
		items, err := app.queries.GetItems(timeout, column.ID)
		cancel()
		if err != nil {
			return fiber.ErrInternalServerError
		}

		newColumn := kanbanColumn{
			board.Slug,
			items,
			column,
		}
		kanbanColumns = append(kanbanColumns, newColumn)
	}

	return app.Render(c, "pages/kanban", map[string]interface{}{
		"Title":   board.Name,
		"Name":    board.Name,
		"Slug":    board.Slug,
		"Columns": kanbanColumns,
	}, "base")
}

func (app *Application) newColumn(c *fiber.Ctx) error {
	slug := c.Params("slug")
	return app.Render(c, "partials/new-column", map[string]interface{}{"Slug": slug})
}

func (app *Application) createColumn(c *fiber.Ctx) error {
	slug := c.Params("slug")
	name := c.FormValue("name")

	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	board, err := app.queries.GetBoard(timeout, slug)
	cancel()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	timeout, cancel = context.WithTimeout(c.UserContext(), 5*time.Second)
	count, err := app.queries.CountColumns(timeout, board.ID)
	cancel()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	timeout, cancel = context.WithTimeout(c.UserContext(), 5*time.Second)
	id, err := app.queries.CreateColumn(timeout, data.CreateColumnParams{
		Name:         name,
		ElementOrder: count,
		BoardID:      board.ID,
	})
	cancel()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return app.Render(c, "column", map[string]interface{}{"Name": name, "Slug": board.Slug, "ID": id})
}

func (app *Application) deleteColumn(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	err = app.queries.DeleteColumn(timeout, id)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}

func (app *Application) editColumn(c *fiber.Ctx) error {
	slug := c.Params("slug")
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	column, err := app.queries.GetColumn(timeout, id)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return app.Render(c, "partials/edit-column", map[string]interface{}{"Slug": slug, "ID": id, "Name": column.Name})
}

func (app *Application) updateColumn(c *fiber.Ctx) error {
	name := c.FormValue("name")
	slug := c.Params("slug")
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	column, err := app.queries.GetColumn(timeout, id)
	cancel()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	timeout, cancel = context.WithTimeout(c.UserContext(), 5*time.Second)
	err = app.queries.UpdateColumn(timeout, data.UpdateColumnParams{
		Name:         name,
		ElementOrder: column.ElementOrder,
		ID:           id,
	})
	cancel()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return app.Render(c, "partials/column-name", map[string]interface{}{"Slug": slug, "ID": id, "Name": name})
}

func (app *Application) updateColumnOrder(c *fiber.Ctx) error {
	columns := struct {
		Cols []int `json:"columns"`
	}{}
	err := c.BodyParser(&columns)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	for index, id := range columns.Cols {
		timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
		err := app.queries.SetColumnOrder(timeout, data.SetColumnOrderParams{
			ID:           int64(id),
			ElementOrder: int64(index),
		})
		cancel()
		if err != nil {
			return fiber.ErrInternalServerError
		}
	}

	return nil
}

func (app *Application) createItem(c *fiber.Ctx) error {
	name := c.FormValue("name")
	colIdStr := c.Params("id")
	colId, err := strconv.ParseInt(colIdStr, 10, 64)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	count, err := app.queries.CountItems(timeout, colId)
	cancel()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	timeout, cancel = context.WithTimeout(c.UserContext(), 5*time.Second)
	id, err := app.queries.CreateItem(timeout, data.CreateItemParams{
		Name:         name,
		ElementOrder: count,
		ColumnID:     colId,
	})
	cancel()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return app.Render(c, "partials/item", map[string]interface{}{"Name": name, "ID": id})
}

func (app *Application) deleteItem(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	err = app.queries.DeleteItem(timeout, id)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}

func (app *Application) editItem(c *fiber.Ctx) error {
	id := c.Params("id")
	name := c.Query("name")
	return app.Render(c, "partials/edit-item", map[string]interface{}{"ID": id, "Name": name})
}

func (app *Application) updateItem(c *fiber.Ctx) error {
	name := c.FormValue("name")
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	err = app.queries.UpdateItem(timeout, data.UpdateItemParams{
		Name: name,
		ID:   id,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return app.Render(c, "item-name", map[string]interface{}{"ID": id, "Name": name})
}

func (app *Application) updateItemOrder(c *fiber.Ctx) error {
	items := struct {
		ColumnId int64 `json:"columnId"`
		Items    []int `json:"items"`
	}{}

	err := c.BodyParser(&items)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	for index, id := range items.Items {
		timeout, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
		err := app.queries.SetItemOrder(timeout, data.SetItemOrderParams{
			ElementOrder: int64(index),
			ColumnID:     items.ColumnId,
			ID:           int64(id),
		})
		cancel()
		if err != nil {
			return fiber.ErrInternalServerError
		}
	}

	return nil
}
