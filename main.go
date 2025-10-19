package main

import (
	"log"
	"valorant-app/bot"
	"valorant-app/config"
	"valorant-app/database"
	"valorant-app/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database.InitDB(cfg)

	// Initialize bot
	telegramBot, err := bot.NewBot(cfg)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	// Check if we should use webhook or polling
	useWebhook := cfg.WebhookURL != "" && cfg.WebhookURL != "http://localhost:8080"
	useNgrok := cfg.NgrokURL != ""

	if useNgrok {
		// Use ngrok URL for webhook
		webhookURL := cfg.NgrokURL + "/webhook"
		if err := telegramBot.SetWebhook(webhookURL); err != nil {
			log.Fatal("Failed to set webhook:", err)
		}
		log.Println("Bot configured with ngrok webhook:", webhookURL)
		useWebhook = true
	} else if useWebhook {
		// Set webhook for production
		if err := telegramBot.SetWebhook(cfg.WebhookURL + "/webhook"); err != nil {
			log.Fatal("Failed to set webhook:", err)
		}
		log.Println("Bot configured with webhook:", cfg.WebhookURL+"/webhook")
	} else {
		// Use polling for local development
		log.Println("Bot configured for local development (polling mode)")
	}

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := r.Group("/api")
	{
		// User routes
		api.GET("/users/:telegram_id", handlers.GetUser)
		api.POST("/users", handlers.CreateUser)
		api.PUT("/users/:telegram_id", handlers.UpdateUser)

		// Team routes
		api.GET("/teams", handlers.GetTeams)
		api.GET("/teams/:id", handlers.GetTeam)
		api.POST("/teams", handlers.CreateTeam)
		api.POST("/teams/:team_id/join/:telegram_id", handlers.JoinTeam)
		api.POST("/teams/leave/:telegram_id", handlers.LeaveTeam)

		// Role routes
		api.GET("/teams/:team_id/roles", handlers.GetTeamRoles)
		api.POST("/teams/:team_id/roles", handlers.CreateRole)
		api.POST("/teams/:team_id/users/:user_id/roles/:role_id", handlers.AssignRole)
		api.DELETE("/teams/:team_id/users/:user_id/roles/:role_id", handlers.RemoveRole)
		api.GET("/teams/:team_id/users/:user_id/roles", handlers.GetUserRoles)
	}

	// Webhook for Telegram bot (only if using webhook)
	if useWebhook {
		r.POST("/webhook", func(c *gin.Context) {
			telegramBot.HandleWebhook(c.Writer, c.Request)
		})
	}

	// Serve static files for the web app
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")

	// Serve the main web app
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Valorant App",
		})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)

	// Start bot polling in background if not using webhook
	if !useWebhook {
		go telegramBot.StartPolling()
	}

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
