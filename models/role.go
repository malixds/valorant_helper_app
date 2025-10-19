package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	TeamID      uint           `json:"team_id" gorm:"not null"`
	Team        Team           `json:"team" gorm:"foreignKey:TeamID"`
	Permissions []Permission   `json:"permissions" gorm:"many2many:role_permissions;"`
	Users       []User         `json:"users" gorm:"many2many:user_roles;"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"unique;not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Предопределенные роли
const (
	RoleOwner   = "owner"   // Владелец команды
	RoleAdmin   = "admin"   // Администратор
	RoleCaptain = "captain" // Капитан
	RoleMember  = "member"  // Участник
)

// Предопределенные права
const (
	PermissionManageTeam    = "manage_team"    // Управление командой
	PermissionManageRoles   = "manage_roles"   // Управление ролями
	PermissionKickMembers   = "kick_members"   // Исключение участников
	PermissionInviteMembers = "invite_members" // Приглашение участников
	PermissionViewMembers   = "view_members"   // Просмотр участников
	PermissionEditProfile   = "edit_profile"   // Редактирование профиля
)
