package routes

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Blitz-Cloud/ettiHelper/types"
	"github.com/davecgh/go-spew/spew"
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

	app.Post("/api/admin/last-sync", func(c *fiber.Ctx) error {
		currentTime := time.Now().UTC().Local().Format(time.RFC3339)
		os.WriteFile("./sync.txt", []byte(currentTime), 0777)
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/last-sync", func(c *fiber.Ctx) error {
		wd, err := os.Getwd()
		if err != nil {
			serverLogger.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		content, err := os.ReadFile(
			filepath.Join(

				wd,
				"./sync.txt"))

		if err != nil {
			serverLogger.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		spew.Dump(content)

		return c.SendString(string(content))
	})
	// apiGroup := app.Group("/api", middleware.ValidateJwtMiddleware)
	apiGroup := app.Group("/api")

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
