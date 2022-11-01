package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/mattn/go-colorable"
	"github.com/sabilhaq/devcode-todo-go/database"
	"github.com/sabilhaq/devcode-todo-go/handler"
	"github.com/sabilhaq/devcode-todo-go/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var lw *log.Logger = log.New(os.Stderr, "\r\n", log.LstdFlags)
	if runtime.GOOS == "windows" {
		lw.SetOutput(colorable.NewColorableStderr())
	}

	newLogger := logger.New(
		lw,
		logger.Config{
			LogLevel: logger.Error, // Log level
			Colorful: true,         // Disable color
		},
	)

	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local`, os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), "3306", os.Getenv("MYSQL_DBNAME"))
	database.DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	database.DBConn.AutoMigrate(
		&models.Activity{},
		&models.Todo{},
	)
}

func setupRoutes(app *fiber.App) {
	app.Use(fiberLogger.New())

	app.Get("/activity-groups", handler.GetActivities)
	app.Post("/activity-groups", handler.CreateActivity)
	app.Get("/activity-groups/:id", handler.GetActivity)
	app.Patch("/activity-groups/:id", handler.UpdateActivity)
	app.Delete("/activity-groups/:id", handler.DeleteActivity)

	app.Get("/todo-items", handler.GetTodos)
	app.Post("/todo-items", handler.CreateTodo)
	app.Get("/todo-items/:id", handler.GetTodo)
	app.Patch("/todo-items/:id", handler.UpdateTodo)
	app.Delete("/todo-items/:id", handler.DeleteTodo)
}

func main() {
	app := fiber.New()
	initDatabase()
	setupRoutes(app)

	log.Fatal(app.Listen(":3030"))
}
