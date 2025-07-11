package routes

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/Blitz-Cloud/ettiHelper/middleware"
	"github.com/Blitz-Cloud/ettiHelper/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func RegisterApiRouter(app *fiber.App, serverLogger *log.Logger) {

	var labRoot utils.FsNode
	var labsPosts []utils.BlogPost
	var blogPostRoot utils.FsNode
	var blogPosts []utils.BlogPost

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}

	labsPostsFolder := os.Getenv("labsFolder")
	serverLogger.Printf("The labsPosts folder location is set to: %s", labsPostsFolder)
	utils.Explorer(labsPostsFolder, &labRoot, ".md", &labsPosts, utils.MdContentParser)
	serverLogger.Printf("Explorer a gasit %d coduri scrise la laboratoare", len(labsPosts))

	utils.SortBlogPostsInDescendingOrderByDate(&labsPosts)
	serverLogger.Println("Finished sorting labsPosts posts")

	blogFolder := os.Getenv("blogFolder")
	serverLogger.Printf("The blog folder location is set to: %s", blogFolder)
	utils.Explorer(blogFolder, &blogPostRoot, ".md", &blogPosts, utils.MdContentParser)
	serverLogger.Printf("Explorer found %d blog posts", len(blogPosts))
	utils.SortBlogPostsInDescendingOrderByDate(&blogPosts)
	serverLogger.Println("Finished sorting blog posts")

	apiGroup := app.Group("/api", middleware.ValidateJwtMiddleware)

	apiGroup.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Auth is working")
	})

	apiGroup.Get("/:postType/posts", func(c *fiber.Ctx) error {
		postType := c.Params("postType")
		switch postType {
		case "blog":
			return c.JSON(blogPosts)
		case "labs":
			return c.JSON(labsPosts)
		}
		return c.SendStatus(404)
	})

	apiGroup.Get("/:postType/post/:date/:name", func(c *fiber.Ctx) error {
		postType := c.Params("postType")
		switch postType {
		case "blog":
			date := c.Params("date")
			name, _ := url.QueryUnescape(c.Params("name"))
			post := utils.BlogPost{}
			for i := 0; i < len(blogPosts); i++ {
				if blogPosts[i].Date == date && blogPosts[i].Title == name {
					fmt.Println("True")
					post = blogPosts[i]
				}
			}
			return c.JSON(post)
		case "labs":
			date := c.Params("date")
			name, _ := url.QueryUnescape(c.Params("name"))
			lab := utils.BlogPost{}
			for i := 0; i < len(labsPosts); i++ {
				if labsPosts[i].Date == date && labsPosts[i].Title == name {
					fmt.Println("True")
					lab = labsPosts[i]
				}
			}
			return c.JSON(lab)
		}
		return c.SendStatus(404)
	})
}
