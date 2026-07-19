package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/logger"
)

// Queue name constants.
const (
	QueueEmail        = "web3sphere.email"
	QueueNotification = "web3sphere.notification"
	QueueAnalytics    = "web3sphere.analytics"
	QueueEscrow       = "web3sphere.escrow"
	QueuePayment      = "web3sphere.payment"
	QueueBlockchain   = "web3sphere.blockchain"
)

// Dead letter queue suffix.
const DLQSuffix = ".dlq"

// Client wraps the RabbitMQ connection and channel.
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     *configs.RabbitMQConfig
	log     *logger.Logger
}

// MessageHandler is a callback for processing consumed messages.
type MessageHandler func(ctx context.Context, body []byte) error

// New creates a new RabbitMQ client and declares all queues.
func New(cfg *configs.RabbitMQConfig, log *logger.Logger) (*Client, error) {
	conn, err := amqp.Dial(cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}

	if err := ch.Qos(cfg.PrefetchCount, 0, false); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	client := &Client{conn: conn, channel: ch, cfg: cfg, log: log}

	// Declare all queues with their DLQs
	queues := []string{
		QueueEmail, QueueNotification, QueueAnalytics,
		QueueEscrow, QueuePayment, QueueBlockchain,
	}

	for _, q := range queues {
		if err := client.declareQueueWithDLQ(q); err != nil {
			client.Close()
			return nil, err
		}
	}

	log.Info("RabbitMQ connected and queues declared")
	return client, nil
}

// declareQueueWithDLQ creates a queue and its dead-letter queue.
func (c *Client) declareQueueWithDLQ(name string) error {
	dlqName := name + DLQSuffix

	// Declare DLQ first
	_, err := c.channel.QueueDeclare(dlqName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ %s: %w", dlqName, err)
	}

	// Declare main queue with DLQ routing
	args := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": dlqName,
	}
	_, err = c.channel.QueueDeclare(name, true, false, false, false, args)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", name, err)
	}

	return nil
}

// Publish sends a message to a queue.
func (c *Client) Publish(ctx context.Context, queue string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return c.channel.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
	})
}

// Consume starts consuming messages from a queue.
func (c *Client) Consume(ctx context.Context, queue string, handler MessageHandler) error {
	msgs, err := c.channel.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to start consumer for %s: %w", queue, err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.log.Infof("Consumer for %s stopped", queue)
				return
			case msg, ok := <-msgs:
				if !ok {
					c.log.Warnf("Channel closed for queue %s", queue)
					return
				}
				if err := handler(ctx, msg.Body); err != nil {
					c.log.Errorf("Failed to process message from %s: %v", queue, err)
					msg.Nack(false, false) // Send to DLQ
				} else {
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

// Close shuts down the RabbitMQ connection.
func (c *Client) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	c.log.Info("RabbitMQ connection closed")
}
