package live

import (
	l "github.com/Live-Quiz-Project/Backend/internal/live/v1"
	"github.com/Live-Quiz-Project/Backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func LiveRoutes(r *gin.RouterGroup, h *l.Handler) {
	liveR := r.Group("/live")
	liveR.Use(middleware.UserRequiredAuthentication)
	liveR.POST("", h.CreateLiveQuizSession)
	// liveR.GET("", h.GetLiveQuizSessions)
	// liveR.GET("/:id", h.GetLiveQuizSessionByID)
	// liveR.PUT("/:id", h.UpdateLiveQuizSession)
	// liveR.DELETE("/:id", h.DeleteLiveQuizSession)
	liveR.GET("/:id/end", h.EndLiveQuizSession)
	liveR.GET("/:id/check", h.CheckLiveQuizSessionAvailability)
	liveR.GET("/host", h.GetHost)
	liveR.GET("/participants", h.GetParticipants)

	r.GET("live/join/:code", h.JoinLiveQuizSession)
}
