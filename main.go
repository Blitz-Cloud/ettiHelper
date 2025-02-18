package main

import (
	"log"
	"os"

	"github.com/Blitz-Cloud/ettiHelper/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/mustache/v2"
	"github.com/joho/godotenv"
)

type Post struct {
	title   string
	data    string
	content string
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env file not loaded or missing")
	}
	serverLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	// initializarea bazei de date

	// initializing the fiber app and setting the view engine
	engine := mustache.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layout/main",
	})
	app.Static("/static", "./static")

	// logging
	app.Use(logger.New())

	// routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/login")
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		date := c.Cookies("testC")
		if date != "" {
			return c.Redirect("/labs")
		}
		return c.Render("login", fiber.Map{})
	})

	// microsoft flow
	routes.RegisterMicrosoftOAuth(app)

	// register
	routes.RegisterLabsRoutes(app, serverLogger)
	routes.RegisterTipizateRoutes(app, serverLogger)
	routes.RegisterBlogRoutes(app, serverLogger)

	app.Listen(":3000")
}
