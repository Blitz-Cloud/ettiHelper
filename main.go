package main

import (
	"log"
	"os"

	"github.com/Blitz-Cloud/ettiHelper/routes"
	"github.com/Blitz-Cloud/ettiHelper/utils"
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

	// inca este incompleta

	//legacy login

	app.Get("/login", func(c *fiber.Ctx) error {
		date := c.Cookies("testC")
		if date != "" {
			return c.Redirect("/labs")
		}
		return c.Render("login", fiber.Map{})
	})

	// placeholder for content
	app.Get("/blog/recommendation-for-english-presentation", func(c *fiber.Ctx) error {
		data, err := os.ReadFile("./content/englishPresentationRecommendations.md")
		if err != nil {
			log.Fatal("Cant read the file")
		}
		return c.Render("blogPost", fiber.Map{
			"content": string(utils.Md2Html(data)),
		})
	})
	// app.Post("/login", func(c *fiber.Ctx) error {
	// 	data, err := url.ParseQuery(string(c.Body()))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if data["password"][0] == "h3lloId" {
	// 		c.Cookie(&fiber.Cookie{
	// 			Name:     "testC",
	// 			Value:    "test",
	// 			Expires:  time.Now().Add(7 * 24 * time.Hour),
	// 			HTTPOnly: true,
	// 		})
	// 	}
	// 	fmt.Println(data["password"][0])
	// 	return c.Redirect("/posts")
	// })

	// new login flow

	// microsoft flow
	routes.RegisterMicrosoftOAuth(app)
	routes.RegisterLabsRoutes(app, serverLogger)

	// cleaning required here
	// authGroup.Get("/post/:day", func(c *fiber.Ctx) error {
	// 	day := c.Params("day")
	// 	examplesByDay := make([]string, len(examples))
	// 	for i := 0; i < len(examples); i++ {
	// 		if day == examples[i].Date {
	// 			examplesByDay = append(examplesByDay, examples[i].Name)
	// 		}
	// 	}
	// 	spew.Dump(examplesByDay)
	// 	return c.Render("posts", fiber.Map{"posts": examplesByDay})
	// })

	app.Listen(":3000")
}
