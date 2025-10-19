package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	TelegramID int64          `json:"telegram_id" gorm:"uniqueIndex;not null"`
	Username   string         `json:"username"`
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	TeamID     *uint          `json:"team_id"`
	Team       *Team          `json:"team" gorm:"foreignKey:TeamID"`
	Roles      []Role         `json:"roles" gorm:"many2many:user_roles;"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
