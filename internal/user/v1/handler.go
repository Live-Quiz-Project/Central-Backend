package v1

import (
	"net/http"
	"os"
	"time"

	"github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

// ---------- Auth related handlers ---------- //
func (h *Handler) LogIn(c *gin.Context) {
	var req LogInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, refreshToken, err := h.Service.LogIn(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("token", refreshToken, 60*60*24*7, "/", os.Getenv("HOST"), false, true)
	c.JSON(http.StatusOK, res)
}

func (h *Handler) LogOut(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if the refresh token is valid

	claims, err := util.DecodeToken(refreshToken, os.Getenv("REFRESH_TOKEN_SECRET"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	accessToken, err := util.GenerateToken(claims.UserID, claims.Name, claims.DisplayName, claims.DisplayColor, claims.DisplayEmoji, time.Now().Add(15*time.Minute), os.Getenv("ACCESS_TOKEN_SECRET"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (h *Handler) DecodeToken(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	uname, ok := c.Get("name")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	displayName, ok := c.Get("display_name")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	displayColor, ok := c.Get("display_color")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	displayEmoji, ok := c.Get("display_emoji")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, util.Claims{
		UserID:       userID,
		Name:         uname.(string),
		DisplayName:  displayName.(string),
		DisplayColor: displayColor.(string),
		DisplayEmoji: displayEmoji.(string),
	})
}

// ---------- User related handlers ---------- //
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, refreshToken, err := h.Service.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("token", refreshToken, 60*60*24*7, "/", os.Getenv("HOST"), false, true)
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetUsers(c *gin.Context) {
	if uid, ok := c.Get("uid"); ok {
		id, err := uuid.Parse(uid.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid id",
			})
			return
		}

		res, err := h.Service.GetUserByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, res)
		return
	}

	if _, ok := c.Get("admin"); ok {
		res, err := h.Service.GetUsers(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "unauthorized",
	})
}

func (h *Handler) GetUserByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	if _, ok := c.Get("uid"); ok {
		res, err := h.Service.GetUserByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "unauthorized",
	})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	if _, ok := c.Get("uid"); ok {
		res, err := h.Service.UpdateUser(c.Request.Context(), &req, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "unauthorized",
	})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	if uid, ok := c.Get("uid"); ok {
		if uid != id.String() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}
		err := h.Service.DeleteUser(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted user"})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "unauthorized",
	})
}

// ---------- Admin related handlers ---------- //
func (h *Handler) RestoreUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	if _, ok := c.Get("admin"); ok {
		err := h.Service.RestoreUser(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Successfully restored user"})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "unauthorized",
	})
}
