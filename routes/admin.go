package routes

// import (
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"time"

// 	"github.com/Blitz-Cloud/ettiHelper/middleware"
// 	"github.com/Blitz-Cloud/ettiHelper/types"
// 	"github.com/davecgh/go-spew/spew"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// 	"gorm.io/gorm"
// )

// func RegisterAdminRoutes(app *fiber.App, serverLogger *log.Logger) {

// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading env file")
// 	}

// 	adminGroup := app.Group("/api/admin")

// 	//UNPROTECTED
// 	adminGroup.Get("/last-sync", func(c *fiber.Ctx) error {
// 		wd, err := os.Getwd()
// 		if err != nil {
// 			serverLogger.Println(err)
// 			return c.SendStatus(fiber.StatusInternalServerError)
// 		}
// 		content, err := os.ReadFile(
// 			filepath.Join(

// 				wd,
// 				"./sync.txt"))

// 		if err != nil {
// 			serverLogger.Println(err)
// 			return c.SendStatus(fiber.StatusInternalServerError)
// 		}
// 		spew.Dump(content)

// 		return c.SendString(string(content))
// 	})

// 	adminGroup.Use(middleware.AdminRouteProtector)

// 	adminGroup.Post("/post/:postType", func(c *fiber.Ctx) error {
// 		postType := c.Params("postType")
// 		db, ok := c.Locals("db").(*gorm.DB)
// 		if !ok {
// 			serverLogger.Fatal("Error accessing db con from admin route")
// 		}
// 		switch postType {
// 		case "blog":
// 			var data types.Blog
// 			c.BodyParser(&data)
// 			spew.Dump(data)
// 			result := db.Create(&data)
// 			if result.Error != nil {
// 				return c.SendStatus(fiber.StatusInternalServerError)
// 			}

// 		case "lab":
// 			var data types.Lab
// 			c.BodyParser(&data)
// 			spew.Dump(data)
// 			result := db.Create(&data)
// 			if result.Error != nil {
// 				return c.SendStatus(fiber.StatusInternalServerError)
// 			}
// 		}

// 		return c.SendStatus(200)
// 	})

// 	// protected
// 	adminGroup.Post("/last-sync", func(c *fiber.Ctx) error {
// 		currentTime := time.Now().UTC().Local().Format(time.RFC3339)
// 		os.WriteFile("./sync.txt", []byte(currentTime), 0777)
// 		return c.SendStatus(fiber.StatusOK)
// 	})
// }
