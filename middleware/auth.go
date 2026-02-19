package middleware

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func IsValidMSJWT(token string) bool {
	ctx := context.Background()
	keys, err := keyfunc.NewDefaultCtx(ctx, []string{"https://login.microsoftonline.com/upb.ro/discovery/keys"})
	// spew.Dump(keys)
	if err != nil {
		log.Fatalf("Failed to create a keyfunc.Keyfunc from the server's URL.\nError: %s", err)
	}
	parsedToken, err := jwt.Parse(token, keys.Keyfunc)
	if err != nil {
		log.Fatalf("Failed to create a keyfunc.Keyfunc from the server's URL.\nError: %s", err)
	}
	if parsedToken.Valid {
		return true
	}
	return false
}

func RouteProtector(c *fiber.Ctx) error {
	// tenant, _ := c.Locals("tenant").(types.Namespace)
	// if tenant.AuthFlow == "none" {
	// 	return c.Next()
	// }
	fmt.Println("Route protector hit")
	date := c.Cookies("testC")
	if date != "" {
		fmt.Println("Access Allowed")
		return c.Next()
	}
	return c.SendStatus(401)
	// return c.Redirect("/login")
}

func AdminRouteProtector(c *fiber.Ctx) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}
	if c.Get("Authorization") == "" || c.Get("Authorization") != fmt.Sprintf("Bearer %s", os.Getenv("adminToken")) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}

func ValidateJwtMiddleware(c *fiber.Ctx) error {
	fmt.Println("JWT Verification Hit")
	authorizationHeader := c.GetReqHeaders()["Authorization"]
	if len(authorizationHeader) == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	idToken := strings.Split(authorizationHeader[0], " ")
	if len(idToken) == 1 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if IsValidMSJWT(idToken[1]) {
		return c.Next()
	}

	return c.Redirect(fmt.Sprintf("/login?from=%s", c.OriginalURL()))
}
