package consumer
import ("context";"encoding/json"; amqp "github.com/rabbitmq/amqp091-go"; "github.com/yourname/ecommerce/notification-service/internal/email"; "go.uber.org/zap")
type OrderPlacedEvent struct { OrderID,UserEmail string; TotalAmount float64 }
type PaymentSucceededEvent struct { OrderID,UserEmail string; Amount float64 }
type PaymentFailedEvent struct { OrderID,UserEmail,Reason string }
type NotificationConsumer struct { conn *amqp.Connection; channel *amqp.Channel; emailSvc *email.SMTPService; logger *zap.Logger }
func New(url string, emailSvc *email.SMTPService, logger *zap.Logger) (*NotificationConsumer, error) {
	conn,err:=amqp.Dial(url); if err!=nil{return nil,err}
	ch,err:=conn.Channel(); if err!=nil{conn.Close();return nil,err}
	ch.Qos(1,0,false)
	return &NotificationConsumer{conn:conn,channel:ch,emailSvc:emailSvc,logger:logger},nil
}
func (c *NotificationConsumer) Start(ctx context.Context) error {
	c.channel.ExchangeDeclare("ecommerce.dlx","topic",true,false,false,false,nil)
	for qName,key:=range map[string]string{"notification.order.placed":"order.placed","notification.payment.succeeded":"payment.succeeded","notification.payment.failed":"payment.failed"} {
		q,err:=c.channel.QueueDeclare(qName,true,false,false,false,amqp.Table{"x-dead-letter-exchange":"ecommerce.dlx"})
		if err!=nil{return err}
		c.channel.QueueBind(q.Name,key,"ecommerce",false,nil)
		msgs,err:=c.channel.Consume(q.Name,"",false,false,false,false,nil); if err!=nil{return err}
		go func(k string,deliveries <-chan amqp.Delivery){
			for{select{case <-ctx.Done():return; case msg,ok:=<-deliveries:if !ok{return};c.handle(k,msg)}}
		}(key,msgs)
	}
	c.logger.Info("notification consumer started"); return nil
}
func (c *NotificationConsumer) handle(key string, msg amqp.Delivery) {
	switch key {
	case "order.placed":
		var ev OrderPlacedEvent; if err:=json.Unmarshal(msg.Body,&ev);err!=nil{msg.Nack(false,false);return}
		if err:=c.emailSvc.SendOrderConfirmation(ev.UserEmail,ev.OrderID,ev.TotalAmount);err!=nil{msg.Nack(false,true);return}
	case "payment.succeeded":
		var ev PaymentSucceededEvent; if err:=json.Unmarshal(msg.Body,&ev);err!=nil{msg.Nack(false,false);return}
		if err:=c.emailSvc.SendPaymentReceipt(ev.UserEmail,ev.OrderID,ev.Amount);err!=nil{msg.Nack(false,true);return}
	case "payment.failed":
		var ev PaymentFailedEvent; if err:=json.Unmarshal(msg.Body,&ev);err!=nil{msg.Nack(false,false);return}
		if err:=c.emailSvc.SendPaymentFailure(ev.UserEmail,ev.OrderID,ev.Reason);err!=nil{msg.Nack(false,true);return}
	}
	msg.Ack(false)
}
func (c *NotificationConsumer) Close() { if c.channel!=nil{c.channel.Close()}; if c.conn!=nil{c.conn.Close()} }
