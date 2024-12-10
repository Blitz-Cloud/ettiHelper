package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache/v2"
)

type Post struct {
	title   string
	data    string
	content string
}

var exampleRoot FsNode
var examples []Example

func main() {
	// initializarea bazei de date
	Explorer("/home/ionut/facultate/seminar", &exampleRoot, &examples)
	SortDescending(examples)
	// pentru debug
	//spew.Dump(examples)

	engine := mustache.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/login")
	})

	// inca este incompleta

	app.Get("/login", func(c *fiber.Ctx) error {
		date := c.Cookies("testC")
		if date != "" {
			return c.Redirect("/posts")
		}
		return c.Render("login", fiber.Map{})
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		data, err := url.ParseQuery(string(c.Body()))
		if err != nil {
			return err
		}
		if data["password"][0] == "h3lloId" {
			c.Cookie(&fiber.Cookie{
				Name:     "testC",
				Value:    "test",
				Expires:  time.Now().Add(7 * 24 * time.Hour),
				HTTPOnly: true,
			})
		}
		fmt.Println(data["password"][0])
		return c.Redirect("/posts")
	})

	app.Get("/post/:day", func(c *fiber.Ctx) error {
		day := c.Params("day")
		examplesByDay := make([]string, len(examples))
		for i := 0; i < len(examples); i++ {
			if day == examples[i].Date {
				examplesByDay = append(examplesByDay, examples[i].Name)
			}
		}
		return c.Render("post", fiber.Map{"posts": examplesByDay})
	})

	app.Get("/posts", func(c *fiber.Ctx) error {
		date := c.Cookies("testC")
		if date != "" {
			days := make([]string, len(examples))
			for _, file := range examples {
				ok := 1
				for i := 0; i < len(days); i++ {
					if days[i] == file.Date {
						ok = 0
					}
				}
				if ok == 1 {
					days = append(days, file.Date)
				}
			}

			return c.Render("post", fiber.Map{"posts": examples})
		} else {
			return c.Redirect("/login")
		}
	})

	app.Listen(":3000")
	fmt.Printf("Hello World")
}
