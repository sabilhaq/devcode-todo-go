package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID              int            `json:"id" gorm:"primarykey"`
	ActivityGroupID int            `json:"activity_group_id" validate:"required"`
	Title           string         `json:"title" validate:"required"`
	IsActive        string         `json:"is_active"`
	Priority        string         `json:"priority"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type CreateTodoResponse struct {
	ID              int            `json:"id"`
	ActivityGroupID int            `json:"activity_group_id" validate:"required"`
	Title           string         `json:"title" validate:"required"`
	IsActive        bool           `json:"is_active"`
	Priority        string         `json:"priority"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at"`
}

type GetTodoResponse struct {
	ID              int            `json:"id"`
	ActivityGroupID string         `json:"activity_group_id"`
	Title           string         `json:"title"`
	IsActive        bool           `json:"is_active"`
	Priority        string         `json:"priority"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at"`
}
