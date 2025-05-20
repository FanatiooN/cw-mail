package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/mail-service/config"
	"github.com/mail-service/database"
	_ "github.com/mail-service/docs" // Импорт сгенерированных docs
	"github.com/mail-service/queue"
	"github.com/mail-service/routes"
)

// @title          Mail Service API
// @version        1.0
// @description    API сервиса обмена сообщениями
// @termsOfService http://swagger.io/terms/

// @contact.name  Mail Service API Support
// @contact.url   https://github.com/mail-service
// @contact.email support@mail-service.com

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен с префиксом 'Bearer '

func main() {

	migrate := flag.Bool("migrate", false, "Выполнить миграцию базы данных")
	flag.Parse()


	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}


	db, err := database.InitDB(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}


	if *migrate {
		if err := database.Migrate(db); err != nil {
			log.Fatalf("Ошибка миграции базы данных: %v", err)
		}
	}


	notifyQueue, err := queue.NewNotificationQueue(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к RabbitMQ: %v", err)
	}
	defer notifyQueue.Close()


	router := gin.Default()


	router.LoadHTMLGlob(filepath.Join("templates", "*.html"))


	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})


	routes.SetupRoutes(router, db, cfg, notifyQueue)


	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Сервер запущен на %s", serverAddr)
	log.Printf("Документация Swagger доступна по адресу http://localhost:%s/swagger/index.html", cfg.Server.Port)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
