package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mail-service/config"
	"github.com/mail-service/controllers"
	"github.com/mail-service/middleware"
	"github.com/mail-service/queue"
	"gorm.io/gorm"
)


func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config, notifyQueue *queue.NotificationQueue) {

	authController := controllers.NewAuthController(db, cfg)
	userController := controllers.NewUserController(db)
	messageController := controllers.NewMessageController(db, notifyQueue, nil)


	api := router.Group("/api")
	{

		public := api.Group("")
		{
			public.POST("/auth/register", authController.Register)
			public.POST("/auth/login", authController.Login)
		}


		protected := api.Group("")
		protected.Use(middleware.JWTAuthMiddleware(cfg, db))
		{

			users := protected.Group("/users")
			{
				users.GET("/me", userController.GetCurrentUser)
			}


			auth := protected.Group("/auth")
			{
				auth.GET("/me", userController.GetCurrentUser)
			}


			messages := protected.Group("/messages")
			{
				messages.POST("", messageController.SendMessage)
				messages.GET("/inbox", messageController.GetInbox)
				messages.GET("/sent", messageController.GetSent)
				messages.GET("/spam", messageController.GetSpam)
				messages.GET("/trash", messageController.GetTrash)
				messages.GET("/:id", messageController.GetMessageByID)
				messages.PUT("/:id/label", messageController.UpdateLabel)
			}
		}
	}
}
