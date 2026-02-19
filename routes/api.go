package routes

import (
	"errors"

	"github.com/Blitz-Cloud/ettiHelper/types"
	"github.com/Blitz-Cloud/ettiHelper/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterApiRouterUnprotected(app *fiber.App) {

	apiPost := app.Group("/api")

	apiPost.Get("/namespace", func(c *fiber.Ctx) error {
		tenant, ok := c.Locals("tenant").(types.Namespace)
		if !ok {
			utils.Log.Fatal("Error accessing tenant")
		}
		return c.JSON(tenant)
	})

	apiPost.Get("/namespaces", func(c *fiber.Ctx) error {
		db, ok := c.Locals("db").(*gorm.DB)
		if !ok {
			utils.Log.Fatal("Error accessing db con from admin route")
		}
		namespaces := []types.Namespace{}
		if err := db.Find(&namespaces).Error; err != nil {
			utils.Log.Fatal(err.Error())
		}
		return c.JSON(namespaces)
	})
}

func RegisterApiRouterProtected(app *fiber.App) {

	apiPost := app.Group("/api")

	apiPost.Get("/categories", func(c *fiber.Ctx) error {
		db, ok := c.Locals("db").(*gorm.DB)
		if !ok {
			utils.Log.Fatal("Error accessing db con from admin route")
		}
		tenant, ok := c.Locals("tenant").(types.Namespace)
		if !ok {
			utils.Log.Fatal("Error accessing tenant")
		}

		categories := []types.Category{}
		if err := db.Find(&categories, "namespace_Id = ?", tenant.ID).Error; err != nil {
			utils.Log.Fatal(err.Error())
		}
		return c.JSON(categories)
	})

	apiPost.Get("/:categoryId/posts", func(c *fiber.Ctx) error {
		db, ok := c.Locals("db").(*gorm.DB)
		categoryId := c.Params("categoryId")
		if len(categoryId) == 0 {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}
		if !ok {
			utils.Log.Fatal("Error accessing db con from admin route")
		}
		posts := []types.Post{}
		if err := db.Find(&posts, "category_id = ?", categoryId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.SendStatus(fiber.StatusNotFound)
			}
			utils.Log.Error(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(posts)
	})

	apiPost.Get("/post/:postId", func(c *fiber.Ctx) error {
		db, ok := c.Locals("db").(*gorm.DB)
		postId := c.Params("postId")
		if len(postId) == 0 {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}
		if !ok {
			utils.Log.Fatal("Error accessing db con from admin route")
		}

		post := types.Post{}

		if err := db.Find(&post, "id = ?", postId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.SendStatus(fiber.StatusNotFound)
			}
			utils.Log.Error(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(post)
	})
}
