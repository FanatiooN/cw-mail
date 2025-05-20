package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mail-service/config"
	"github.com/mail-service/middleware"
	"github.com/mail-service/models"
	"gorm.io/gorm"
)


type AuthController struct {
	DB     *gorm.DB
	Config *config.Config
}


type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
}


type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}


type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}


func NewAuthController(db *gorm.DB, cfg *config.Config) *AuthController {
	return &AuthController{
		DB:     db,
		Config: cfg,
	}
}


// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} TokenResponse "JWT токен"
// @Failure 400 {object} map[string]string "Неверные данные запроса"
// @Failure 409 {object} map[string]string "Пользователь с таким email уже существует"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверные данные запроса"})
		return
	}


	email := strings.ToLower(strings.TrimSpace(req.Email))


	var existingUser models.User
	if result := ac.DB.Where("email = ?", email).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "пользователь с таким email уже существует"})
		return
	}


	user, err := models.CreateUser(ac.DB, email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать пользователя"})
		return
	}


	token, err := middleware.GenerateToken(user, ac.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сгенерировать токен"})
		return
	}

	c.JSON(http.StatusCreated, TokenResponse{Token: token})
}


// @Summary Вход в систему
// @Description Аутентифицирует пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для входа"
// @Success 200 {object} TokenResponse "JWT токен"
// @Failure 400 {object} map[string]string "Неверные данные запроса"
// @Failure 401 {object} map[string]string "Неверные учетные данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверные данные запроса"})
		return
	}


	email := strings.ToLower(strings.TrimSpace(req.Email))


	user, err := models.FindUserByEmail(ac.DB, email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверные учетные данные"})
		return
	}


	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверные учетные данные"})
		return
	}


	token, err := middleware.GenerateToken(user, ac.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сгенерировать токен"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{Token: token})
}
