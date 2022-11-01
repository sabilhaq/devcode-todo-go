package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sabilhaq/devcode-todo-go/database"
	"github.com/sabilhaq/devcode-todo-go/models"
	"github.com/sabilhaq/devcode-todo-go/utils"
)

func GetTodos(c *gin.Context) {
	activityID, _ := strconv.Atoi(c.Query("activity_group_id"))

	db := database.DBConn
	if activityID != 0 {
		db = db.Where("activity_group_id = ?", activityID)
	}
	var result []models.Todo
	db.Find(&result)

	todos := []models.GetTodoResponse{}
	for _, todo := range result {
		var isActive bool
		if todo.IsActive == "1" {
			isActive = true
		}
		todos = append(todos, models.GetTodoResponse{
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

	c.JSON(http.StatusOK, models.Response{
		Status:  "Success",
		Message: "Success",
		Data:    todos,
	})
}

func CreateTodo(c *gin.Context) {
	todo := new(models.Todo)
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
			Data:    map[string]interface{}{},
		})
	}

	errors := utils.ValidateStruct(*todo)
	if errors != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: fmt.Sprintf("%v cannot be null", errors[0].FailedField),
			Data:    map[string]interface{}{},
		})
	}

	db := database.DBConn
	todo.IsActive = "1"
	todo.Priority = "very-high"

	db.Create(&todo)

	c.JSON(http.StatusCreated, models.Response{
		Status:  "Success",
		Message: "Success",
		Data: models.CreateTodoResponse{
			ID:              todo.ID,
			ActivityGroupID: todo.ActivityGroupID,
			Title:           todo.Title,
			IsActive:        true,
			Priority:        todo.Priority,
			CreatedAt:       todo.CreatedAt,
			UpdatedAt:       todo.UpdatedAt,
			DeletedAt:       todo.DeletedAt,
		},
	})
}

func GetTodo(c *gin.Context) {
	db := database.DBConn
	var todo models.Todo
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Status:  http.StatusText(http.StatusNotFound),
			Message: fmt.Sprintf("Todo with ID %v Not Found", id),
			Data:    map[string]interface{}{},
		})
	}

	var isActive bool
	if todo.IsActive == "1" {
		isActive = true
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  "Success",
		Message: "Success",
		Data: models.GetTodoResponse{
			ID:              todo.ID,
			ActivityGroupID: strconv.Itoa(todo.ActivityGroupID),
			Title:           todo.Title,
			IsActive:        isActive,
			Priority:        todo.Priority,
			CreatedAt:       todo.CreatedAt,
			UpdatedAt:       todo.UpdatedAt,
			DeletedAt:       todo.DeletedAt,
		},
	})
}

func UpdateTodo(c *gin.Context) {
	req := new(models.Todo)
	_ = c.ShouldBindJSON(&req)

	db := database.DBConn
	id, _ := strconv.Atoi(c.Param("id"))
	todo := new(models.Todo)
	if err := db.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Status:  http.StatusText(http.StatusNotFound),
			Message: fmt.Sprintf("Todo with ID %v Not Found", id),
			Data:    map[string]interface{}{},
		})
	}

	if req.Title != "" {
		todo.Title = req.Title
	}
	if req.ActivityGroupID != 0 {
		todo.ActivityGroupID = req.ActivityGroupID
	}
	if req.IsActive != "" {
		todo.IsActive = req.IsActive
	}
	if req.Priority != "" {
		todo.Priority = req.Priority
	}

	db.Save(&todo)

	c.JSON(http.StatusOK, models.Response{
		Status:  "Success",
		Message: "Success",
		Data:    todo,
	})
}

func DeleteTodo(c *gin.Context) {
	db := database.DBConn
	id, _ := strconv.Atoi(c.Param("id"))

	res := db.Delete(&models.Todo{}, id)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.Response{
			Status:  http.StatusText(http.StatusNotFound),
			Message: fmt.Sprintf("Todo with ID %v Not Found", id),
			Data:    map[string]interface{}{},
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  "Success",
		Message: "Success",
		Data:    map[string]interface{}{},
	})
}
