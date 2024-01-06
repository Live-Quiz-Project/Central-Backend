package dashboard

import (
	"github.com/Live-Quiz-Project/Backend/internal/middleware"
	d "github.com/Live-Quiz-Project/Backend/internal/dashboard/v1"
	"github.com/gin-gonic/gin"
)

func DashboardRoutes(r *gin.RouterGroup, h *d.Handler) {
	dashboard := r.Group("/dashboard")
	dashboard.Use(middleware.UserRequiredAuthentication)

	dashboard.GET("/live/:id", h.GetAnswerResponseByLiveQuizSessionID)
	dashboard.GET("/participant/:id", h.GetAnswerResponseByParticipantID)
	dashboard.GET("/question/:id", h.GetAnswerResponseByQuestionID)
}