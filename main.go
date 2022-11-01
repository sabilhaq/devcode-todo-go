package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
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

// func setupRoutes(app *fiber.App) {
func setupRoutes(app *gin.Engine) {
	// app.Use(fiberLogger.New())

	app.GET("/activity-groups", handler.GetActivities)
	app.POST("/activity-groups", handler.CreateActivity)
	app.GET("/activity-groups/:id", handler.GetActivity)
	app.PATCH("/activity-groups/:id", handler.UpdateActivity)
	app.DELETE("/activity-groups/:id", handler.DeleteActivity)

	app.GET("/todo-items", handler.GetTodos)
	app.POST("/todo-items", handler.CreateTodo)
	app.GET("/todo-items/:id", handler.GetTodo)
	app.PATCH("/todo-items/:id", handler.UpdateTodo)
	app.DELETE("/todo-items/:id", handler.DeleteTodo)
}

func main() {
	app := gin.Default()

	initDatabase()
	setupRoutes(app)

	app.Run(":3030") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
