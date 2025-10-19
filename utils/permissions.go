package utils

import (
	"valorant-app/database"
	"valorant-app/models"
)

// CheckPermission проверяет, есть ли у пользователя определенное право в команде
func CheckPermission(userID, teamID uint, permissionName string) bool {
	var user models.User
	result := database.DB.Preload("Roles.Permissions").Where("id = ? AND team_id = ?", userID, teamID).First(&user)
	if result.Error != nil {
		return false
	}

	for _, role := range user.Roles {
		if role.TeamID != teamID {
			continue
		}
		for _, permission := range role.Permissions {
			if permission.Name == permissionName {
				return true
			}
		}
	}

	return false
}

// HasRole проверяет, есть ли у пользователя определенная роль в команде
func HasRole(userID, teamID uint, roleName string) bool {
	var user models.User
	result := database.DB.Preload("Roles").Where("id = ? AND team_id = ?", userID, teamID).First(&user)
	if result.Error != nil {
		return false
	}

	for _, role := range user.Roles {
		if role.TeamID == teamID && role.Name == roleName {
			return true
		}
	}

	return false
}

// IsTeamOwner проверяет, является ли пользователь владельцем команды
func IsTeamOwner(userID, teamID uint) bool {
	var team models.Team
	result := database.DB.Where("id = ? AND created_by = ?", teamID, userID).First(&team)
	return result.Error == nil
}

// GetUserPermissions получает все права пользователя в команде
func GetUserPermissions(userID, teamID uint) []string {
	var user models.User
	result := database.DB.Preload("Roles.Permissions").Where("id = ? AND team_id = ?", userID, teamID).First(&user)
	if result.Error != nil {
		return []string{}
	}

	permissions := make(map[string]bool)
	for _, role := range user.Roles {
		if role.TeamID != teamID {
			continue
		}
		for _, permission := range role.Permissions {
			permissions[permission.Name] = true
		}
	}

	var resultPermissions []string
	for permission := range permissions {
		resultPermissions = append(resultPermissions, permission)
	}

	return resultPermissions
}
