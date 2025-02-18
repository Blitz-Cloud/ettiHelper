package routes

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/Blitz-Cloud/ettiHelper/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
)

func RegisterBlogRoutes(app *fiber.App, serverLogger *log.Logger) {
	var blogPostRoot utils.FsNode
	var blogPosts []utils.BlogPost
	utils.Explorer("./content", &blogPostRoot, ".md", &blogPosts, utils.MdContentParser)
	serverLogger.Printf("Explorer found %d blog posts", len(blogPosts))
	utils.SortBlogPostsInDescendingOrderByDate(&blogPosts)
	spew.Dump(blogPosts)
	serverLogger.Println("Finished sorting blog posts")
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
	// not finished
	app.Get("/blog", func(c *fiber.Ctx) error {
		return c.Render("Posts", fiber.Map{
			"posts": blogPosts, "linkTo": "blog",
		})
	})
	app.Get("/blog/:date/:name", func(c *fiber.Ctx) error {
		date := c.Params("date")
		name, _ := url.QueryUnescape(c.Params("name"))
		post := utils.BlogPost{}
		for i := 0; i < len(blogPosts); i++ {
			if blogPosts[i].Date == date && blogPosts[i].Title == name {
				fmt.Println("True")
				post = blogPosts[i]
			}
		}
		return c.Render("blogPost", fiber.Map{"content": post.HtmlContent})
	})

}
