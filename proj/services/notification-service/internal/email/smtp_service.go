package email
import ("fmt";"net/smtp";"go.uber.org/zap")
type Config struct { Host,Port,Username,Password,From string }
type SMTPService struct { cfg Config; logger *zap.Logger }
func New(cfg Config, logger *zap.Logger) *SMTPService { return &SMTPService{cfg:cfg,logger:logger} }
func (s *SMTPService) SendOrderConfirmation(to,orderID string, amount float64) error {
	return s.send(to,"Order Confirmed - "+orderID,fmt.Sprintf("Order %s placed. Total: $%.2f",orderID,amount))
}
func (s *SMTPService) SendPaymentReceipt(to,orderID string, amount float64) error {
	return s.send(to,"Payment Receipt - "+orderID,fmt.Sprintf("Payment $%.2f received for order %s",amount,orderID))
}
func (s *SMTPService) SendPaymentFailure(to,orderID,reason string) error {
	return s.send(to,"Payment Failed - "+orderID,fmt.Sprintf("Payment failed for order %s: %s",orderID,reason))
}
func (s *SMTPService) send(to,subject,body string) error {
	auth:=smtp.PlainAuth("",s.cfg.Username,s.cfg.Password,s.cfg.Host)
	msg:=[]byte("To: "+to+"\r\nSubject: "+subject+"\r\n\r\n"+body)
	if err:=smtp.SendMail(s.cfg.Host+":"+s.cfg.Port,auth,s.cfg.From,[]string{to},msg);err!=nil{s.logger.Error("smtp failed",zap.Error(err));return err}
	return nil
}
