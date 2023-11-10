package router

import (
	"time"

	"github.com/Live-Quiz-Project/Backend/internal/live"
	l "github.com/Live-Quiz-Project/Backend/internal/live/v1"
	"github.com/Live-Quiz-Project/Backend/internal/quiz"
	q "github.com/Live-Quiz-Project/Backend/internal/quiz/v1"
	"github.com/Live-Quiz-Project/Backend/internal/response"
	res "github.com/Live-Quiz-Project/Backend/internal/response/v1"
	"github.com/Live-Quiz-Project/Backend/internal/user"
	u "github.com/Live-Quiz-Project/Backend/internal/user/v1"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Initialize(userHandler *u.Handler, quizHandler *q.Handler, responseHandler *res.Handler, liveHandler *l.Handler) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:5173" || origin == "http://localhost:5174" || origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	v1 := r.Group("/v1")
	user.UserRoutes(v1, userHandler)
	quiz.QuizRoutes(v1, quizHandler)
	live.LiveRoutes(v1, liveHandler)
	response.ResponseRoutes(v1, responseHandler)
}

func Run(addr string) error {
	return r.Run(addr)
}
