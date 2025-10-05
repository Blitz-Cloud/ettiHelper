package routes

import (
	"log"

	"github.com/Blitz-Cloud/ettiHelper/middleware"
	"github.com/Blitz-Cloud/ettiHelper/types"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func RegisterApiRouter(app *fiber.App, serverLogger *log.Logger) {

	// i leave this here in case i need it
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}

	apiGroup := app.Group("/api", middleware.ValidateJwtMiddleware)
	// apiGroup := app.Group("/api")

	apiGroup.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Auth is working")
	})

	apiGroup.Get("/:postType/posts", func(c *fiber.Ctx) error {
		postType := c.Params("postType")
		db, ok := c.Locals("db").(*gorm.DB)
		if ok {
			switch postType {
			case "blog":
				var posts []types.Blog
				db.Find(&posts)
				return c.JSON(posts)
			case "labs":
				var labs []types.Lab
				db.Find(&labs)
				return c.JSON(labs)
			}
		} else {
			serverLogger.Fatal("Error accessing the db con from gofiber context")
		}

		return c.SendStatus(404)
	})

	apiGroup.Get("/:postType/post/:id", func(c *fiber.Ctx) error {

		postType := c.Params("postType")
		id := c.Params("id")
		db, ok := c.Locals("test").(*gorm.DB)
		if ok {
			switch postType {
			case "blog":
				var post types.Blog
				db.First(&post, id)
				return c.JSON(post)
			case "labs":
				var lab types.Lab
				db.First(&lab, id)
				return c.JSON(lab)
			}
		} else {
			serverLogger.Fatal("Error accessing the db con from gofiber context")
		}
		return c.SendStatus(404)
	})

}
