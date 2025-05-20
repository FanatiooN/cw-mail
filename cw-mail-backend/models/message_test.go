package models

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Ошибка подключения к тестовой БД: %v", err)
	}

	err = db.AutoMigrate(&User{}, &Message{})
	if err != nil {
		t.Fatalf("Ошибка миграции: %v", err)
	}

	// Создаем тестовых пользователей
	sender, err := CreateUser(db, "sender@example.com", "password")
	if err != nil {
		t.Fatalf("Ошибка создания отправителя: %v", err)
	}

	receiver, err := CreateUser(db, "receiver@example.com", "password")
	if err != nil {
		t.Fatalf("Ошибка создания получателя: %v", err)
	}

	return db
}

func TestSendMessage(t *testing.T) {
	db := setupTestDB(t)

	// Получаем отправителя
	var sender User
	err := db.Where("email = ?", "sender@example.com").First(&sender).Error
	if err != nil {
		t.Fatalf("Ошибка поиска отправителя: %v", err)
	}

	// Отправляем сообщение
	msg, err := SendMessage(db, sender.ID, "receiver@example.com", "Тестовое сообщение", "Текст сообщения", 0)
	if err != nil {
		t.Fatalf("Ошибка отправки сообщения: %v", err)
	}

	// Проверки
	if msg.Subject != "Тестовое сообщение" {
		t.Errorf("Неверная тема сообщения: %s", msg.Subject)
	}

	if msg.SenderID != sender.ID {
		t.Errorf("Неверный ID отправителя: %d", msg.SenderID)
	}

	if msg.Label != "inbox" {
		t.Errorf("Неверная метка сообщения: %s", msg.Label)
	}
}

func TestGetInboxMessages(t *testing.T) {
	db := setupTestDB(t)

	// Получаем пользователей
	var sender User
	var receiver User
	db.Where("email = ?", "sender@example.com").First(&sender)
	db.Where("email = ?", "receiver@example.com").First(&receiver)

	// Отправляем тестовое сообщение
	_, err := SendMessage(db, sender.ID, "receiver@example.com", "Тест входящих", "Текст", 0)
	if err != nil {
		t.Fatalf("Ошибка отправки сообщения: %v", err)
	}

	// Получаем входящие сообщения
	messages, err := GetInboxMessages(db, receiver.ID)
	if err != nil {
		t.Fatalf("Ошибка получения входящих: %v", err)
	}

	// Проверка
	if len(messages) != 1 {
		t.Errorf("Ожидалось 1 сообщение, получено %d", len(messages))
	}

	if len(messages) > 0 && messages[0].Subject != "Тест входящих" {
		t.Errorf("Неверная тема сообщения: %s", messages[0].Subject)
	}
}
