package router

import (
	"os"
	"strings"
	"time"

	"github.com/Live-Quiz-Project/Backend/internal/dashboard"
	d "github.com/Live-Quiz-Project/Backend/internal/dashboard/v1"
	"github.com/Live-Quiz-Project/Backend/internal/live"
	l "github.com/Live-Quiz-Project/Backend/internal/live/v1"
	"github.com/Live-Quiz-Project/Backend/internal/quiz"
	q "github.com/Live-Quiz-Project/Backend/internal/quiz/v1"
	"github.com/Live-Quiz-Project/Backend/internal/user"
	u "github.com/Live-Quiz-Project/Backend/internal/user/v1"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Initialize(userHandler *u.Handler, quizHandler *q.Handler, liveHandler *l.Handler, dashboardHandler *d.Handler) {
	r = gin.Default()
	allowOriginsEnv := os.Getenv("ALLOW_ORIGINS")
	allowOrigins := strings.Split(allowOriginsEnv, ",")

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Cookie"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			// return origin == "http://localhost:5173" || origin == "http://localhost:5174" || origin == "http://localhost:3000" || origin == "http://127.0.0.1:5173"
			for _, allowedOrigin := range allowOrigins {
				if allowedOrigin == origin {
					return true
				}
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	}))

	v1 := r.Group("/v1")
	user.UserRoutes(v1, userHandler)
	quiz.QuizRoutes(v1, quizHandler)
	live.LiveRoutes(v1, liveHandler)
	dashboard.DashboardRoutes(v1, dashboardHandler)
}

func Run(addr string) error {
	return r.Run(addr)
}
