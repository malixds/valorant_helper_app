package handlers

import (
	"net/http"
	"strconv"
	"valorant-app/database"
	"valorant-app/models"

	"github.com/gin-gonic/gin"
)

func GetTeams(c *gin.Context) {
	var teams []models.Team
	result := database.DB.Preload("Members").Find(&teams)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}

	c.JSON(http.StatusOK, teams)
}

func GetTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var team models.Team
	result := database.DB.Preload("Members").First(&team, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	c.JSON(http.StatusOK, team)
}

func CreateTeam(c *gin.Context) {
	var team models.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&team)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team"})
		return
	}

	// Создаем роль владельца для команды
	ownerRole := models.Role{
		Name:        models.RoleOwner,
		Description: "Владелец команды",
		TeamID:      team.ID,
	}

	// Находим все права для роли владельца
	var permissions []models.Permission
	database.DB.Find(&permissions)
	ownerRole.Permissions = permissions

	database.DB.Create(&ownerRole)

	// Назначаем роль владельца создателю команды
	if team.CreatedBy != 0 {
		var user models.User
		database.DB.First(&user, team.CreatedBy)
		database.DB.Model(&user).Association("Roles").Append(&ownerRole)
	}

	c.JSON(http.StatusCreated, team)
}

func JoinTeam(c *gin.Context) {
	telegramIDStr := c.Param("telegram_id")
	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid telegram ID"})
		return
	}

	teamIDStr := c.Param("team_id")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	// Check if user exists
	var user models.User
	result := database.DB.Where("telegram_id = ?", telegramID).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if team exists
	var team models.Team
	result = database.DB.First(&team, teamID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	// Update user's team
	user.TeamID = &team.ID
	result = database.DB.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join team"})
		return
	}

	// Находим или создаем роль участника для команды
	var memberRole models.Role
	result = database.DB.Where("name = ? AND team_id = ?", models.RoleMember, team.ID).First(&memberRole)
	if result.Error != nil {
		// Создаем роль участника для команды
		memberRole = models.Role{
			Name:        models.RoleMember,
			Description: "Участник команды",
			TeamID:      team.ID,
		}

		// Находим права для роли участника
		var permissions []models.Permission
		database.DB.Where("name IN ?", []string{models.PermissionViewMembers, models.PermissionEditProfile}).Find(&permissions)
		memberRole.Permissions = permissions

		database.DB.Create(&memberRole)
	}

	// Назначаем роль участника
	database.DB.Model(&user).Association("Roles").Append(&memberRole)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined team"})
}

func LeaveTeam(c *gin.Context) {
	telegramIDStr := c.Param("telegram_id")
	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid telegram ID"})
		return
	}

	var user models.User
	result := database.DB.Where("telegram_id = ?", telegramID).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.TeamID = nil
	result = database.DB.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave team"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully left team"})
}
