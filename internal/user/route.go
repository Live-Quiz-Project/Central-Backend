package user

import (
	"github.com/Live-Quiz-Project/Backend/internal/middleware"
	u "github.com/Live-Quiz-Project/Backend/internal/user/v1"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup, h *u.Handler) {
	r.POST("/login", h.LogIn)
	r.GET("/logout", h.LogOut)
	r.GET("/refresh", h.RefreshToken)
	r.GET("/decode", middleware.UserRequiredAuthentication, h.DecodeToken)
	r.POST("/google-signin", h.GoogleSignIn)
	r.POST("/otp", h.SendOTPCode)
	r.POST("/verify-otp", h.VerifyOTPCode)

	userR := r.Group("/users")
	userR.POST("", h.CreateUser)
	userR.Use(middleware.UserRequiredAuthentication)
	userR.GET("", h.GetUsers)
	userR.GET("/:id", h.GetUserByID)
	userR.PATCH("/:id", h.UpdateUser)
	userR.DELETE("/:id", h.DeleteUser)
	userR.PATCH("/reset-password/:id", h.ChangePassword)

	admin := r.Group("/admin")
	admin.GET("/restore/:id", h.RestoreUser)
}
