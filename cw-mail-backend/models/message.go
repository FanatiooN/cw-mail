package models

import (
	"time"

	"gorm.io/gorm"
)


type Message struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SenderID   uint      `json:"sender_id" gorm:"index"`
	ReceiverID uint      `json:"receiver_id" gorm:"index"`
	Subject    string    `json:"subject"`
	Body       string    `json:"body"`
	IsRead     bool      `json:"is_read" gorm:"default:false"`
	Label      string    `json:"label" gorm:"default:'inbox'"`
	ReadLimit  int       `json:"read_limit" gorm:"default:0"`
	ReadCount  int       `json:"read_count" gorm:"default:0"`
	ExpiresAt  time.Time `json:"expires_at,omitempty" gorm:"index"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	Sender     User      `json:"sender" gorm:"foreignKey:SenderID"`
	Receiver   User      `json:"receiver" gorm:"foreignKey:ReceiverID"`
}


func SendMessage(db *gorm.DB, senderID uint, receiverEmail, subject, body string, readLimit int) (*Message, error) {
	var receiver User
	if err := db.Where("email = ?", receiverEmail).First(&receiver).Error; err != nil {
		return nil, err
	}

	message := &Message{
		SenderID:   senderID,
		ReceiverID: receiver.ID,
		Subject:    subject,
		Body:       body,
		IsRead:     false,
		Label:      "inbox",
		ReadLimit:  readLimit,
		ReadCount:  0,
	}



	if readLimit > 0 {
		message.ExpiresAt = time.Now().Add(24 * time.Hour)
	}

	if err := db.Create(message).Error; err != nil {
		return nil, err
	}

	return message, nil
}


func GetInboxMessages(db *gorm.DB, userID uint) ([]Message, error) {
	var messages []Message
	err := db.Preload("Sender").Preload("Receiver").
		Where("receiver_id = ? AND label = ?", userID, "inbox").
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}


func GetSentMessages(db *gorm.DB, userID uint) ([]Message, error) {
	var messages []Message
	err := db.Preload("Sender").Preload("Receiver").
		Where("sender_id = ?", userID).
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}


func GetSpamMessages(db *gorm.DB, userID uint) ([]Message, error) {
	var messages []Message
	err := db.Preload("Sender").Preload("Receiver").
		Where("receiver_id = ? AND label = ?", userID, "spam").
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}


func GetTrashMessages(db *gorm.DB, userID uint) ([]Message, error) {
	var messages []Message
	err := db.Preload("Sender").Preload("Receiver").
		Where("(receiver_id = ? OR sender_id = ?) AND label = ?", userID, userID, "trash").
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}


func UpdateMessageLabel(db *gorm.DB, messageID uint, userID uint, label string) error {

	var message Message
	if err := db.First(&message, messageID).Error; err != nil {
		return err
	}


	if message.SenderID != userID && message.ReceiverID != userID {
		return gorm.ErrRecordNotFound
	}


	return db.Model(&message).Update("label", label).Error
}


func IncrementMessageReadCount(db *gorm.DB, messageID uint) (bool, error) {
	var message Message
	if err := db.First(&message, messageID).Error; err != nil {
		return false, err
	}


	message.ReadCount++


	if message.ReadLimit > 0 && message.ReadCount >= message.ReadLimit {
		return true, db.Delete(&message).Error
	}


	return false, db.Model(&message).Update("read_count", message.ReadCount).Error
}


func DeleteExpiredMessages(db *gorm.DB) error {
	return db.Where("expires_at < ? AND expires_at IS NOT NULL", time.Now()).Delete(&Message{}).Error
}
