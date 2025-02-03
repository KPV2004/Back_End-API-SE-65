package main

import (
	"fmt"
	"go-server/adapters"
	"go-server/core"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
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

	//Implement Port Hexagonal Arc {Secondary to Primary Port}
	userRepo := adapters.NewGormUserRepository(db)
	userService := core.NewUserService(userRepo)
	userHandler := adapters.NewHttpUserHandler(userService)

	app.Post("/user/register", userHandler.RegisterUser)
	app.Get("/user/getuser/:email", userHandler.GetUser)
	// Migrate the schema
	db.AutoMigrate(&core.User{})
	fmt.Println("Database migration completed!")
	app.Listen((":8000"))
	// newBook := &Book{Name: "Think Again", Author: "adam", Description: "test", price: 200}

	// createBook(db, newBook)
}
