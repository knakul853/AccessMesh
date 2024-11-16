package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"time"
)

type EmailService struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	fromEmail    string
}

func NewEmailService(host string, port int, username, password, from string) *EmailService {
	return &EmailService{
		smtpHost:     host,
		smtpPort:     port,
		smtpUsername: username,
		smtpPassword: password,
		fromEmail:    from,
	}
}

func (s *EmailService) SendVerificationEmail(to, token string) error {
	subject := "Verify Your Email"
	body := fmt.Sprintf("Please verify your email by clicking this link: http://yourdomain.com/verify?token=%s", token)
	
	return s.sendEmail(to, subject, body)
}

func (s *EmailService) SendPasswordResetEmail(to, token string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf("Reset your password by clicking this link: http://yourdomain.com/reset-password?token=%s", token)
	
	return s.sendEmail(to, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
	
	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", s.fromEmail, to, subject, body)

	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)
	return smtp.SendMail(addr, auth, s.fromEmail, []string{to}, []byte(msg))
}

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

type Token struct {
	Token     string
	ExpiresAt time.Time
}

func NewToken(duration time.Duration) (*Token, error) {
	token, err := GenerateToken()
	if err != nil {
		return nil, err
	}
	
	return &Token{
		Token:     token,
		ExpiresAt: time.Now().Add(duration),
	}, nil
}
