package handlers

import (
	"net/http"
	"strconv"
	"valorant-app/database"
	"valorant-app/models"

	"github.com/gin-gonic/gin"
)

// GetTeamRoles получает все роли команды
func GetTeamRoles(c *gin.Context) {
	teamIDStr := c.Param("team_id")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var roles []models.Role
	result := database.DB.Preload("Permissions").Where("team_id = ?", teamID).Find(&roles)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// CreateRole создает новую роль в команде
func CreateRole(c *gin.Context) {
	teamIDStr := c.Param("team_id")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role.TeamID = uint(teamID)
	result := database.DB.Create(&role)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// AssignRole назначает роль пользователю
func AssignRole(c *gin.Context) {
	teamIDStr := c.Param("team_id")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	roleIDStr := c.Param("role_id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	// Проверяем, что роль принадлежит команде
	var role models.Role
	result := database.DB.Where("id = ? AND team_id = ?", roleID, teamID).First(&role)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found in this team"})
		return
	}

	// Проверяем, что пользователь в команде
	var user models.User
	result = database.DB.Where("id = ? AND team_id = ?", userID, teamID).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in this team"})
		return
	}

	// Назначаем роль
	err = database.DB.Model(&user).Association("Roles").Append(&role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully"})
}

// RemoveRole убирает роль у пользователя
func RemoveRole(c *gin.Context) {
	teamIDStr := c.Param("team_id")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	roleIDStr := c.Param("role_id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var user models.User
	result := database.DB.Where("id = ? AND team_id = ?", userID, teamID).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in this team"})
		return
	}

	var role models.Role
	result = database.DB.Where("id = ? AND team_id = ?", roleID, teamID).First(&role)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found in this team"})
		return
	}

	// Убираем роль
	err = database.DB.Model(&user).Association("Roles").Delete(&role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role removed successfully"})
}

// GetUserRoles получает роли пользователя в команде
func GetUserRoles(c *gin.Context) {
	teamIDStr := c.Param("team_id")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	result := database.DB.Preload("Roles", "team_id = ?", teamID).Where("id = ? AND team_id = ?", userID, teamID).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in this team"})
		return
	}

	c.JSON(http.StatusOK, user.Roles)
}
