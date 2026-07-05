package email

import (
	"fmt"
	"net/smtp"
	"go.uber.org/zap"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type SMTPService struct {
	config Config
	logger *zap.Logger
}

func New(cfg Config, logger *zap.Logger) *SMTPService {
	return &SMTPService{config: cfg, logger: logger}
}

func (s *SMTPService) SendOrderConfirmation(to, orderID string, amount float64) error {
	subject := "Order Confirmation - " + orderID
	body    := fmt.Sprintf("Your order %s has been placed successfully. Total: $%.2f", orderID, amount)
	return s.send(to, subject, body)
}

func (s *SMTPService) SendPaymentReceipt(to, orderID string, amount float64) error {
	subject := "Payment Receipt - " + orderID
	body    := fmt.Sprintf("Payment of $%.2f received for order %s.", amount, orderID)
	return s.send(to, subject, body)
}

func (s *SMTPService) SendPaymentFailure(to, orderID, reason string) error {
	subject := "Payment Failed - " + orderID
	body    := fmt.Sprintf("Payment for order %s failed. Reason: %s", orderID, reason)
	return s.send(to, subject, body)
}

func (s *SMTPService) send(to, subject, body string) error {
	addr := s.config.Host + ":" + s.config.Port
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	msg  := []byte("To: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + body)
	if err := smtp.SendMail(addr, auth, s.config.From, []string{to}, msg); err != nil {
		s.logger.Error("smtp send failed", zap.String("to", to), zap.Error(err))
		return err
	}
	return nil
}
