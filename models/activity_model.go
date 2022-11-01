package models

import (
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	ID        int            `json:"id" gorm:"primarykey"`
	Email     string         `json:"email"`
	Title     string         `json:"title" validate:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
