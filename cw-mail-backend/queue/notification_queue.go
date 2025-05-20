package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mail-service/config"
	amqp "github.com/rabbitmq/amqp091-go"
)


const (
	ExchangeName   = "mail_notifications"
	QueueName      = "notifications"
	RoutingKey     = "new_message"
	PublishTimeout = 5 * time.Second
)


type NewMessageNotification struct {
	MessageID  uint      `json:"message_id"`
	SenderID   uint      `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	Timestamp  time.Time `json:"timestamp"`
}


type NotificationQueue struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}


func NewNotificationQueue(cfg *config.Config) (*NotificationQueue, error) {
	conn, err := amqp.Dial(cfg.GetRabbitMQURI())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}


	err = ch.ExchangeDeclare(
		ExchangeName, // имя
		"direct",     // тип
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // аргументы
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare an exchange: %w", err)
	}


	_, err = ch.QueueDeclare(
		QueueName, // имя
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // аргументы
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}


	err = ch.QueueBind(
		QueueName,    // имя очереди
		RoutingKey,   // routing key
		ExchangeName, // имя exchange
		false,        // no-wait
		nil,          // аргументы
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &NotificationQueue{
		Connection: conn,
		Channel:    ch,
	}, nil
}


func (nq *NotificationQueue) PublishNewMessageNotification(messageID, senderID, receiverID uint) error {
	notification := NewMessageNotification{
		MessageID:  messageID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Timestamp:  time.Now(),
	}

	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), PublishTimeout)
	defer cancel()

	err = nq.Channel.PublishWithContext(
		ctx,
		ExchangeName, // exchange
		RoutingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish notification: %w", err)
	}

	return nil
}


func (nq *NotificationQueue) Close() error {
	if err := nq.Channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := nq.Connection.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}
