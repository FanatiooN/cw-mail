package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mail-service/config"
	"github.com/mail-service/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{})
	return db
}

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userRole       string
		allowedRoles   []string
		expectedStatus int
	}{
		{"admin can access admin route", models.RoleAdmin, []string{models.RoleAdmin}, http.StatusOK},
		{"moderator can access moderator route", models.RoleModerator, []string{models.RoleModerator}, http.StatusOK},
		{"admin can access admin or moderator route", models.RoleAdmin, []string{models.RoleAdmin, models.RoleModerator}, http.StatusOK},
		{"moderator can access admin or moderator route", models.RoleModerator, []string{models.RoleAdmin, models.RoleModerator}, http.StatusOK},
		{"user cannot access admin route", models.RoleUser, []string{models.RoleAdmin}, http.StatusForbidden},
		{"user cannot access moderator route", models.RoleUser, []string{models.RoleModerator}, http.StatusForbidden},
		{"user cannot access admin or moderator route", models.RoleUser, []string{models.RoleAdmin, models.RoleModerator}, http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			// Middleware который устанавливает роль в контекст
			router.Use(func(c *gin.Context) {
				c.Set("user_role", tt.userRole)
				c.Next()
			})

			router.GET("/test", RequireRole(tt.allowedRoles...), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Ожидался статус %d, получен %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRequireAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		userRole       string
		expectedStatus int
	}{
		{models.RoleAdmin, http.StatusOK},
		{models.RoleModerator, http.StatusForbidden},
		{models.RoleUser, http.StatusForbidden},
		{"invalid_role", http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run("role_"+tt.userRole, func(t *testing.T) {
			router := gin.New()

			router.Use(func(c *gin.Context) {
				c.Set("user_role", tt.userRole)
				c.Next()
			})

			router.GET("/admin", RequireAdmin(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "admin access"})
			})

			req := httptest.NewRequest("GET", "/admin", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Для роли %s ожидался статус %d, получен %d", tt.userRole, tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRequireAdminOrModerator(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		userRole       string
		expectedStatus int
	}{
		{models.RoleAdmin, http.StatusOK},
		{models.RoleModerator, http.StatusOK},
		{models.RoleUser, http.StatusForbidden},
		{"invalid_role", http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run("role_"+tt.userRole, func(t *testing.T) {
			router := gin.New()

			router.Use(func(c *gin.Context) {
				c.Set("user_role", tt.userRole)
				c.Next()
			})

			router.GET("/admin-or-mod", RequireAdminOrModerator(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "admin or moderator access"})
			})

			req := httptest.NewRequest("GET", "/admin-or-mod", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Для роли %s ожидался статус %d, получен %d", tt.userRole, tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRoleMiddlewareWithoutRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/test", RequireRole(models.RoleAdmin), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusUnauthorized, w.Code)
	}
}

func TestGenerateTokenWithInvalidRole(t *testing.T) {
	cfg := &config.Config{
		JWT: struct {
			Secret     string
			Expiration time.Duration
		}{
			Secret:     "test_secret",
			Expiration: time.Hour,
		},
	}


	user := &models.User{
		ID:    1,
		Email: "test@example.com",
		Role:  "invalid_role",
	}

	_, err := GenerateToken(user, cfg)
	if err == nil {
		t.Error("Ожидалась ошибка при генерации токена для пользователя с невалидной ролью")
	}
}

func TestJWTAuthMiddlewareWithInvalidRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	cfg := &config.Config{
		JWT: struct {
			Secret     string
			Expiration time.Duration
		}{
			Secret:     "test_secret",
			Expiration: time.Hour,
		},
	}

	user, _ := models.CreateUser(db, "test@example.com", "password123")

	router := gin.New()
	router.Use(JWTAuthMiddleware(cfg, db))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   "invalid_role",

	token, err := generateTestToken(claims, cfg.JWT.Secret)
	if err != nil {
		t.Fatalf("Ошибка создания тестового токена: %v", err)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус %d для токена с невалидной ролью, получен %d", http.StatusUnauthorized, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "недействительная роль в токене" {
		t.Errorf("Неверное сообщение об ошибке: %v", response["error"])
	}
}

func TestRoleMismatchBetweenTokenAndDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	cfg := &config.Config{
		JWT: struct {
			Secret     string
			Expiration time.Duration
		}{
			Secret:     "test_secret",
			Expiration: time.Hour,
		},
	}

	admin, _ := models.CreateUserWithRole(db, "admin@example.com", "password123", models.RoleAdmin)

	db.Model(admin).Update("role", models.RoleUser)

	router := gin.New()
	router.Use(JWTAuthMiddleware(cfg, db))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	claims := &JWTClaims{
		UserID: admin.ID,
		Email:  admin.Email,
		Role:   models.RoleAdmin,
	}

	token, err := generateTestToken(claims, cfg.JWT.Secret)
	if err != nil {
		t.Fatalf("Ошибка создания тестового токена: %v", err)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус %d при несоответствии ролей, получен %d", http.StatusUnauthorized, w.Code)
	}
}

func TestRoleAttacksThroughJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	cfg := &config.Config{
		JWT: struct {
			Secret     string
			Expiration time.Duration
		}{
			Secret:     "test_secret",
			Expiration: time.Hour,
		},
	}

	user, _ := models.CreateUser(db, "user@example.com", "password123")

	router := gin.New()
	router.Use(JWTAuthMiddleware(cfg, db))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	maliciousRoles := []string{
		"'; DROP TABLE users; --",
		"<script>alert('xss')</script>",
		"admin\x00user",
		"ADMIN",
		"admin ",
		" admin",
	}

	for _, role := range maliciousRoles {
		claims := &JWTClaims{
			UserID: user.ID,
			Email:  user.Email,
			Role:   role,
		}

		token, err := generateTestToken(claims, cfg.JWT.Secret)
		if err != nil {
			continue 
		}

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Потенциально опасная роль %q прошла валидацию", role)
		}
	}
}

func generateTestToken(claims *JWTClaims, secret string) (string, error) {
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.NotBefore = jwt.NewNumericDate(time.Now())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
