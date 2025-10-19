package seeders

import (
	"log"
	"valorant-app/models"

	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) {
	log.Println("Starting database seeding...")

	seedPermissions(db)
	seedDefaultRoles(db)

	log.Println("Database seeding completed!")
}

func seedPermissions(db *gorm.DB) {
	permissions := []models.Permission{
		{Name: models.PermissionManageTeam, Description: "Управление командой"},
		{Name: models.PermissionManageRoles, Description: "Управление ролями"},
		{Name: models.PermissionKickMembers, Description: "Исключение участников"},
		{Name: models.PermissionInviteMembers, Description: "Приглашение участников"},
		{Name: models.PermissionViewMembers, Description: "Просмотр участников"},
		{Name: models.PermissionEditProfile, Description: "Редактирование профиля"},
	}

	for _, permission := range permissions {
		var existingPermission models.Permission
		result := db.Where("name = ?", permission.Name).First(&existingPermission)
		if result.Error != nil {
			// Permission doesn't exist, create it
			db.Create(&permission)
			log.Printf("Created permission: %s", permission.Name)
		}
	}
}

func seedDefaultRoles(db *gorm.DB) {
	// Создаем базовые права для каждой роли
	ownerPermissions := []string{
		models.PermissionManageTeam,
		models.PermissionManageRoles,
		models.PermissionKickMembers,
		models.PermissionInviteMembers,
		models.PermissionViewMembers,
		models.PermissionEditProfile,
	}

	adminPermissions := []string{
		models.PermissionManageTeam,
		models.PermissionKickMembers,
		models.PermissionInviteMembers,
		models.PermissionViewMembers,
		models.PermissionEditProfile,
	}

	captainPermissions := []string{
		models.PermissionKickMembers,
		models.PermissionInviteMembers,
		models.PermissionViewMembers,
		models.PermissionEditProfile,
	}

	memberPermissions := []string{
		models.PermissionViewMembers,
		models.PermissionEditProfile,
	}

	// Создаем роли с правами
	createRoleWithPermissions(db, models.RoleOwner, "Владелец команды", "Полный доступ ко всем функциям команды", ownerPermissions)
	createRoleWithPermissions(db, models.RoleAdmin, "Администратор", "Управление командой и участниками", adminPermissions)
	createRoleWithPermissions(db, models.RoleCaptain, "Капитан", "Управление участниками команды", captainPermissions)
	createRoleWithPermissions(db, models.RoleMember, "Участник", "Базовые права участника", memberPermissions)
}

func createRoleWithPermissions(db *gorm.DB, roleName, description, fullDescription string, permissionNames []string) {
	// Проверяем, существует ли роль
	var existingRole models.Role
	result := db.Where("name = ?", roleName).First(&existingRole)
	if result.Error == nil {
		// Роль уже существует, обновляем права
		updateRolePermissions(db, &existingRole, permissionNames)
		return
	}

	// Создаем новую роль
	role := models.Role{
		Name:        roleName,
		Description: description,
	}

	// Находим права по именам
	var permissions []models.Permission
	db.Where("name IN ?", permissionNames).Find(&permissions)
	role.Permissions = permissions

	db.Create(&role)
	log.Printf("Created role: %s", roleName)
}

func updateRolePermissions(db *gorm.DB, role *models.Role, permissionNames []string) {
	// Находим права по именам
	var permissions []models.Permission
	db.Where("name IN ?", permissionNames).Find(&permissions)

	// Обновляем права роли
	db.Model(role).Association("Permissions").Replace(permissions)
	log.Printf("Updated permissions for role: %s", role.Name)
}
