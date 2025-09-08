package routes

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Blitz-Cloud/ettiHelper/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterAdminRoutes(app *fiber.App, serverLogger *log.Logger) {

	adminGroup := app.Group("/api/admin")

	adminGroup.Post("/post/:postType", func(c *fiber.Ctx) error {
		postType := c.Params("postType")
		db, ok := c.Locals("db").(*gorm.DB)
		if !ok {
			serverLogger.Fatal("Error accessing db con from admin route")
		}
		switch postType {
		case "blog":
			var data types.Blog
			c.BodyParser(&data)
			spew.Dump(data)
			result := db.Create(&data)
			if result.Error != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

		case "lab":
			var data types.Lab
			c.BodyParser(&data)
			spew.Dump(data)
			result := db.Create(&data)
			if result.Error != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		return c.SendStatus(200)
	})

	// adminGroup.Get("/", func(c *fiber.Ctx) error {
	// 	currentTime := time.Now()
	// 	data := types.Lab{
	// 		Title:              "Test Data",
	// 		Description:        "Hello world",
	// 		Date:               currentTime,
	// 		Tags:               "test",
	// 		UniYearAndSemester: 12,
	// 		Content:            "```cpp #include <iostream>```",
	// 	}
	// 	return c.JSON(data)
	// })

	adminGroup.Get("/last-sync", func(c *fiber.Ctx) error {
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

	// protected
	adminGroup.Post("/last-sync", func(c *fiber.Ctx) error {
		currentTime := time.Now().Local().UTC().Format(time.RFC3339)
		os.WriteFile("./sync.json", []byte(currentTime), 0777)
		return c.SendStatus(fiber.StatusOK)
	})
}
