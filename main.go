package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Blitz-Cloud/ettiHelper/middleware"
	"github.com/Blitz-Cloud/ettiHelper/routes"
	"github.com/davecgh/go-spew/spew"
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

var exampleRoot FsNode
var examples []Example
var tipizateRoot FsNode
var tipizate []Tipizat

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env file not loaded or missing")
	}

	// initializarea bazei de date
	Explorer("/home/ionut/facultate/seminar", &exampleRoot, &examples)
	TipizatExplorer("/home/ionut/facultate/tipizate", &tipizateRoot, &tipizate)
	spew.Dump(tipizate[0])
	SortDescending(&examples)

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
			return c.Redirect("/posts")
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
			"content": string(Md2Html(data)),
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

	authGroup := app.Group("/", middleware.RouteProtector)

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

	authGroup.Get("/tipizate", func(c *fiber.Ctx) error {
		days := make([]string, len(tipizate))
		for _, file := range tipizate {
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

		return c.Render("tipizate", fiber.Map{"posts": tipizate,
			"Title": "Posts"})
	})

	authGroup.Get("/tipizat/:name", func(c *fiber.Ctx) error {
		fmt.Println("Here")
		name := c.Params("name")
		example := new(Tipizat)
		previousPost := ""
		nextPost := ""
		for i := 0; i < len(tipizate); i++ {
			if tipizate[i].Name == name {

				example = &tipizate[i]
				if i == 0 {

					nextPost = fmt.Sprintf("%s", tipizate[i].Name)
				} else {
					nextPost = fmt.Sprintf("%s", tipizate[i-1].Name)
				}
				if i == len(examples)-1 {

					previousPost = fmt.Sprintf("%s", tipizate[i].Name)
				} else {
					previousPost = fmt.Sprintf("%s", tipizate[i+1].Name)
				}
				break
			}
		}

		return c.Render("tipizat", fiber.Map{
			"post":         example,
			"previousPost": previousPost,
			"nextPost":     nextPost,
		})
	})
	authGroup.Get("/posts", func(c *fiber.Ctx) error {
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

		return c.Render("posts", fiber.Map{"posts": examples,
			"Title": "Posts"})
	})

	authGroup.Get("/post/:date/:name", func(c *fiber.Ctx) error {
		date := c.Params("date")
		name := c.Params("name")
		example := new(Example)
		previousPost := ""
		nextPost := ""
		for i := 0; i < len(examples); i++ {
			if examples[i].Name == name && examples[i].Date == date {

				example = &examples[i]
				if i == 0 {

					nextPost = fmt.Sprintf("%s/%s", examples[i].Date, examples[i].Name)
				} else {
					nextPost = fmt.Sprintf("%s/%s", examples[i-1].Date, examples[i-1].Name)
				}
				if i == len(examples)-1 {

					previousPost = fmt.Sprintf("%s/%s", examples[i].Date, examples[i].Name)
				} else {
					previousPost = fmt.Sprintf("%s/%s", examples[i+1].Date, examples[i+1].Name)
				}
				break
			}
		}

		return c.Render("post", fiber.Map{
			"post":         example,
			"previousPost": previousPost,
			"nextPost":     nextPost,
		})
	})
	authGroup.Get("/api/post/:date/:name", func(c *fiber.Ctx) error {
		date := c.Params("date")
		name := c.Params("name")
		example := new(Example)
		for i := 0; i < len(examples); i++ {
			if examples[i].Name == name && examples[i].Date == date {

				example = &examples[i]
				break
			}
		}
		return c.SendString(example.Content)
	})

	authGroup.Get("/api/tipizat/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		example := new(Tipizat)
		for i := 0; i < len(tipizate); i++ {
			if tipizate[i].Name == name {

				example = &tipizate[i]
				break
			}
		}
		return c.SendString(example.Content)
	})
	app.Listen(":3000")
	fmt.Printf("Hello World")
}
