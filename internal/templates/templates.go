package templates

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hunterwilkins2/go-kanban/internal/config"
)

func Render(cfg *config.Config, c *fiber.Ctx, name string, bind map[string]interface{}, layouts ...string) error {
	if bind == nil {
		bind = map[string]interface{}{}
	}

	bind["HotReload"] = cfg.HotReload
	bind["CurrentYear"] = time.Now().Year()

	// if len(layouts) == 0 {
	// 	layouts = append(layouts, "base")
	// }

	return c.Render(name, bind, layouts...)
}
