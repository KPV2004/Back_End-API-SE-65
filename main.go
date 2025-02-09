package main

import (
	"fmt"
	"go-server/adapters"
	"go-server/core"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

func main() {

	app := fiber.New()
	// Configure your PostgreSQL database details here

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Change this to your frontend domain for security
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // add Logger
	})

	if err != nil {
		panic("failed to connect to database")
	}

	env_err := godotenv.Load()
	if env_err != nil {
		// log.Fatal("Error loading .env file")
		panic("Error loading .env file")
	}

	smtpHost := os.Getenv("MAILER_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("MAILER_PORT")) // Convert port to int
	smtpUser := os.Getenv("MAILER_USERNAME")
	smtpPass := os.Getenv("MAILER_PASSWORD")

	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	//Implement Port Hexagonal Arc {Secondary to Primary Port}
	userRepo := adapters.NewGormUserRepository(db)
	userService := core.NewUserService(userRepo)
	// userHandler := adapters.NewHttpUserHandler(userService)
	emailRepo := adapters.NewEmailRepository(dialer)
	emailService := core.NewEmailService(emailRepo)
	// Assuming core.NewEmailService(emailRepo) creates an email service
	userHandler := adapters.NewHttpUserHandler(userService, emailService)

	app.Post("/user/register", userHandler.RegisterUser)
	app.Get("/user/getuser/:email", userHandler.GetUser)
	app.Get("/user/genotp/:email", userHandler.GenOTP)
	app.Post("/user/verifyotp", userHandler.VerifyOTP)

	app.Post("/admin/register", userHandler.RegisterAdmin)
	app.Post("/admin/login", userHandler.LoginAdmin)

	// Migrate the schema
	db.AutoMigrate(&core.User{})
	db.AutoMigrate(&core.Admin{})
	db.AutoMigrate(&core.Verification{})
	fmt.Println("Database migration completed!")
	app.Listen(("0.0.0.0:8000"))
	// newBook := &Book{Name: "Think Again", Author: "adam", Description: "test", price: 200}

	// createBook(db, newBook)
}
