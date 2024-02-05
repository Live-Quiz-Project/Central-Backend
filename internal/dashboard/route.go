package dashboard

import (
	d "github.com/Live-Quiz-Project/Backend/internal/dashboard/v1"
	"github.com/Live-Quiz-Project/Backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func DashboardRoutes(r *gin.RouterGroup, h *d.Handler) {
	dashboard := r.Group("/dashboard")
	dashboard.Use(middleware.UserRequiredAuthentication)

	dashboard.GET("/live/:id", h.GetAnswerResponseByLiveQuizSessionID)
	// dashboard.GET("/participant/:id", h.GetAnswerResponseByParticipantID)
	// dashboard.GET("/question/:id", h.GetAnswerResponseByQuestionID)
	dashboard.GET("", h.GetDashboardHistoryByUserID)
	dashboard.GET("/question/:id", h.GetDashboardQuestionViewByID)
	dashboard.GET("/answer/:id", h.GetDashboardAnswerViewByID)
}
