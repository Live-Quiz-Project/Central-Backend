package live

import (
	l "github.com/Live-Quiz-Project/Backend/internal/live/v1"
	"github.com/Live-Quiz-Project/Backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func LiveRoutes(r *gin.RouterGroup, h *l.Handler) {
	r.GET("live/check/:code", h.CheckLiveQuizSessionAvailability)
	r.GET("live/join/:code", h.JoinLiveQuizSession)

	liveR := r.Group("/live")
	liveR.Use(middleware.UserRequiredAuthentication)
	liveR.POST("", h.CreateLiveQuizSession)
	// liveR.GET("", h.GetLiveQuizSessions)
	// liveR.GET("/:id", h.GetLiveQuizSessionByID)
	// liveR.PUT("/:id", h.UpdateLiveQuizSession)
	// liveR.DELETE("/:id", h.DeleteLiveQuizSession)
	liveR.GET("/:id/end", h.EndLiveQuizSession)
	liveR.GET("/:id/host", h.GetHost)
	liveR.GET("/:id/participants", h.GetParticipants)

	liveR.GET("/:id/cache", h.GetLiveQuizSessionCache)
	liveR.GET("/:id/cache/responses", h.GetLiveQuizSessionResponsesCache)
}
