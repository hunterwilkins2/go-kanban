package routes

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (app *Application) hotReload(c *websocket.Conn) {
	for {
		log.Debug("hot reload connection successful")
		msgType, _, err := c.ReadMessage()
		if err != nil {
			return
		}

		if err = c.WriteMessage(msgType, []byte("connected")); err != nil {
			return
		}
	}
}

func (app *Application) testAlive(c *fiber.Ctx) error {
	return c.SendString("hot reload alive")
}
