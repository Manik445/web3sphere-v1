package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/logger"
)

// Topic constants.
const (
	TopicUserEvents     = "web3sphere.user.events"
	TopicPaymentEvents  = "web3sphere.payment.events"
	TopicBlockchain     = "web3sphere.blockchain.events"
	TopicAnalytics      = "web3sphere.analytics.events"
	TopicNotifications  = "web3sphere.notification.events"
)

// Producer wraps the Sarama sync producer.
type Producer struct {
	producer sarama.SyncProducer
	log      *logger.Logger
}

// NewProducer creates a new Kafka producer.
func NewProducer(cfg *configs.KafkaConfig, log *logger.Logger) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(cfg.BrokerList(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	log.Info("Kafka producer connected")
	return &Producer{producer: producer, log: log}, nil
}

// Publish sends a message to a Kafka topic.
func (p *Producer) Publish(ctx context.Context, topic string, key string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal Kafka message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(body),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to publish to Kafka topic %s: %w", topic, err)
	}

	p.log.Debugf("Kafka message published to %s [partition=%d, offset=%d]", topic, partition, offset)
	return nil
}

// Close shuts down the Kafka producer.
func (p *Producer) Close() {
	if err := p.producer.Close(); err != nil {
		p.log.Errorf("Failed to close Kafka producer: %v", err)
	} else {
		p.log.Info("Kafka producer closed")
	}
}

// Consumer wraps the Sarama consumer group.
type Consumer struct {
	group sarama.ConsumerGroup
	log   *logger.Logger
}

// MessageHandler handles consumed Kafka messages.
type MessageHandler func(ctx context.Context, key, value []byte) error

// consumerGroupHandler implements sarama.ConsumerGroupHandler.
type consumerGroupHandler struct {
	handler MessageHandler
	log     *logger.Logger
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h.handler(session.Context(), msg.Key, msg.Value); err != nil {
			h.log.Errorf("Failed to process Kafka message: %v", err)
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}

// NewConsumer creates a new Kafka consumer group.
func NewConsumer(cfg *configs.KafkaConfig, log *logger.Logger) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}

	if cfg.AutoOffsetReset == "earliest" {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	group, err := sarama.NewConsumerGroup(cfg.BrokerList(), cfg.GroupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer group: %w", err)
	}

	log.Info("Kafka consumer group connected")
	return &Consumer{group: group, log: log}, nil
}

// Consume starts consuming from the given topics.
func (c *Consumer) Consume(ctx context.Context, topics []string, handler MessageHandler) error {
	h := &consumerGroupHandler{handler: handler, log: c.log}

	go func() {
		for {
			if err := c.group.Consume(ctx, topics, h); err != nil {
				c.log.Errorf("Kafka consumer error: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

// Close shuts down the Kafka consumer.
func (c *Consumer) Close() {
	if err := c.group.Close(); err != nil {
		c.log.Errorf("Failed to close Kafka consumer: %v", err)
	} else {
		c.log.Info("Kafka consumer closed")
	}
}
