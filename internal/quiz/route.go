package quiz

import (
	"github.com/Live-Quiz-Project/Backend/internal/middleware"
	q "github.com/Live-Quiz-Project/Backend/internal/quiz/v1"
	"github.com/gin-gonic/gin"
)

func QuizRoutes(r *gin.RouterGroup, h *q.Handler) {
	quizR := r.Group("/quizzes")
	quizR.Use(middleware.UserRequiredAuthentication)
	quizR.POST("", h.CreateQuiz)
	quizR.GET("", h.GetQuizzes)
	quizR.GET("/:id", h.GetQuizByID)
	quizR.PUT("/:id", h.UpdateQuiz)
	quizR.DELETE("/:id", h.DeleteQuiz)
}
