package routes

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/Blitz-Cloud/ettiHelper/middleware"
// 	"github.com/Blitz-Cloud/ettiHelper/utils"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// )

// func RegisterTipizateRoutes(app *fiber.App, serverLogger *log.Logger) {
// 	var tipizateRoot utils.FsNode
// 	var tipizate []utils.BlogPost

// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading env file")
// 	}

// 	tipizateFolder := os.Getenv("tipizateFolder")
// 	serverLogger.Printf("The blog folder location is set to: %s", tipizateFolder)

// 	utils.Explorer(tipizateFolder, &tipizateRoot, ".md", &tipizate, utils.MdContentParser)
// 	serverLogger.Printf("Explorer found %d c code examples", len(tipizate))

// 	authGroup := app.Group("/", middleware.RouteProtector)

// 	authGroup.Get("/tipizate", func(c *fiber.Ctx) error {
// 		days := make([]string, len(tipizate))
// 		for _, file := range tipizate {
// 			ok := 1
// 			for i := 0; i < len(days); i++ {
// 				if days[i] == file.Date {
// 					ok = 0
// 				}
// 			}
// 			if ok == 1 {
// 				days = append(days, file.Date)
// 			}
// 		}

// 		return c.Render("Posts", fiber.Map{"posts": tipizate,
// 			"Title": "Posts", "linkTo": "tipizat"})
// 	})

// 	authGroup.Get("/tipizat/:date/:name", func(c *fiber.Ctx) error {
// 		fmt.Println("Here")
// 		fmt.Println(len(tipizate))
// 		name := c.Params("name")
// 		example := new(utils.BlogPost)
// 		previousPost := ""
// 		nextPost := ""
// 		for i := 0; i < len(tipizate); i++ {
// 			if tipizate[i].Title == name {

// 				example = &tipizate[i]
// 				if i == 0 {

// 					nextPost = fmt.Sprintf("%s/%s", tipizate[i].Date, tipizate[i].Title)
// 				} else {
// 					nextPost = fmt.Sprintf("%s/%s", tipizate[i-1].Date, tipizate[i-1].Title)
// 				}
// 				if i == len(tipizate)-1 {

// 					previousPost = fmt.Sprintf("%s/%s", tipizate[i].Date, tipizate[i].Title)
// 				} else {
// 					previousPost = fmt.Sprintf("%s/%s", tipizate[i+1].Date, tipizate[i+1].Title)
// 				}
// 				break
// 			}
// 		}

// 		return c.Render("lab", fiber.Map{
// 			"post":         example,
// 			"linkTo":       "tipizat",
// 			"previousPost": previousPost,
// 			"nextPost":     nextPost,
// 		})
// 	})
// 	authGroup.Get("/api/tipizat/:name", func(c *fiber.Ctx) error {
// 		name := c.Params("name")
// 		example := new(utils.BlogPost)
// 		for i := 0; i < len(tipizate); i++ {
// 			if tipizate[i].Title == name {

// 				example = &tipizate[i]
// 				break
// 			}
// 		}
// 		return c.SendString(example.Content)
// 	})
// }
