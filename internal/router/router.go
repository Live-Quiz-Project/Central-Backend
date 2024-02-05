package router

import (
	"time"

	"github.com/Live-Quiz-Project/Backend/internal/live"
	l "github.com/Live-Quiz-Project/Backend/internal/live/v1"
	"github.com/Live-Quiz-Project/Backend/internal/quiz"
	q "github.com/Live-Quiz-Project/Backend/internal/quiz/v1"
	"github.com/Live-Quiz-Project/Backend/internal/user"
	u "github.com/Live-Quiz-Project/Backend/internal/user/v1"
	"github.com/Live-Quiz-Project/Backend/internal/dashboard"
	d "github.com/Live-Quiz-Project/Backend/internal/dashboard/v1"
	"github.com/Live-Quiz-Project/Backend/internal/leaderboard"
	lb "github.com/Live-Quiz-Project/Backend/internal/leaderboard/v1"


	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Initialize(userHandler *u.Handler, quizHandler *q.Handler, liveHandler *l.Handler, dashboardHandler *d.Handler, leaderboardHandler *lb.Handler) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Cookie"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:5173" || origin == "http://localhost:5174" || origin == "http://localhost:3000" || origin == "http://127.0.0.1:5173"
		},
		MaxAge: 12 * time.Hour,
	}))

	v1 := r.Group("/v1")
	user.UserRoutes(v1, userHandler)
	quiz.QuizRoutes(v1, quizHandler)
	live.LiveRoutes(v1, liveHandler)
	dashboard.DashboardRoutes(v1, dashboardHandler)
	leaderboard.LeaderboardRoutes(v1, leaderboardHandler)
}

func Run(addr string) error {
	return r.Run(addr)
}
