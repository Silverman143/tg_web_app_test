package handler

import (
	"time"

	"project-2x/pkg/service"
	"project-2x/pkg/telegramBot"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
	telegramBot *telegramBot.Bot
}

func NewHandler(s *service.Service, bot *telegramBot.Bot) *Handler{
	return &Handler{
		service: s,
	telegramBot: bot,
	}
}

func (h *Handler) InitRouts(bot *telegramBot.Bot) *gin.Engine{
	router := gin.New()
	
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	   }))

	router.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Expose-Headers", "Content-Length")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "43200")
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	auth:= router.Group("/auth") 
	{
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		api.GET("/profiles/my_profile/", h.GetUserData)

		api.GET("/daily-bonus", h.GetDailyBonys)
		api.POST("/daily-bonus/claim", h.ClaimDailyBonys)

		api.GET("/leaderboard/all-time", h.GetAllTimeLeaderboard)
		api.GET("/leaderboard/month", h.GetCurrentMonthLeaderboard)
		api.GET("/leaderboard/week", h.GetCurrentWeekLeaderboard)

		api.GET("/wallet/balance", h.GetBalance)
		api.GET("/wallet/transactions/all", h.GetAllTransactions)
		api.GET("/wallet/transactions/in", h.GetPositiveTransactions)
		api.GET("/wallet/transactions/out", h.GetNegativeTransactions)

		api.POST("/payment/stars-invoice-link", h.GetStarsPaymentURL)

		api.GET("/referral/url", h.GetReferralURL)
		api.POST("/referral/friends", h.GetUserReferrals)
	}

	// Добавьте маршрут для обработки вебхуков Telegram
	router.POST("/web_hook/"+bot.API.Self.UserName, h.TelegramBotHandler)
	
	return router
}