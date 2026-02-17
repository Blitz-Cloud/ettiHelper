package middleware

import (
	"fmt"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Global JWKS provider to avoid re-fetching on every request
var kf keyfunc.Keyfunc

func InitAuth(tenantID string) {
	// Microsoft's JWKS endpoint
	jwksURL := fmt.Sprintf("https://login.microsoftonline.com/%s/discovery/v2.0/keys", tenantID)

	var err error
	// Fetch and automatically refresh keys in the background every hour
	kf, err = keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch Microsoft JWKS: %v", err))
	}
}

func AuthMiddleware(tenantID string, clientID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and Validate
		token, err := jwt.Parse(tokenString, kf.Keyfunc,
			jwt.WithValidMethods([]string{"RS256"}), // Microsoft uses RS256
			jwt.WithAudience(clientID),              // Must match your App Client ID
			jwt.WithIssuer(fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", tenantID)),
		)

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
		}

		// Double-check Tenant ID (tid) to ensure user is in the correct Org
		if claims["tid"] != tenantID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Wrong organization"})
		}

		// Success: Store user info in context for later use in routes
		c.Locals("user_email", claims["preferred_username"])
		c.Locals("user_oid", claims["oid"]) // Unique ID for the user in AD

		return c.Next()
	}
}
