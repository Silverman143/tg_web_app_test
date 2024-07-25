package handler

import (
	"net/http"
	"project-2x/pkg/database"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getLeaderboardData(c *gin.Context, leaderboardType string) {
	telegramId, exists := GetTelegramID(c)
	if exists != nil {
		NewErrorResponse(c, http.StatusUnauthorized, "id is not exist in context")
		return
	}


	var (
		leaderboard []database.LeaderboardEntry
		rank        int
		err         error
	)

	switch leaderboardType {
	case "allTime":
		leaderboard, rank, err = h.service.LeaderBords.GetAllTimeLeaderboard(telegramId)
	case "currentMonth":
		leaderboard, rank, err = h.service.LeaderBords.GetCurrentMonthLeaderboard(telegramId)
	case "currentWeek":
		leaderboard, rank, err = h.service.LeaderBords.GetCurrentWeekLeaderboard(telegramId)
	default:
		NewErrorResponse(c, http.StatusBadRequest, "invalid leaderboard type")
		return
	}

	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "Can't get leaderboard: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"rank":       rank,
		"leaderboard": leaderboard,
	})
}

func (h *Handler) GetAllTimeLeaderboard(c *gin.Context) {
	h.getLeaderboardData(c, "allTime")
}

func (h *Handler) GetCurrentMonthLeaderboard(c *gin.Context) {
	h.getLeaderboardData(c, "currentMonth")
}

func (h *Handler) GetCurrentWeekLeaderboard(c *gin.Context) {
	h.getLeaderboardData(c, "currentWeek")
}