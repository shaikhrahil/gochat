package templates

import "github.com/gofiber/fiber/v2"

func Render(name string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render(name, fiber.Map{})
	}
}
