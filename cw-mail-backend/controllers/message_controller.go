package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mail-service/models"
	"github.com/mail-service/queue"
	"gorm.io/gorm"
)


type MessageController struct {
	DB          *gorm.DB
	NotifyQueue *queue.NotificationQueue
}


type SendMessageRequest struct {
	ReceiverEmail string `json:"receiver_email" binding:"required,email" example:"receiver@example.com"`
	Subject       string `json:"subject" binding:"required" example:"Важное сообщение"`
	Body          string `json:"body" binding:"required" example:"Текст сообщения содержит важную информацию"`
	ReadLimit     int    `json:"read_limit" example:"1"` // необязательное поле, 0 означает без ограничений
}


type UpdateLabelRequest struct {
	Label string `json:"label" binding:"required" example:"trash"`
}


func NewMessageController(db *gorm.DB, notifyQueue *queue.NotificationQueue, _ interface{}) *MessageController {
	return &MessageController{
		DB:          db,
		NotifyQueue: notifyQueue,
	}
}


// @Summary Отправить сообщение
// @Description Отправляет сообщение другому пользователю
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SendMessageRequest true "Данные для отправки сообщения"
// @Success 201 {object} models.Message "Созданное сообщение"
// @Failure 400 {object} map[string]string "Неверные данные запроса"
// @Failure 401 {object} map[string]string "Пользователь не аутентифицирован"
// @Failure 404 {object} map[string]string "Получатель не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /messages/send [post]
func (mc *MessageController) SendMessage(c *gin.Context) {
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверные данные запроса"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}


	if req.ReadLimit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "лимит прочтений не может быть отрицательным"})
		return
	}


	tx := mc.DB.Begin()

	message, err := models.SendMessage(tx, userID.(uint), req.ReceiverEmail, req.Subject, req.Body, req.ReadLimit)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "не удалось отправить сообщение: получатель не найден"})
		return
	}


	err = mc.NotifyQueue.PublishNewMessageNotification(message.ID, message.SenderID, message.ReceiverID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось отправить уведомление"})
		return
	}


	tx.Commit()

	c.JSON(http.StatusCreated, message)
}


// @Summary Получить входящие сообщения
// @Description Возвращает список входящих сообщений текущего пользователя
// @Tags messages
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Message "Список входящих сообщений"
// @Failure 401 {object} map[string]string "Пользователь не аутентифицирован"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /messages/inbox [get]
func (mc *MessageController) GetInbox(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	messages, err := models.GetInboxMessages(mc.DB, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить входящие сообщения"})
		return
	}

	c.JSON(http.StatusOK, messages)
}


// @Summary Получить отправленные сообщения
// @Description Возвращает список отправленных сообщений текущего пользователя
// @Tags messages
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Message "Список отправленных сообщений"
// @Failure 401 {object} map[string]string "Пользователь не аутентифицирован"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /messages/sent [get]
func (mc *MessageController) GetSent(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	messages, err := models.GetSentMessages(mc.DB, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить отправленные сообщения"})
		return
	}

	c.JSON(http.StatusOK, messages)
}


// @Summary Получить спам-сообщения
// @Description Возвращает список спам-сообщений текущего пользователя
// @Tags messages
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Message "Список спам-сообщений"
// @Failure 401 {object} map[string]string "Пользователь не аутентифицирован"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /messages/spam [get]
func (mc *MessageController) GetSpam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	messages, err := models.GetSpamMessages(mc.DB, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить спам-сообщения"})
		return
	}

	c.JSON(http.StatusOK, messages)
}


// @Summary Получить удаленные сообщения
// @Description Возвращает список удаленных сообщений текущего пользователя
// @Tags messages
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Message "Список удаленных сообщений"
// @Failure 401 {object} map[string]string "Пользователь не аутентифицирован"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /messages/trash [get]
func (mc *MessageController) GetTrash(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	messages, err := models.GetTrashMessages(mc.DB, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить удаленные сообщения"})
		return
	}

	c.JSON(http.StatusOK, messages)
}


// @Summary Обновить метку сообщения
// @Description Обновляет метку сообщения (например, пометить как спам или удаленное)
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID сообщения"
// @Param request body UpdateLabelRequest true "Данные для обновления метки"
// @Success 200 {object} map[string]string "Успешное обновление метки"
// @Failure 400 {object} map[string]string "Неверные данные запроса"
// @Failure 401 {object} map[string]string "Пользователь не аутентифицирован"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Сообщение не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /messages/{id}/label [put]
func (mc *MessageController) UpdateLabel(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	messageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	var req UpdateLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверные данные запроса"})
		return
	}


	validLabels := map[string]bool{
		"inbox": true,
		"spam":  true,
		"trash": true,
	}

	if !validLabels[req.Label] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверная метка"})
		return
	}


	err = models.UpdateMessageLabel(mc.DB, uint(messageID), userID.(uint), req.Label)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "сообщение не найдено или доступ запрещен"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось обновить метку сообщения"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "метка успешно обновлена"})
}


// @Summary Получить сообщение по ID
// @Description Возвращает детали сообщения по его идентификатору
// @Tags messages
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID сообщения"
// @Success 200 {object} models.Message "Сообщение"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 401 {object} map[string]string "Пользователь не аутентифицирован"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Сообщение не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /messages/{id} [get]
func (mc *MessageController) GetMessageByID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	messageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	var message models.Message
	if err := mc.DB.First(&message, messageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "сообщение не найдено"})
		return
	}


	if message.SenderID != userID.(uint) && message.ReceiverID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет доступа к этому сообщению"})
		return
	}


	if message.ReceiverID == userID.(uint) {

		if !message.IsRead {

			mc.DB.Model(&message).Update("is_read", true)


			shouldDelete, err := models.IncrementMessageReadCount(mc.DB, uint(messageID))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось обновить счетчик прочтений"})
				return
			}


			if shouldDelete {
				c.JSON(http.StatusOK, gin.H{
					"message": "Это сообщение было прочитано последний раз и удалено",
					"subject": message.Subject,
					"deleted": true,
				})
				return
			}


			if err := mc.DB.First(&message, messageID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "сообщение не найдено"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, message)
}


// @Summary Удалить просроченные сообщения
// @Description Удаляет сообщения, которые не были прочитаны в течение 24 часов
// @Tags messages
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Успешная очистка"
// @Failure 401 {object} map[string]string "Пользователь не аутентифицирован"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /messages/cleanup [post]
func (mc *MessageController) CleanupExpiredMessages(c *gin.Context) {
	err := models.DeleteExpiredMessages(mc.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось удалить просроченные сообщения"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "просроченные сообщения успешно удалены"})
}
