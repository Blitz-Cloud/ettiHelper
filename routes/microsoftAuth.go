package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

func RegisterMicrosoftOAuth(app *fiber.App) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env file not loaded.\nCant create OAuth client without")
	}
	callbackURL := ""
	if os.Getenv("isProd") == "true" {
		callbackURL = "https://ettih.blitzcloud.me/auth/microsoft/callback"
	} else {
		callbackURL = "http://localhost:3000/auth/microsoft/callback"
	}

	oauth2Config := &oauth2.Config{
		ClientID:     os.Getenv("microsoftClientId"),
		ClientSecret: os.Getenv("microsoftClientSecret"),
		RedirectURL:  callbackURL,
		Scopes:       []string{"User.Read"},
		Endpoint:     microsoft.AzureADEndpoint("upb.ro"),
	}
	verifier := oauth2.GenerateVerifier()

	router := app.Group("/auth/etti")
	router.Get("/", func(c *fiber.Ctx) error {
		authUrl := oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
		return c.Redirect(authUrl)
	})

	router.Get("/callback", func(c *fiber.Ctx) error {
		code := c.Query("code")
		state := c.Query("state")

		// Verify the state matches
		spew.Dump(state)
		if state != "state" {
			return c.Status(http.StatusBadRequest).SendString("Invalid state parameter")
		}

		// Get token from the authorization code
		token, err := oauth2Config.Exchange(c.Context(), code, oauth2.VerifierOption(verifier))
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to exchange token")
		}

		// Use the token to get user info
		client := oauth2Config.Client(c.Context(), token)
		resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to fetch user data")
		}
		defer resp.Body.Close()

		// You can use the response to get the user's info
		// For example, print the body
		var userInfo map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to parse user data")
		}
		c.Cookie(&fiber.Cookie{
			Name:     "testC",
			Value:    "test",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			HTTPOnly: true,
		})
		return c.Redirect("/labs")

	})

}
