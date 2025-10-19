package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	CreatedBy   uint   `json:"created_by"`
	// Creator     User           `json:"creator" gorm:"foreignKey:CreatedBy"`
	Members   []User         `json:"members" gorm:"foreignKey:TeamID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
