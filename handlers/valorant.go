package handlers

import (
	"net/http"
	"strconv"
	"valorant-app/database"
	"valorant-app/models"

	"github.com/gin-gonic/gin"
)

// AddValorantPlayer добавляет Valorant аккаунт к пользователю
func AddValorantPlayer(c *gin.Context) {
	telegramIDStr := c.Param("telegram_id")
	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid telegram ID"})
		return
	}

	var request struct {
		GameName string `json:"game_name" binding:"required"`
		Tag      string `json:"tag" binding:"required"`
		Region   string `json:"region" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Находим пользователя
	var user models.User
	result := database.DB.Where("telegram_id = ?", telegramID).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Создаем Valorant игрока
	valorantPlayer := models.ValorantPlayer{
		UserID:   user.ID,
		GameName: request.GameName,
		Tag:      request.Tag,
		Region:   request.Region,
	}

	result = database.DB.Create(&valorantPlayer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Valorant player"})
		return
	}

	c.JSON(http.StatusCreated, valorantPlayer)
}

// GetValorantPlayer получает Valorant аккаунт пользователя
func GetValorantPlayer(c *gin.Context) {
	telegramIDStr := c.Param("telegram_id")
	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid telegram ID"})
		return
	}

	var user models.User
	result := database.DB.Preload("ValorantPlayers").Where("telegram_id = ?", telegramID).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ValorantPlayers)
}

// GetTeamValorantPlayers получает всех Valorant игроков команды
func GetTeamValorantPlayers(c *gin.Context) {
	teamIDStr := c.Param("team_id")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var players []models.ValorantPlayer
	result := database.DB.Preload("User").Where("user_id IN (SELECT id FROM users WHERE team_id = ?)", teamID).Find(&players)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team players"})
		return
	}

	c.JSON(http.StatusOK, players)
}

// SyncPlayerData синхронизирует данные игрока с Valorant API
func SyncPlayerData(c *gin.Context) {
	telegramIDStr := c.Param("telegram_id")
	_, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid telegram ID"})
		return
	}

	// TODO: Реализовать синхронизацию с Valorant API
	// 1. Получить данные игрока из API
	// 2. Обновить данные в базе данных
	// 3. Получить матчи игрока
	// 4. Сохранить статистику

	c.JSON(http.StatusOK, gin.H{"message": "Player data synced successfully"})
}

// GetPlayerStats получает статистику игрока
func GetPlayerStats(c *gin.Context) {
	telegramIDStr := c.Param("telegram_id")
	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid telegram ID"})
		return
	}

	var user models.User
	result := database.DB.Preload("ValorantPlayers.Stats").Where("telegram_id = ?", telegramID).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ValorantPlayers)
}

// GetTeamStats получает статистику команды
func GetTeamStats(c *gin.Context) {
	teamIDStr := c.Param("team_id")
	_, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	// TODO: Реализовать агрегацию статистики команды
	// 1. Получить всех игроков команды
	// 2. Агрегировать их статистику
	// 3. Вернуть общую статистику команды

	c.JSON(http.StatusOK, gin.H{"message": "Team stats endpoint - to be implemented"})
}
