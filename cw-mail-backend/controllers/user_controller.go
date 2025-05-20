package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type UserController struct {
	DB *gorm.DB
}


type UserResponse struct {
	UserID uint `json:"user_id" example:"1"`
}


func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		DB: db,
	}
}


// @Summary Получить информацию о текущем пользователе
// @Description Возвращает идентификатор текущего пользователя
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse "Информация о пользователе"
// @Failure 401 {object} map[string]string "Пользователь не аутентифицирован"
// @Router /users/me [get]
func (uc *UserController) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		UserID: userID.(uint),
	})
}
