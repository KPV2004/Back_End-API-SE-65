package middleware

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// Global Firebase app instance
var app *firebase.App

// Initialize Firebase Admin SDK
func InitFirebase() {
	env_err := godotenv.Load()
	if env_err != nil {
		// log.Fatal("Error loading .env file")
		panic("Error loading .env file")
	}

	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS"))
	var err error
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase: %v", err)
	}
}

// Verify Firebase ID Token
func VerifyIDToken(idToken string) (*auth.Token, error) {
	if app == nil {
		return nil, fmt.Errorf("Firebase not initialized")
	}

	ctx := context.Background()
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error getting Auth client: %v", err)
	}

	// Verify ID Token
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("Invalid or expired token: %v", err)
	}

	return token, nil
}

// AuthMiddleware - Protect routes by verifying Firebase ID token
func AuthMiddleware(c *fiber.Ctx) error {
	// Extract token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header is missing"})
	}

	// Ensure token has "Bearer " prefix
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}

	// Extract the token part
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Verify the ID token
	token, err := VerifyIDToken(tokenStr)
	if err != nil {
		log.Println("Token verification failed:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Store user ID in context
	c.Locals("userID", token.UID)

	// Proceed to the next middleware or handler
	return c.Next()
}
