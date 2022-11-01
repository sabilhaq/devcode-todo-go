package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sabilhaq/devcode-todo-go/database"
	"github.com/sabilhaq/devcode-todo-go/models"
	"github.com/sabilhaq/devcode-todo-go/utils"
)

func GetActivities(c *fiber.Ctx) error {
	db := database.DBConn
	var activities []models.Activity
	db.Find(&activities)

	return c.JSON(models.Response{
		Status:  "Success",
		Message: "Success",
		Data:    activities,
	})
}

func CreateActivity(c *fiber.Ctx) error {
	activity := new(models.Activity)
	if err := c.BodyParser(&activity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
			Data:    map[string]interface{}{},
		})
	}

	errors := utils.ValidateStruct(*activity)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: fmt.Sprintf("%v cannot be null", errors[0].FailedField),
			Data:    map[string]interface{}{},
		})
	}

	db := database.DBConn
	db.Create(&activity)

	return c.JSON(models.Response{
		Status:  "Success",
		Message: "Success",
		Data:    activity,
	})
}

func GetActivity(c *fiber.Ctx) error {
	db := database.DBConn
	var activity models.Activity
	id, _ := strconv.Atoi(c.Params("id"))
	if err := db.First(&activity, id).Error; err != nil {
		return c.JSON(models.Response{
			Status:  http.StatusText(http.StatusNotFound),
			Message: fmt.Sprintf("Activity with ID %v Not Found", id),
			Data:    map[string]interface{}{},
		})
	}

	return c.JSON(models.Response{
		Status:  "Success",
		Message: "Success",
		Data:    activity,
	})
}

func UpdateActivity(c *fiber.Ctx) error {
	req := new(models.Activity)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
			Data:    map[string]interface{}{},
		})
	}

	errors := utils.ValidateStruct(*req)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: fmt.Sprintf("%v cannot be null", errors[0].FailedField),
			Data:    map[string]interface{}{},
		})
	}

	db := database.DBConn
	id, _ := strconv.Atoi(c.Params("id"))
	activity := new(models.Activity)
	if err := db.First(&activity, id).Error; err != nil {
		return c.JSON(models.Response{
			Status:  http.StatusText(http.StatusNotFound),
			Message: fmt.Sprintf("Activity with ID %v Not Found", id),
			Data:    map[string]interface{}{},
		})
	}

	activity.Title = req.Title
	if req.Email != "" {
		activity.Email = req.Email
	}

	db.Save(&activity)

	return c.JSON(models.Response{
		Status:  "Success",
		Message: "Success",
		Data:    activity,
	})
}

func DeleteActivity(c *fiber.Ctx) error {
	db := database.DBConn
	id, _ := strconv.Atoi(c.Params("id"))

	res := db.Delete(&models.Activity{}, id)
	if res.RowsAffected == 0 {
		return c.JSON(models.Response{
			Status:  http.StatusText(http.StatusNotFound),
			Message: fmt.Sprintf("Activity with ID %v Not Found", id),
			Data:    map[string]interface{}{},
		})
	}
	db.Where("activity_group_id = ?", id).Delete(&models.Todo{})

	return c.JSON(models.Response{
		Status:  "Success",
		Message: "Success",
		Data:    map[string]interface{}{},
	})
}
