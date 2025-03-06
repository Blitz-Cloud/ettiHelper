package routes

import (
	"fmt"
	"log"
	"os"

	"github.com/Blitz-Cloud/ettiHelper/middleware"
	"github.com/Blitz-Cloud/ettiHelper/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func RegisterLabsRoutes(app *fiber.App, serverLogger *log.Logger) {

	// initializarea asa zisei baze de date
	// sper sa pot face transferul spre o baza de date adevarata
	var labsRoot utils.FsNode
	var labs []utils.BlogPost

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}

	labsFolder := os.Getenv("labsFolder")
	serverLogger.Printf("The labs folder location is set to: %s", labsFolder)

	utils.Explorer(labsFolder, &labsRoot, ".md", &labs, utils.MdContentParser)
	serverLogger.Printf("Explorer a gasit %d coduri scrise la laboratoare", len(labs))
	utils.SortBlogPostsInDescendingOrderByDate(&labs)
	serverLogger.Println("Finished sorting labs posts")

	authGroup := app.Group("/", middleware.RouteProtector)

	// authGroup.Get("/allLabs", func(c *fiber.Ctx) error {
	// 	return c.Render("Posts", fiber.Map{"posts": labs})
	// })

	authGroup.Get("/labs/:uniYearAndSemester", func(c *fiber.Ctx) error {
		days := make([]string, len(labs))
		for _, file := range labs {
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

		return c.Render("Posts", fiber.Map{"posts": labs,
			"Title": "Posts", "linkTo": "lab"})
	})

	authGroup.Get("/lab/:date/:name", func(c *fiber.Ctx) error {
		date := c.Params("date")
		name := c.Params("name")
		lab := new(utils.BlogPost)
		previousPost := ""
		nextPost := ""
		for i := 0; i < len(labs); i++ {
			if labs[i].Title == name && labs[i].Date == date {

				lab = &labs[i]
				if i == 0 {

					nextPost = fmt.Sprintf("%s/%s", labs[i].Date, labs[i].Title)
				} else {
					nextPost = fmt.Sprintf("%s/%s", labs[i-1].Date, labs[i-1].Title)
				}
				if i == len(labs)-1 {

					previousPost = fmt.Sprintf("%s/%s", labs[i].Date, labs[i].Title)
				} else {
					previousPost = fmt.Sprintf("%s/%s", labs[i+1].Date, labs[i+1].Title)
				}
				break
			}
		}
		lab.Content = string(utils.Md2Html([]byte(lab.Content)))
		return c.Render("lab", fiber.Map{
			"post":         lab,
			"linkTo":       "lab",
			"previousPost": previousPost,
			"nextPost":     nextPost,
		})
	})
	// authGroup.Get("/api/post/:date/:name", func(c *fiber.Ctx) error {
	// 	date := c.Params("date")
	// 	name := c.Params("name")
	// 	lab := new(utils.BlogPost)
	// 	for i := 0; i < len(labs); i++ {
	// 		if labs[i].Title == name && labs[i].Date == date {

	// 			lab = &labs[i]
	// 			break
	// 		}
	// 	}
	// 	return c.SendString(lab.Content)
	// })

}
