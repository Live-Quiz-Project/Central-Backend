package quiz

import (
	"github.com/Live-Quiz-Project/Backend/internal/middleware"
	q "github.com/Live-Quiz-Project/Backend/internal/quiz/v1"
	"github.com/gin-gonic/gin"
)

func QuizRoutes(r *gin.RouterGroup, h *q.Handler) {
	quizR := r.Group("/quizzes")
	quizR.Use(middleware.UserRequiredAuthentication)
	quizR.POST("", h.CreateQuiz)           // Use for Create Quiz (include nest inside)
	quizR.GET("", h.GetQuizzes)            // Use for Get All Quiz From User
	quizR.GET("/:id", h.GetQuizByID)       // Use for Get All Quiz Detail By QuizID
	quizR.PUT("/:id", h.UpdateQuiz)        // Use for Update Quiz
	quizR.DELETE("/:id", h.DeleteQuiz)     // Use for Soft Delete Quiz
	quizR.PATCH("/:id", h.RestoreQuiz) 		 // Use for Restore Quiz

	quizR.GET("/history", h.GetQuizHistories)
	quizR.GET("/history/:id", h.GetQuizHistoryByID)
}
