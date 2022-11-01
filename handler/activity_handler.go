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

func GetActivities(c *gin.Context) {
	db := database.DBConn
	var activities []models.Activity
	db.Find(&activities)

	c.JSON(http.StatusOK, models.Response{
		Status:  "Success",
		Message: "Success",
		Data:    activities,
	})
}

func CreateActivity(c *gin.Context) {
	activity := new(models.Activity)
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
			Data:    map[string]interface{}{},
		})
	}

	errors := utils.ValidateStruct(*activity)
	if errors != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: fmt.Sprintf("%v cannot be null", errors[0].FailedField),
			Data:    map[string]interface{}{},
		})
	}

	db := database.DBConn
	db.Create(&activity)

	c.JSON(http.StatusCreated, models.Response{
		Status:  "Success",
		Message: "Success",
		Data:    activity,
	})
}

func GetActivity(c *gin.Context) {
	db := database.DBConn
	var activity models.Activity
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Status:  http.StatusText(http.StatusNotFound),
			Message: fmt.Sprintf("Activity with ID %v Not Found", id),
			Data:    map[string]interface{}{},
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  "Success",
		Message: "Success",
		Data:    activity,
	})
}

func UpdateActivity(c *gin.Context) {
	req := new(models.Activity)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
			Data:    map[string]interface{}{},
		})
	}

	errors := utils.ValidateStruct(*req)
	if errors != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: fmt.Sprintf("%v cannot be null", errors[0].FailedField),
			Data:    map[string]interface{}{},
		})
	}

	db := database.DBConn
	id, _ := strconv.Atoi(c.Param("id"))
	activity := new(models.Activity)
	if err := db.First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
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

	c.JSON(http.StatusOK, models.Response{
		Status:  "Success",
		Message: "Success",
		Data:    activity,
	})
}

func DeleteActivity(c *gin.Context) {
	db := database.DBConn
	id, _ := strconv.Atoi(c.Param("id"))

	res := db.Delete(&models.Activity{}, id)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.Response{
			Status:  http.StatusText(http.StatusNotFound),
			Message: fmt.Sprintf("Activity with ID %v Not Found", id),
			Data:    map[string]interface{}{},
		})
	}
	db.Where("activity_group_id = ?", id).Delete(&models.Todo{})

	c.JSON(http.StatusOK, models.Response{
		Status:  "Success",
		Message: "Success",
		Data:    map[string]interface{}{},
	})
}
