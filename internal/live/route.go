package live

import (
	l "github.com/Live-Quiz-Project/Backend/internal/live/v1"
	"github.com/Live-Quiz-Project/Backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func LiveRoutes(r *gin.RouterGroup, h *l.Handler) {
	r.POST("live", middleware.UserRequiredAuthentication, h.CreateLiveQuizSession)
	liveR := r.Group("/live/:code")
	liveR.GET("/check", h.CheckLiveQuizSessionAvailability)
	liveR.GET("/join", h.JoinLiveQuizSession)
	liveR.GET("/mod", h.UpdateModerator)
	liveR.GET("/interrupt", h.InterruptCountdown)
	liveR.GET("/end", middleware.UserRequiredAuthentication, h.EndLiveQuizSession)
}
