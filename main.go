package main

import (
	"log"
	"os"

	// "github.com/Blitz-Cloud/ettiHelper/routes"
	"github.com/Blitz-Cloud/ettiHelper/routes"
	"github.com/Blitz-Cloud/ettiHelper/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/mustache/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Post struct {
	title   string
	data    string
	content string
}

// Claims struct to represent the JWT claims
type Claims struct {
	Aud string `json:"aud"`
	Iss string `json:"iss"`
	jwt.RegisteredClaims
}

type DBConType struct{}

var DBCon DBConType

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env file not loaded or missing")
	}
	serverLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	// initializarea bazei de date

	// initializing the fiber app and setting the view engine
	engine := mustache.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layout/main",
	})
	app.Static("/static", "./static")
	app.Static("/assets", "./build/client/assets")

	db, err := gorm.Open(sqlite.Open("./ettiContent.db"), &gorm.Config{})
	db.AutoMigrate(&types.Lab{}, &types.Blog{})
	if err != nil {
		// serverLogger.Println(err)
		serverLogger.Fatal(err)
	}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	// logging
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		// trebuie sa adaug prod si development aici
		AllowOrigins:     "http://localhost:5173,http://localhost:3000,https://ettiui.netlify.app",
		AllowCredentials: true,
	}))

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.SendString("LoginPage")
	})

	routes.RegisterApiRouter(app, serverLogger)
	// routes.RegisterAdminRoutes(app, serverLogger)
	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("./build/client/index.html")
	})
	app.Listen(":3000")
}
