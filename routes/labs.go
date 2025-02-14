package routes

import (
	"fmt"
	"log"

	"github.com/Blitz-Cloud/ettiHelper/middleware"
	"github.com/Blitz-Cloud/ettiHelper/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterLabsRoutes(app *fiber.App, serverLogger *log.Logger) {

	// initializarea asa zisei baze de date
	// sper sa pot face transferul spre o baza de date adevarata
	var exampleRoot utils.FsNode
	var examples []utils.Example
	utils.Explorer("/home/ionut/facultate/seminar", &exampleRoot, ".c", &examples, utils.LabsContentParser)
	utils.SortDescending(&examples)
	serverLogger.Printf("Explorer a gasit %d coduri scris la laboratoare", len(examples))

	authGroup := app.Group("/", middleware.RouteProtector)

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

		return c.Render("posts", fiber.Map{"posts": examples,
			"Title": "Posts"})
	})

	authGroup.Get("/lab/:date/:name", func(c *fiber.Ctx) error {
		date := c.Params("date")
		name := c.Params("name")
		example := new(utils.Example)
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
		example := new(utils.Example)
		for i := 0; i < len(examples); i++ {
			if examples[i].Name == name && examples[i].Date == date {

				example = &examples[i]
				break
			}
		}
		return c.SendString(example.Content)
	})

}
