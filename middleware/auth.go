package middleware

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
)

func RouteProtector(c *fiber.Ctx) error {
	fmt.Println("Route protector hit")
	date := c.Cookies("testC")
	spew.Dump(date)
	if date != "" {
		fmt.Println("Access Allowed")
		return c.Next()
	}
	return c.Redirect("/login")
}
