package messaging

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *zap.Logger
}

func NewPublisher(amqpURL string, logger *zap.Logger) (*Publisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil { return nil, err }
	ch, err := conn.Channel()
	if err != nil { conn.Close(); return nil, err }
	err = ch.ExchangeDeclare("ecommerce", "topic", true, false, false, false, nil)
	if err != nil { ch.Close(); conn.Close(); return nil, err }
	return &Publisher{conn: conn, channel: ch, logger: logger}, nil
}

func (p *Publisher) Publish(ctx context.Context, routingKey string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil { return err }
	err = p.channel.PublishWithContext(ctx, "ecommerce", routingKey, false, false,
		amqp.Publishing{ContentType: "application/json", DeliveryMode: amqp.Persistent, Body: body})
	if err != nil {
		p.logger.Error("publish failed", zap.String("key", routingKey), zap.Error(err))
		return err
	}
	p.logger.Info("event published", zap.String("key", routingKey))
	return nil
}

func (p *Publisher) Close() {
	if p.channel != nil { p.channel.Close() }
	if p.conn != nil    { p.conn.Close() }
}
