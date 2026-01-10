package main

import (
	"log"

	"github.com/Blitz-Cloud/ettiHelper/routes"
	"github.com/Blitz-Cloud/ettiHelper/types"
	"github.com/Blitz-Cloud/ettiHelper/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// "github.com/Blitz-Cloud/ettiHelper/routes"

func main() {
	utils.Log.Info("SERVER Started")
	app := fiber.New(fiber.Config{})

	db, err := gorm.Open(sqlite.Open("./ettiContent.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&types.Namespace{}, &types.Category{}, &types.Post{})
	DB, err := utils.InMemoryDB("/home/ionut/project/content/ettiContent/root", "/home/ionut/project/content/ettiContent/etti")
	if err != nil {

		utils.Log.Error(err.Error())
	}
	utils.Log.Info("IMPORTANT")
	err = utils.SeedFromInMemory(db, DB)
	if err != nil {
		utils.Log.Error(err.Error())
	}

	// if err != nil {
	// 	utils.Log.Error(err.Error())
	// }
	err = godotenv.Load()
	if err != nil {
		utils.Log.Fatal("ENV file not loaded or missing")
	}
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", DB)

		return c.Next()
	})

	// // // logging
	// app.Use(logger.New(logger.Config{}))

	// app.Use(cors.New(cors.Config{
	// 	// trebuie sa adaug prod si development aici
	// 	AllowOrigins:     "http://localhost:5173,http://localhost:3000,https://ettiui.netlify.app",
	// 	AllowCredentials: true,
	// }))

	// // // app.Get("/login", func(c *fiber.Ctx) error {
	// // // 	return c.SendString("LoginPage")
	// // // })

	routes.RegisterApiRouter(app)
	// // // routes.RegisterAdminRoutes(app, serverLogger)
	// // app.Get("*", func(c *fiber.Ctx) error {
	// // 	return c.SendFile("./build/client/index.html")
	// // })

	// currentTime := time.Now().UTC().Local().Format(time.RFC3339)
	// os.WriteFile("./sync.txt", []byte(currentTime), 0777)
	app.Listen(":3000")
}
