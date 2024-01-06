package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/Live-Quiz-Project/Backend/internal/cache"
	"github.com/Live-Quiz-Project/Backend/internal/db"
	"github.com/Live-Quiz-Project/Backend/internal/env"
	l "github.com/Live-Quiz-Project/Backend/internal/live/v1"
	q "github.com/Live-Quiz-Project/Backend/internal/quiz/v1"
	"github.com/Live-Quiz-Project/Backend/internal/router"
	u "github.com/Live-Quiz-Project/Backend/internal/user/v1"
	d "github.com/Live-Quiz-Project/Backend/internal/dashboard/v1"
)

func main() {
	env.Initialize()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Error initializing database connection: %v", err)
	}

	cacheConn, er := cache.NewCache(ctx)
	if er != nil {
		log.Fatalf("Error initializing cache connection: %v", er)
	}

	defer func() {
		dbConn.Close()
		cacheConn.Close()
		cancel()
	}()

	uRepo := u.NewRepository(dbConn.GetDB())
	uServ := u.NewService(uRepo)
	userHandler := u.NewHandler(uServ)

	qRepo := q.NewRepository(dbConn.GetDB())
	qServ := q.NewService(qRepo)
	quizHandler := q.NewHandler(qServ)

	hub := l.NewHub()
	lRepo := l.NewRepository(dbConn.GetDB(), cacheConn.GetCache())
	lServ := l.NewService(lRepo)

	liveHandler := l.NewHandler(hub, lServ, qServ)

	dashboardRepo := d.NewRepository(dbConn.GetDB())
	dashboardServ := d.NewService(dashboardRepo)
	dashboardHandler := d.NewHandler(dashboardServ)

	go hub.Run()
	router.Initialize(userHandler, quizHandler, liveHandler, dashboardHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
