package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	Activity struct {
		ID        int            `json:"id" gorm:"primarykey"`
		Email     string         `json:"email"`
		Title     string         `json:"title" validate:"required"`
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt time.Time      `json:"updated_at"`
		DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	}

	Todo struct {
		ID              int            `json:"id" gorm:"primarykey"`
		ActivityGroupID int            `json:"activity_group_id" validate:"required"`
		Title           string         `json:"title" validate:"required"`
		IsActive        string         `json:"is_active" default:"1"`
		Priority        string         `json:"priority" default:"very-high"`
		CreatedAt       time.Time      `json:"created_at"`
		UpdatedAt       time.Time      `json:"updated_at"`
		DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	}

	GetTodoResponse struct {
		ID              int            `json:"id" gorm:"primarykey"`
		ActivityGroupID string         `json:"activity_group_id" validate:"required"`
		Title           string         `json:"title" validate:"required"`
		IsActive        bool           `json:"is_active"`
		Priority        string         `json:"priority"`
		CreatedAt       time.Time      `json:"created_at"`
		UpdatedAt       time.Time      `json:"updated_at"`
		DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	}

	Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}

	ErrorResponse struct {
		FailedField string
		Tag         string
		Value       string
	}
)

var validate = validator.New()

func ValidateStruct(s interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			var errorField string
			if err.Field() == "Title" {
				errorField = "title"
			} else if err.Field() == "ActivityGroupID" {
				errorField = "activity_group_id"
			}

			element.FailedField = errorField
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local`, os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), "3306", os.Getenv("MYSQL_DBNAME"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to mysql database")
	}

	// Migrate
	db.AutoMigrate(
		&Activity{},
		&Todo{},
	)

	app := fiber.New()

	// Activity handler
	activityGroups := app.Group("/activity-groups")

	activityGroups.Get("", func(c *fiber.Ctx) error {
		activities := new([]Activity)
		db.Find(&activities)
		return c.Status(http.StatusOK).JSON(Response{
			Status:  "Success",
			Message: "Success",
			Data:    activities,
		})
	})

	activityGroups.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		activityDB := new(Activity)
		db.First(&activityDB, id)
		idInt, _ := strconv.Atoi(id)
		if int(activityDB.ID) != idInt {
			return c.Status(http.StatusNotFound).JSON(Response{
				Status:  "Not Found",
				Message: fmt.Sprintf("Activity with ID %v Not Found", id),
				Data:    map[string]any{},
			})
		}

		return c.Status(http.StatusOK).JSON(Response{
			Status:  "Success",
			Message: "Success",
			Data:    activityDB,
		})
	})

	activityGroups.Post("", func(c *fiber.Ctx) error {
		activity := new(Activity)
		if err := c.BodyParser(activity); err != nil {
			return err
		}
		if err := ValidateStruct(activity); err != nil {
			return c.Status(http.StatusBadRequest).JSON(Response{
				Status:  "Bad Request",
				Message: "title cannot be null",
				Data:    map[string]any{},
			})
		}

		db.Create(&activity)
		return c.Status(http.StatusCreated).JSON(Response{
			Status:  "Success",
			Message: "Success",
			Data:    activity,
		})
	})

	activityGroups.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		activityDB := new(Activity)
		db.First(&activityDB, id)
		idInt, _ := strconv.Atoi(id)
		if int(activityDB.ID) != idInt {
			return c.Status(http.StatusNotFound).JSON(Response{
				Status:  "Not Found",
				Message: fmt.Sprintf("Activity with ID %v Not Found", id),
				Data:    map[string]any{},
			})
		}

		db.Delete(&activityDB)
		db.Where("activity_id = ?", activityDB.ID).Delete(&Todo{})
		return c.Status(http.StatusOK).JSON(Response{
			Status:  "Success",
			Message: "Success",
			Data:    map[string]any{},
		})
	})

	activityGroups.Patch("/:id", func(c *fiber.Ctx) error {
		activity := new(Activity)
		if err := c.BodyParser(activity); err != nil {
			return err
		}
		if err := ValidateStruct(activity); err != nil {
			return c.Status(http.StatusBadRequest).JSON(Response{
				Status:  "Bad Request",
				Message: "title cannot be null",
				Data:    map[string]any{},
			})
		}

		id := c.Params("id")
		activityDB := new(Activity)
		db.First(&activityDB, id)
		idInt, _ := strconv.Atoi(id)
		if int(activityDB.ID) != idInt {
			return c.Status(http.StatusNotFound).JSON(Response{
				Status:  "Not Found",
				Message: fmt.Sprintf("Activity with ID %v Not Found", id),
				Data:    map[string]any{},
			})
		}

		activityDB.Title = activity.Title
		if activity.Email != "" {
			activityDB.Email = activity.Email
		}

		db.Save(&activityDB)
		return c.Status(http.StatusOK).JSON(Response{
			Status:  "Success",
			Message: "Success",
			Data:    activityDB,
		})
	})

	// Todo handler
	todoItems := app.Group("/todo-items")

	todoItems.Get("", func(c *fiber.Ctx) error {
		activityID := c.Query("activity_group_id")
		activityIDInt := 0
		todos := []Todo{}
		if activityID != "" {
			activityIDInt, _ = strconv.Atoi(activityID)
			db.Where("activity_id=?", activityIDInt).Find(&todos)
		} else {
			db.Find(&todos)
		}
		todosResp := []GetTodoResponse{}
		for _, todo := range todos {
			isActive := true
			if todo.IsActive != "1" {
				isActive = false
			}
			todosResp = append(todosResp, GetTodoResponse{
				ID:              todo.ID,
				ActivityGroupID: strconv.Itoa(todo.ActivityGroupID),
				Title:           todo.Title,
				IsActive:        isActive,
				Priority:        todo.Priority,
				CreatedAt:       todo.CreatedAt,
				UpdatedAt:       todo.UpdatedAt,
				DeletedAt:       todo.DeletedAt,
			})
		}

		return c.Status(http.StatusOK).JSON(Response{
			Status:  "Success",
			Message: "Success",
			Data:    todosResp,
		})
	})

	todoItems.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		todoDB := new(Todo)
		res := db.First(&todoDB, id)

		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return c.Status(http.StatusNotFound).JSON(Response{
					Status:  "Not Found",
					Message: fmt.Sprintf("Todo with ID %v Not Found", id),
					Data:    map[string]any{},
				})
			}
		}
		todoResp := GetTodoResponse{}
		isActive := true
		if todoDB.IsActive != "1" {
			isActive = false
		}
		todoResp = GetTodoResponse{
			ID:              todoDB.ID,
			ActivityGroupID: strconv.Itoa(todoDB.ActivityGroupID),
			Title:           todoDB.Title,
			IsActive:        isActive,
			Priority:        todoDB.Priority,
			CreatedAt:       todoDB.CreatedAt,
			UpdatedAt:       todoDB.UpdatedAt,
			DeletedAt:       todoDB.DeletedAt,
		}
		return c.Status(http.StatusOK).JSON(Response{
			Status:  "Success",
			Message: "Success",
			Data:    todoResp,
		})
	})

	todoItems.Post("", func(c *fiber.Ctx) error {
		todo := new(Todo)
		if err := c.BodyParser(todo); err != nil {
			return err
		}
		if err := ValidateStruct(todo); err != nil {
			return c.Status(http.StatusBadRequest).JSON(Response{
				Status:  "Bad Request",
				Message: fmt.Sprintf("%v cannot be null", err[0].FailedField),
				Data:    map[string]any{},
			})
		}

		todo.IsActive = "1"
		todo.Priority = "very-high"

		result := db.Create(&todo)
		if result.Error != nil {
			return c.Status(http.StatusBadRequest).JSON(Response{
				Status:  "Bad Request",
				Message: "activity_group_id not found",
				Data:    map[string]any{},
			})
		}

		type createTodoResponse struct {
			ID              int            `json:"id"`
			ActivityGroupID int            `json:"activity_group_id" validate:"required"`
			Title           string         `json:"title" validate:"required"`
			IsActive        bool           `json:"is_active" default:"1"`
			Priority        string         `json:"priority" default:"very-high"`
			CreatedAt       time.Time      `json:"created_at"`
			UpdatedAt       time.Time      `json:"updated_at"`
			DeletedAt       gorm.DeletedAt `json:"deleted_at"`
		}
		todoResp := createTodoResponse{}
		isActive := false
		if todo.IsActive == "1" {
			isActive = true
		}
		todoResp = createTodoResponse{
			ID:              todo.ID,
			ActivityGroupID: todo.ActivityGroupID,
			Title:           todo.Title,
			IsActive:        isActive,
			Priority:        todo.Priority,
			CreatedAt:       todo.CreatedAt,
			UpdatedAt:       todo.UpdatedAt,
			DeletedAt:       todo.DeletedAt,
		}
		return c.Status(http.StatusCreated).JSON(Response{
			Status:  "Success",
			Message: "Success",
			Data:    todoResp,
		})
	})

	todoItems.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		todoDB := new(Todo)
		db.First(&todoDB, id)
		idInt, _ := strconv.Atoi(id)
		if int(todoDB.ID) != idInt {
			return c.Status(http.StatusNotFound).JSON(Response{
				Status:  "Not Found",
				Message: fmt.Sprintf("Todo with ID %v Not Found", id),
				Data:    map[string]any{},
			})
		}

		db.Delete(&todoDB)
		return c.Status(http.StatusOK).JSON(Response{
			Status:  "Success",
			Message: "Success",
			Data:    map[string]any{},
		})
	})

	todoItems.Patch("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		todoDB := new(Todo)
		res := db.First(&todoDB, id)
		if res.Error != nil {
			return c.Status(http.StatusNotFound).JSON(Response{
				Status:  "Not Found",
				Message: fmt.Sprintf("Todo with ID %v Not Found", id),
				Data:    map[string]any{},
			})
		}
		if strconv.Itoa(todoDB.ID) != id {
			return c.Status(http.StatusNotFound).JSON(Response{
				Status:  "Not Found",
				Message: fmt.Sprintf("Todo with ID %v Not Found", id),
				Data:    map[string]any{},
			})
		}

		todo := new(Todo)
		_ = c.BodyParser(todo)

		if todo.ActivityGroupID != 0 {
			todoDB.ActivityGroupID = todo.ActivityGroupID
		}
		if todo.IsActive != "" {
			todoDB.IsActive = todo.IsActive
		}
		if todo.Priority != "" {
			todoDB.Priority = todo.Priority
		}

		if todo.ActivityGroupID == 0 {
			todo.ActivityGroupID = todoDB.ActivityGroupID
		}

		if todo.Title != "" {
			todoDB.Title = todo.Title
		}

		db.Save(todoDB)
		return c.Status(http.StatusOK).JSON(Response{
			Status:  "Success",
			Message: "Success",
			Data:    todoDB,
		})
	})

	log.Fatal(app.Listen(":3030"))
}
