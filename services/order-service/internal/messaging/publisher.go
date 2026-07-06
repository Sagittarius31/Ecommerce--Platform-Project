package messaging
import ("context";"encoding/json"; amqp "github.com/rabbitmq/amqp091-go"; "go.uber.org/zap")
type Publisher struct { conn *amqp.Connection; channel *amqp.Channel; logger *zap.Logger }
func NewPublisher(url string, logger *zap.Logger) (*Publisher, error) {
	conn, err := amqp.Dial(url); if err != nil { return nil, err }
	ch, err := conn.Channel(); if err != nil { conn.Close(); return nil, err }
	ch.ExchangeDeclare("ecommerce","topic",true,false,false,false,nil)
	return &Publisher{conn:conn,channel:ch,logger:logger}, nil
}
func (p *Publisher) Publish(ctx context.Context, key string, payload interface{}) error {
	body, _ := json.Marshal(payload)
	return p.channel.PublishWithContext(ctx,"ecommerce",key,false,false,amqp.Publishing{ContentType:"application/json",DeliveryMode:amqp.Persistent,Body:body})
}
func (p *Publisher) Close() { if p.channel != nil { p.channel.Close() }; if p.conn != nil { p.conn.Close() } }
