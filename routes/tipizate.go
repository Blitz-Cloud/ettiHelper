package routes

import (
	"fmt"
	"log"

	"github.com/Blitz-Cloud/ettiHelper/middleware"
	"github.com/Blitz-Cloud/ettiHelper/utils"
	"github.com/gofiber/fiber/v2"
)

var tipizateRoot utils.FsNode
var tipizate []utils.Tipizat

func RegisterTipizateRoutes(app *fiber.App, serverLogger *log.Logger) {

	utils.Explorer("/home/ionut/facultate/tipizate", &tipizateRoot, ".c", &tipizate, utils.ClangCodeExamplesParser)
	serverLogger.Printf("Explorer found %d c code examples", len(tipizate))

	authGroup := app.Group("/", middleware.RouteProtector)

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
		fmt.Println(len(tipizate))
		name := c.Params("name")
		example := new(utils.Tipizat)
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
				if i == len(tipizate)-1 {

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
	authGroup.Get("/api/tipizat/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		example := new(utils.Tipizat)
		for i := 0; i < len(tipizate); i++ {
			if tipizate[i].Name == name {

				example = &tipizate[i]
				break
			}
		}
		return c.SendString(example.Content)
	})
}
