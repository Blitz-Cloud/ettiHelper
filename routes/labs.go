package routes

import (
	"fmt"
	"log"
	"os"

	"github.com/Blitz-Cloud/ettiHelper/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func RegisterLabsRoutes(app *fiber.App, serverLogger *log.Logger) {

	// initializarea asa zisei baze de date
	// sper sa pot face transferul spre o baza de date adevarata
	var exampleRoot utils.FsNode
	var examples []utils.BlogPost

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}

	labsFolder := os.Getenv("labsFolder")
	serverLogger.Printf("The labs folder location is set to: %s", labsFolder)

	utils.Explorer(labsFolder, &exampleRoot, ".md", &examples, utils.MdContentParser)
	serverLogger.Printf("Explorer a gasit %d coduri scrise la laboratoare", len(examples))
	utils.SortBlogPostsInDescendingOrderByDate(&examples)
	serverLogger.Println("Finished sorting labs posts")

	authGroup := app.Group("/")

	authGroup.Get("/allLabs", func(c *fiber.Ctx) error {
		return c.Render("Posts", fiber.Map{"posts": examples})
	})

	authGroup.Get("/labs", func(c *fiber.Ctx) error {
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

		return c.Render("Posts", fiber.Map{"posts": examples,
			"Title": "Posts", "linkTo": "lab"})
	})

	authGroup.Get("/lab/:date/:name", func(c *fiber.Ctx) error {
		date := c.Params("date")
		name := c.Params("name")
		example := new(utils.BlogPost)
		previousPost := ""
		nextPost := ""
		for i := 0; i < len(examples); i++ {
			if examples[i].Title == name && examples[i].Date == date {

				example = &examples[i]
				if i == 0 {

					nextPost = fmt.Sprintf("%s/%s", examples[i].Date, examples[i].Title)
				} else {
					nextPost = fmt.Sprintf("%s/%s", examples[i-1].Date, examples[i-1].Title)
				}
				if i == len(examples)-1 {

					previousPost = fmt.Sprintf("%s/%s", examples[i].Date, examples[i].Title)
				} else {
					previousPost = fmt.Sprintf("%s/%s", examples[i+1].Date, examples[i+1].Title)
				}
				break
			}
		}
		return c.Render("lab", fiber.Map{
			"post":         example,
			"linkTo":       "lab",
			"previousPost": previousPost,
			"nextPost":     nextPost,
		})
	})
	// authGroup.Get("/api/lab/:date/:name", func(c *fiber.Ctx) error {
	// 	date := c.Params("date")
	// 	name := c.Params("name")
	// 	example := new(utils.BlogPost)
	// 	for i := 0; i < len(examples); i++ {
	// 		if examples[i].Title == name && examples[i].Date == date {

	// 			example = &examples[i]
	// 			break
	// 		}
	// 	}
	// 	return c.SendString(example.Content)
	// })
	authGroup.Get("/api/lab/:date/:name", func(c *fiber.Ctx) error {
		date := c.Params("date")
		name := c.Params("name")
		example := new(utils.BlogPost)
		for i := 0; i < len(examples); i++ {
			if examples[i].Title == name && examples[i].Date == date {

				example = &examples[i]
				break
			}
		}
		return c.JSON(example)
	})
}
