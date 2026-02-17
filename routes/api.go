package routes

import (
	"errors"

	"github.com/Blitz-Cloud/ettiHelper/middleware"
	"github.com/Blitz-Cloud/ettiHelper/types"
	"github.com/Blitz-Cloud/ettiHelper/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterApiRouter(app *fiber.App) {

	// apiPost := app.Group("/api", middleware.ValidateJwtMiddleware)
	apiPost := app.Group("/api")

	// apiPost.Post("/admin/last-sync", func(c *fiber.Ctx) error {
	// 	currentTime := time.Now().UTC().Local().Format(time.RFC3339)
	// 	os.WriteFile("./sync.txt", []byte(currentTime), 0777)
	// 	return c.SendStatus(fiber.StatusOK)
	// })

	// este necesara o functie de middleware care sa verifice pt toate rutele daca utilizatorul este autentificat in cazul in care se incearca actualizarea sau accesarea unei baze de date a carui access este restrictionat

	// implement since? pentru a prelua update urile 3

	// apiPost.Get("/db-status", func(c *fiber.Ctx) error {
	// 	db, ok := c.Locals("db").(*gorm.DB)
	// 	if !ok {
	// 		utils.Log.Fatal("Error accessing db con from admin route")
	// 	}
	// 	var namespaceCount int64
	// 	var categoryCount int64
	// 	var postCount int64
	// 	db.Table("namespaces").Count(&namespaceCount)
	// 	db.Table("categories").Count(&categoryCount)
	// 	db.Table("posts").Count(&postCount)
	// 	return c.JSON(map[string]int64{
	// 		"namespaces": namespaceCount,
	// 		"categories": categoryCount,
	// 		"posts":      postCount,
	// 	})
	// })

	apiPost.Get("/namespace", func(c *fiber.Ctx) error {
		tenant, ok := c.Locals("tenant").(types.Namespace)
		if !ok {
			utils.Log.Fatal("Error accessing tenant")
		}
		return c.JSON(tenant)
	})

	apiPost.Use(middleware.RouteProtector)
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

	// apiPost.Get("/posts", func(c *fiber.Ctx) error {
	// 	db, ok := c.Locals("db").(*gorm.DB)
	// 	if !ok {
	// 		utils.Log.Fatal("Error accessing db con from admin route")
	// 	}
	// 	posts := []types.Post{}
	// 	if err := db.Find(&posts).Error; err != nil {
	// 		utils.Log.Fatal(err.Error())

	// 	}
	// 	return c.JSON(posts)
	// })

	// apiPost.Get("/:nameSpaceId/categories", func(c *fiber.Ctx) error {
	// 	db, ok := c.Locals("db").(*gorm.DB)
	// 	nameSpaceId := c.Params("nameSpaceId")
	// 	if len(nameSpaceId) == 0 {
	// 		return c.SendStatus(fiber.StatusUnprocessableEntity)
	// 	}

	// 	if !ok {
	// 		utils.Log.Fatal("Error accessing db con from admin route")
	// 	}

	// 	categories := []types.Category{}
	// 	if err := db.Joins("Namespace").Where("namespace_id = ?", nameSpaceId).Find(&categories).Error; err != nil {
	// 		if errors.Is(err, gorm.ErrRecordNotFound) {
	// 			return c.SendStatus(fiber.StatusNotFound)
	// 		}
	// 		utils.Log.Error(err.Error())
	// 		return c.SendStatus(fiber.StatusInternalServerError)
	// 	}
	// 	return c.JSON(categories)
	// })

	// apiPost.Get("/:nameSpaceId/:categoryId", func(c *fiber.Ctx) error {
	// 	db, ok := c.Locals("db").(*gorm.DB)
	// 	nameSpaceId := c.Params("nameSpaceId")
	// 	categoryId := c.Params("categoryId")
	// 	utils.Log.Dump(categoryId)
	// 	if len(nameSpaceId) == 0 {
	// 		return c.SendStatus(fiber.StatusUnprocessableEntity)
	// 	}
	// 	if len(categoryId) == 0 {
	// 		return c.SendStatus(fiber.StatusUnprocessableEntity)
	// 	}

	// 	if !ok {
	// 		utils.Log.Fatal("Error accessing db con from admin route")
	// 	}

	// 	category := types.Category{}
	// 	if err := db.Joins("Namespace").Preload("Posts").Where("namespace_id = ?", nameSpaceId).First(&category, "categories.id = ?", categoryId).Error; err != nil {
	// 		if errors.Is(err, gorm.ErrRecordNotFound) {
	// 			return c.SendStatus(fiber.StatusNotFound)
	// 		}
	// 		utils.Log.Error(err.Error())
	// 		return c.SendStatus(fiber.StatusInternalServerError)
	// 	}

	// 	return c.JSON(category)
	// })

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

	// apiPost.Get("/:nameSpaceId/:categoryId/:postId", func(c *fiber.Ctx) error {
	// 	db, ok := c.Locals("db").(*gorm.DB)
	// 	nameSpaceId := c.Params("nameSpaceId")
	// 	categoryId := c.Params("categoryId")
	// 	postId := c.Params("postId")
	// 	if len(nameSpaceId) == 0 {
	// 		return c.SendStatus(fiber.StatusUnprocessableEntity)
	// 	}
	// 	if len(categoryId) == 0 {
	// 		return c.SendStatus(fiber.StatusUnprocessableEntity)
	// 	}
	// 	if len(postId) == 0 {
	// 		return c.SendStatus(fiber.StatusUnprocessableEntity)
	// 	}

	// 	if !ok {
	// 		utils.Log.Fatal("Error accessing db con from admin route")
	// 	}

	// 	posts := types.Post{}
	// 	if err := db.Where("id = ?", postId).Find(&posts).Error; err != nil {
	// 		if errors.Is(err, gorm.ErrRecordNotFound) {
	// 			return c.SendStatus(fiber.StatusNotFound)
	// 		}
	// 		utils.Log.Error(err.Error())
	// 		return c.SendStatus(fiber.StatusInternalServerError)
	// 	}

	// 	return c.JSON(posts)
	// })

}
