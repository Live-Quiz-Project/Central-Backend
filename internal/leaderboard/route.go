package leaderboard

import (
	"github.com/Live-Quiz-Project/Backend/internal/middleware"
	d "github.com/Live-Quiz-Project/Backend/internal/leaderboard/v1"
	"github.com/gin-gonic/gin"
)

func LeaderboardRoutes(r *gin.RouterGroup, h *d.Handler) {
	leaderboard := r.Group("/leaderboard")
	leaderboard.Use(middleware.UserRequiredAuthentication)

	leaderboard.GET("/:id", h.GetPlayersLeaderboard)
}