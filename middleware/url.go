package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func UriRewriter(c *fiber.Ctx) error {
	hostname := c.GetReqHeaders()["Host"][0]
	if strings.Contains(hostname, "api") {
		c.Path("/api" + c.Path())
	}
	return c.Next()
}
